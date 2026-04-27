// Package plugin manages community + local generator plugins. Plugins extend
// DOT by registering additional generators with the executor's Registry and
// additional flow injections (Replace / AddOption / InsertAfter) with the
// flow engine's HookRegistry.
//
// Discovery is on-disk: plugins live under ~/.dot/plugins/<plugin-id>/ and
// declare themselves via a plugin.json manifest. The loader walks that tree
// at startup and instantiates whatever it finds.
package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Installed represents one plugin discovered on disk. Loader populates this
// from the plugin.json next to the binary / module entry point.
type Installed struct {
	ID          string   `json:"id"`      // unique slug; also the on-disk dir name
	Version     string   `json:"version"` // semver
	Description string   `json:"description"`
	Generators  []string `json:"generators"`  // generator names this plugin contributes
	Flows       []string `json:"flows"`       // optional flow IDs this plugin contributes
	EntryPoint  string   `json:"entry_point"` // relative path to executable / wasm / main.go
	Dir         string   `json:"-"`           // absolute install dir (filled by Loader)
}

// PluginsDir returns the directory where plugins live. Honours $DOT_PLUGIN_DIR
// for testing; otherwise resolves to $XDG_DATA_HOME/dot/plugins (or
// ~/.dot/plugins on systems without XDG).
func PluginsDir() (string, error) {
	if env := os.Getenv("DOT_PLUGIN_DIR"); env != "" {
		return env, nil
	}
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "dot", "plugins"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("plugin: home dir: %w", err)
	}
	return filepath.Join(home, ".dot", "plugins"), nil
}

// List enumerates every plugin under PluginsDir(), sorted by ID. Returns an
// empty slice (not an error) when the directory does not exist.
func List() ([]*Installed, error) {
	dir, err := PluginsDir()
	if err != nil {
		return nil, err
	}
	return ListIn(dir)
}

// ListIn is like List but reads from an explicit directory; useful for tests
// and for the `--plugin-dir` CLI flag.
func ListIn(dir string) ([]*Installed, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("plugin: read %s: %w", dir, err)
	}

	var out []*Installed
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pluginRoot := filepath.Join(dir, e.Name())
		manifest := filepath.Join(pluginRoot, "plugin.json")
		data, err := os.ReadFile(manifest)
		if err != nil {
			if os.IsNotExist(err) {
				continue // not a DOT plugin; ignore
			}
			return nil, fmt.Errorf("plugin: read %s: %w", manifest, err)
		}
		var p Installed
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("plugin: parse %s: %w", manifest, err)
		}
		if p.ID == "" {
			p.ID = e.Name()
		}
		p.Dir = pluginRoot
		out = append(out, &p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
