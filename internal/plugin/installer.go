package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// InstallSpec describes what to install. Most callers populate Source (and
// optionally Ref); LocalPath is reserved for development / testing.
//
// Source syntaxes:
//
//	github.com/owner/repo                      → https://github.com/owner/repo.git
//	github.com/owner/repo@v1.2.0               → clone + checkout v1.2.0
//	https://example.com/path/to/repo.git       → clone direct
//	git@github.com:owner/repo.git              → ssh clone (Ref via -ref flag only)
//
// When LocalPath is set it takes precedence over Source — a recursive copy
// happens in place of `git clone`. Useful for iterating on a plugin you are
// authoring without pushing every WIP commit to a remote.
type InstallSpec struct {
	Source    string // remote URL or `github.com/...` shorthand
	Ref       string // optional git ref (tag/branch/commit)
	LocalPath string // dev-only: copy from a local directory instead of cloning

	// OverrideID and OverrideVersion replace whatever the cloned plugin.json
	// declared. Empty means "trust the manifest".
	OverrideID      string
	OverrideVersion string
}

// Install fetches the plugin (clone or local copy), reads its plugin.json,
// validates the ID prefix rule, and moves the result into PluginsDir() under
// the plugin's ID. Returns the resolved Installed metadata so the CLI can
// confirm what landed.
//
// The fetch step writes to a private temp dir first, so a partially-cloned
// or invalid plugin never pollutes ~/.dot/plugins.
func Install(ctx context.Context, spec InstallSpec) (*Installed, error) {
	if spec.LocalPath == "" && spec.Source == "" {
		return nil, fmt.Errorf("plugin: install requires Source or LocalPath")
	}

	staging, err := os.MkdirTemp("", "dot-plugin-install-*")
	if err != nil {
		return nil, fmt.Errorf("plugin: staging dir: %w", err)
	}
	defer os.RemoveAll(staging)

	work := filepath.Join(staging, "src")

	if spec.LocalPath != "" {
		if err := copyDir(spec.LocalPath, work); err != nil {
			return nil, fmt.Errorf("plugin: copy %s: %w", spec.LocalPath, err)
		}
	} else {
		url, ref, err := parseSource(spec.Source, spec.Ref)
		if err != nil {
			return nil, err
		}
		if err := gitClone(ctx, url, ref, work); err != nil {
			return nil, fmt.Errorf("plugin: clone %s: %w", url, err)
		}
	}

	manifest, err := loadStagedManifest(work)
	if err != nil {
		return nil, err
	}

	if spec.OverrideID != "" {
		manifest.ID = spec.OverrideID
	}
	if spec.OverrideVersion != "" {
		manifest.Version = spec.OverrideVersion
	}
	if manifest.ID == "" {
		return nil, fmt.Errorf("plugin: plugin.json is missing required field 'id'")
	}
	if strings.Contains(manifest.ID, ".") {
		return nil, fmt.Errorf("plugin: id %q must not contain '.'", manifest.ID)
	}

	pluginsDir, err := PluginsDir()
	if err != nil {
		return nil, err
	}
	dst := filepath.Join(pluginsDir, manifest.ID)

	// Replace any existing install of the same id; users typically expect
	// `install` to upgrade in place.
	if err := os.RemoveAll(dst); err != nil {
		return nil, fmt.Errorf("plugin: clear %s: %w", dst, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return nil, fmt.Errorf("plugin: mkdir parent: %w", err)
	}
	if err := os.Rename(work, dst); err != nil {
		// Cross-device rename can fail; fall back to copy.
		if err := copyDir(work, dst); err != nil {
			return nil, fmt.Errorf("plugin: install %s: %w", manifest.ID, err)
		}
	}

	manifest.Dir = dst
	if err := writeManifest(dst, manifest); err != nil {
		return nil, err
	}
	return manifest, nil
}

// Uninstall removes a plugin's directory entirely. Missing IDs are a no-op.
func Uninstall(id string) error {
	pluginsDir, err := PluginsDir()
	if err != nil {
		return err
	}
	dst := filepath.Join(pluginsDir, id)
	if _, err := os.Stat(dst); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return os.RemoveAll(dst)
}

// ── Source parsing & git clone ────────────────────────────────────────────

// parseSource normalizes a user-supplied source string into (url, ref).
//
// Accepted forms:
//
//	github.com/owner/repo            → https://github.com/owner/repo.git
//	github.com/owner/repo@v1         → URL above, ref="v1"
//	https://example.com/x.git        → as-is
//	git@github.com:owner/repo.git    → as-is (ref must come via cliRef arg)
//
// cliRef wins over an embedded @ref so users can override on the command line.
func parseSource(source, cliRef string) (url, ref string, err error) {
	source = strings.TrimSpace(source)
	if source == "" {
		return "", "", fmt.Errorf("plugin: empty source")
	}

	// SSH URLs (git@host:path) — never split on @ here; ref must be explicit.
	if strings.HasPrefix(source, "git@") {
		return source, cliRef, nil
	}

	// Split off @ref if present (only for non-SSH inputs).
	embeddedRef := ""
	if i := strings.LastIndex(source, "@"); i >= 0 {
		embeddedRef = source[i+1:]
		source = source[:i]
	}
	resolvedRef := cliRef
	if resolvedRef == "" {
		resolvedRef = embeddedRef
	}

	switch {
	case strings.HasPrefix(source, "https://"), strings.HasPrefix(source, "http://"):
		return source, resolvedRef, nil

	case strings.HasPrefix(source, "github.com/"),
		strings.HasPrefix(source, "gitlab.com/"),
		strings.HasPrefix(source, "bitbucket.org/"):
		// Default to https for the major hosts so `dot plugin install
		// github.com/foo/bar` Just Works without forcing users to spell out
		// the protocol.
		return "https://" + source + ".git", resolvedRef, nil

	default:
		return "", "", fmt.Errorf("plugin: unsupported source %q (try github.com/owner/repo or full https/git URL)", source)
	}
}

// gitClone runs `git clone` and optional `git checkout`. Streams output to
// the user's terminal so clone progress is visible — git progress is the
// one operation where suppressing output is a UX regression.
//
// When ref is empty we use --depth=1 to keep installs fast.
func gitClone(ctx context.Context, url, ref, dst string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git not found in PATH; install git first")
	}

	args := []string{"clone"}
	if ref == "" {
		args = append(args, "--depth=1")
	}
	args = append(args, url, dst)

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Stdout = os.Stderr // keep stdout clean for shell pipes
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone: %w", err)
	}

	if ref != "" {
		co := exec.CommandContext(ctx, "git", "-C", dst, "checkout", ref)
		co.Stdout = os.Stderr
		co.Stderr = os.Stderr
		co.Env = os.Environ()
		if err := co.Run(); err != nil {
			return fmt.Errorf("git checkout %s: %w", ref, err)
		}
	}
	return nil
}

// ── manifest read/write helpers ───────────────────────────────────────────

// loadStagedManifest reads plugin.json from a freshly-cloned directory and
// returns the parsed Installed value (with Dir left blank — the caller fills
// it in once the plugin is moved into place).
func loadStagedManifest(dir string) (*Installed, error) {
	path := filepath.Join(dir, "plugin.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("plugin: %s/plugin.json missing — not a DOT plugin?", dir)
		}
		return nil, fmt.Errorf("plugin: read %s: %w", path, err)
	}
	var p Installed
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("plugin: parse %s: %w", path, err)
	}
	return &p, nil
}

func writeManifest(dir string, p *Installed) error {
	path := filepath.Join(dir, "plugin.json")
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("plugin: marshal manifest: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("plugin: write %s: %w", path, err)
	}
	return nil
}

// copyDir performs a recursive copy of src into dst. Symlinks are followed.
// Existing files in dst are overwritten.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip the .git directory — the plugin store does not need history.
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == ".git" || strings.HasPrefix(rel, ".git"+string(filepath.Separator)) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}
