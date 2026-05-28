package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// versionLineRe matches the Version field in a manifest.go, e.g.:
//
//	Version:     "0.1.0",
var versionLineRe = regexp.MustCompile(`(Version:\s*)"(\d+\.\d+\.\d+)"`)

// runPatch replaces the pinned version of one package inside a generator's
// source file and bumps the manifest Version field by the appropriate
// magnitude:
//   - dependency major bump → manifest minor bump (0.1.0 → 0.2.0)
//   - dependency minor/patch bump → manifest patch bump (0.1.0 → 0.1.1)
//
// The replacement string preserves the constraint prefix from current
// (e.g. current="^4.21.0", latest="5.1.0" → new="^5.1.0").
func runPatch(generatorName, pkg, current, latest string) error {
	if generatorName == "" || pkg == "" || current == "" || latest == "" {
		return fmt.Errorf("--generator, --package, --current, and --latest are all required")
	}

	if err := patchGeneratorFile(generatorName, pkg, current, latest); err != nil {
		return err
	}

	isMajorBump := depBumpIsMajor(current, latest)
	if err := patchManifestVersion(generatorName, isMajorBump); err != nil {
		return err
	}

	return patchGeneratorDocVersion(generatorName)
}

// depBumpIsMajor returns true when latest has a higher major version than the
// base of current (stripping any constraint prefix like ^, ~, >=).
func depBumpIsMajor(current, latest string) bool {
	cMaj, _, _, ok1 := parseSemver(current)
	lMaj, _, _, ok2 := parseSemver(latest)
	return ok1 && ok2 && lMaj > cMaj
}

// patchGeneratorFile rewrites the version string for pkg inside generator.go.
// It uses a regex to tolerate alignment whitespace between the key and value
// (e.g. `"jsonwebtoken":  "^9.0.2"` vs `"express": "^4.21.0"`).
func patchGeneratorFile(generatorName, pkg, current, latest string) error {
	filePath := filepath.Join("generators", generatorName, "generator.go")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", filePath, err)
	}

	newVersion := applyPrefix(current, latest)

	// Allow one or more spaces between the colon and the value to handle
	// tab-aligned map literals. Anchor on the quoted strings to avoid
	// false positives in comments or other generators.
	pattern := fmt.Sprintf(`(%s:\s+)%s`, regexp.QuoteMeta(fmt.Sprintf("%q", pkg)), regexp.QuoteMeta(fmt.Sprintf("%q", current)))
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("build pattern for %q: %w", pkg, err)
	}

	updated := re.ReplaceAllString(string(content), `${1}`+fmt.Sprintf("%q", newVersion))
	if updated == string(content) {
		return fmt.Errorf("patch: %q: %q not found in %s", pkg, current, filePath)
	}

	if err := os.WriteFile(filePath, []byte(updated), 0644); err != nil {
		return fmt.Errorf("write %s: %w", filePath, err)
	}

	fmt.Printf("patched %s: %s %s → %s\n", filePath, pkg, current, newVersion)
	return nil
}

// patchManifestVersion bumps the Version field in manifest.go.
// majorBump=true increments the minor segment (0.1.0 → 0.2.0, resetting patch).
// majorBump=false increments the patch segment (0.1.0 → 0.1.1).
func patchManifestVersion(generatorName string, majorBump bool) error {
	filePath := filepath.Join("generators", generatorName, "manifest.go")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", filePath, err)
	}

	var bumpErr error
	updated := versionLineRe.ReplaceAllStringFunc(string(content), func(match string) string {
		sub := versionLineRe.FindStringSubmatch(match)
		if sub == nil {
			return match
		}
		var bumped string
		var err error
		if majorBump {
			bumped, err = bumpMinor(sub[2])
		} else {
			bumped, err = bumpPatch(sub[2])
		}
		if err != nil {
			bumpErr = err
			return match
		}
		fmt.Printf("bumped manifest version %s → %s in %s\n", sub[2], bumped, filePath)
		return sub[1] + `"` + bumped + `"`
	})
	if bumpErr != nil {
		return fmt.Errorf("bump manifest version in %s: %w", filePath, bumpErr)
	}

	return os.WriteFile(filePath, []byte(updated), 0644)
}

// patchGeneratorDocVersion updates the generator doc version in
// docs/contributor/generators/<name>.md to match the manifest version.
func patchGeneratorDocVersion(generatorName string) error {
	manifestPath := filepath.Join("generators", generatorName, "manifest.go")
	manifestContent, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", manifestPath, err)
	}

	manifestMatch := versionLineRe.FindStringSubmatch(string(manifestContent))
	if manifestMatch == nil {
		return fmt.Errorf("manifest version not found in %s", manifestPath)
	}
	manifestVersion := manifestMatch[2]

	docPath := filepath.Join("docs", "contributor", "generators", generatorName+".md")
	docContent, err := os.ReadFile(docPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", docPath, err)
	}

	versionRowRe := regexp.MustCompile(`(?m)^\|\s*Version\s*\|\s*` + "`" + `([^` + "`" + `]+)` + "`" + `\s*\|\s*$`)
	updated := versionRowRe.ReplaceAllString(string(docContent), "| Version | `"+manifestVersion+"` |")
	if updated == string(docContent) {
		return fmt.Errorf("doc version row not found in %s", docPath)
	}

	if err := os.WriteFile(docPath, []byte(updated), 0644); err != nil {
		return fmt.Errorf("write %s: %w", docPath, err)
	}
	fmt.Printf("updated doc version to %s in %s\n", manifestVersion, docPath)
	return nil
}

// bumpMinor increments the minor segment and resets patch (e.g. "0.1.3" → "0.2.0").
func bumpMinor(version string) (string, error) {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("unexpected version format %q", version)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("non-numeric minor in %q: %w", version, err)
	}
	return parts[0] + "." + strconv.Itoa(minor+1) + ".0", nil
}

// bumpPatch increments the patch segment (e.g. "0.1.0" → "0.1.1").
func bumpPatch(version string) (string, error) {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("unexpected version format %q", version)
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("non-numeric patch in %q: %w", version, err)
	}
	return parts[0] + "." + parts[1] + "." + strconv.Itoa(patch+1), nil
}

// applyPrefix extracts any leading constraint characters from current
// (^, ~, >=, etc.) and prepends them to latest.
func applyPrefix(current, latest string) string {
	prefix := ""
	for _, c := range current {
		if c == '^' || c == '~' || c == '>' || c == '=' || c == '<' {
			prefix += string(c)
		} else {
			break
		}
	}
	return prefix + latest
}
