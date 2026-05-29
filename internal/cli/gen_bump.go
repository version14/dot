package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var manifestVersionRe = regexp.MustCompile(`(Version:\s*)"(\d+\.\d+\.\d+)"`)
var docVersionRowRe = regexp.MustCompile("(?m)^\\|\\s*Version\\s*\\|\\s*`([^`]+)`\\s*\\|\\s*$")

func runGenBump(args []string) int {
	fs := flag.NewFlagSet("gen-bump", flag.ContinueOnError)
	all := fs.Bool("all", false, "bump all generators")
	name := fs.String("name", "", "bump a specific generator by name (comma-separated for multiple)")
	bump := fs.String("bump", "patch", "bump type: patch, minor, or major")
	set := fs.String("set", "", "set an explicit version instead of bumping (e.g. 1.2.3)")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if !*all && *name == "" {
		fmt.Fprintln(os.Stderr, "dot gen-bump: --all or --name <name> is required")
		fs.Usage()
		return 2
	}
	if *all && *name != "" {
		fmt.Fprintln(os.Stderr, "dot gen-bump: --all and --name are mutually exclusive")
		return 2
	}
	if *set != "" && !isValidVersion(*set) {
		fmt.Fprintf(os.Stderr, "dot gen-bump: --set %q is not a valid semver (expected MAJOR.MINOR.PATCH)\n", *set)
		return 2
	}

	var names []string
	if *name != "" {
		for _, n := range strings.Split(*name, ",") {
			n = strings.TrimSpace(n)
			if n != "" {
				names = append(names, n)
			}
		}
	} else {
		found, err := listGeneratorDirs("generators")
		if err != nil {
			fmt.Fprintln(os.Stderr, "dot gen-bump:", err)
			return 1
		}
		names = found
	}

	if len(names) == 0 {
		fmt.Fprintln(os.Stderr, "dot gen-bump: no generators found")
		return 1
	}

	label := *bump
	if *set != "" {
		label = *set
	}
	PrintHeading(fmt.Sprintf("gen-bump (%s)", label))

	anyFailed := false
	for _, n := range names {
		var err error
		if *set != "" {
			err = setGenVersion(n, *set)
		} else {
			err = bumpGenVersion(n, *bump)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "  error: %s: %v\n", n, err)
			anyFailed = true
		}
	}

	if anyFailed {
		return 1
	}
	return 0
}

func listGeneratorDirs(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", root, err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, e.Name(), "manifest.go")); err == nil {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func bumpGenVersion(name, bumpType string) error {
	manifestPath := filepath.Join("generators", name, "manifest.go")
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", manifestPath, err)
	}

	var oldVer, newVer string
	var bumpErr error
	updated := manifestVersionRe.ReplaceAllStringFunc(string(content), func(match string) string {
		if bumpErr != nil {
			return match
		}
		sub := manifestVersionRe.FindStringSubmatch(match)
		if sub == nil {
			return match
		}
		oldVer = sub[2]
		var bumped string
		switch bumpType {
		case "major":
			bumped, bumpErr = bumpMajorVer(sub[2])
		case "minor":
			bumped, bumpErr = bumpMinorVer(sub[2])
		default:
			bumped, bumpErr = bumpPatchVer(sub[2])
		}
		if bumpErr != nil {
			return match
		}
		newVer = bumped
		return sub[1] + `"` + bumped + `"`
	})
	if bumpErr != nil {
		return bumpErr
	}
	if newVer == "" {
		return fmt.Errorf("version field not found in %s", manifestPath)
	}
	fmt.Printf("  %-40s %s → %s\n", name, oldVer, newVer)
	if err := os.WriteFile(manifestPath, []byte(updated), 0644); err != nil {
		return err
	}
	syncGenDocVersion(name, newVer)
	return nil
}

func setGenVersion(name, version string) error {
	manifestPath := filepath.Join("generators", name, "manifest.go")
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", manifestPath, err)
	}

	var oldVer string
	updated := manifestVersionRe.ReplaceAllStringFunc(string(content), func(match string) string {
		sub := manifestVersionRe.FindStringSubmatch(match)
		if sub == nil {
			return match
		}
		oldVer = sub[2]
		return sub[1] + `"` + version + `"`
	})
	if oldVer == "" {
		return fmt.Errorf("version field not found in %s", manifestPath)
	}
	fmt.Printf("  %-40s %s → %s\n", name, oldVer, version)
	if err := os.WriteFile(manifestPath, []byte(updated), 0644); err != nil {
		return err
	}
	syncGenDocVersion(name, version)
	return nil
}

// syncGenDocVersion updates the Version row in the generator's doc page if it exists.
func syncGenDocVersion(name, version string) {
	docPath := filepath.Join("docs", "contributor", "generators", name+".md")
	content, err := os.ReadFile(docPath)
	if err != nil {
		return
	}
	updated := docVersionRowRe.ReplaceAllString(string(content), "| Version | `"+version+"` |")
	if updated == string(content) {
		return
	}
	_ = os.WriteFile(docPath, []byte(updated), 0644)
}

func isValidVersion(v string) bool {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return false
	}
	for _, p := range parts {
		if _, err := strconv.Atoi(p); err != nil {
			return false
		}
	}
	return true
}

func bumpMajorVer(v string) (string, error) {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("unexpected version format %q", v)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("non-numeric major in %q: %w", v, err)
	}
	return strconv.Itoa(major+1) + ".0.0", nil
}

func bumpMinorVer(v string) (string, error) {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("unexpected version format %q", v)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("non-numeric minor in %q: %w", v, err)
	}
	return parts[0] + "." + strconv.Itoa(minor+1) + ".0", nil
}

func bumpPatchVer(v string) (string, error) {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("unexpected version format %q", v)
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("non-numeric patch in %q: %w", v, err)
	}
	return parts[0] + "." + parts[1] + "." + strconv.Itoa(patch+1), nil
}
