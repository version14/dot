package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DepInfo is the result of checking one package against its registry.
type DepInfo struct {
	Latest     string
	Deprecated bool
	Notice     string
}

// RegistryChecker queries one package registry for version and deprecation info.
type RegistryChecker interface {
	Ecosystem() Ecosystem
	Check(pkg, currentVersion string) (DepInfo, error)
}

// checkerFor returns the right RegistryChecker based on which file the
// generator wrote to.
func checkerFor(filename string) (RegistryChecker, bool) {
	client := &http.Client{Timeout: 15 * time.Second}
	switch filename {
	case "package.json":
		return &npmChecker{base: baseChecker{client}}, true
	case "go.mod":
		return &goProxyChecker{base: baseChecker{client}}, true
	case "Cargo.toml":
		return &cargoChecker{base: baseChecker{client}}, true
	case "pom.xml":
		return &mavenChecker{base: baseChecker{client}}, true
	default:
		return nil, false
	}
}

// isOutdated returns true when latestVersion is strictly newer than the base
// of currentConstraint (stripping leading ^~>= characters).
func isOutdated(currentConstraint, latestVersion string) bool {
	cMaj, cMin, cPatch, ok1 := parseSemver(currentConstraint)
	lMaj, lMin, lPatch, ok2 := parseSemver(latestVersion)
	if !ok1 || !ok2 {
		return false
	}
	if lMaj != cMaj {
		return lMaj > cMaj
	}
	if lMin != cMin {
		return lMin > cMin
	}
	return lPatch > cPatch
}

func stripConstraintPrefix(v string) string {
	return strings.TrimLeft(v, "^~>=<v")
}

func parseSemver(v string) (major, minor, patch int, ok bool) {
	v = stripConstraintPrefix(v)
	parts := strings.SplitN(v, ".", 3)
	if len(parts) < 2 {
		return 0, 0, 0, false
	}
	maj, err1 := strconv.Atoi(parts[0])
	min, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, 0, 0, false
	}
	pat := 0
	if len(parts) == 3 {
		patchStr := strings.SplitN(parts[2], "-", 2)[0]
		patchStr = strings.SplitN(patchStr, "+", 2)[0]
		pat, _ = strconv.Atoi(patchStr)
	}
	return maj, min, pat, true
}

// --- shared HTTP helper ---

type baseChecker struct {
	client *http.Client
}

// get fetches url, checks the status code, and returns the response body.
// errPrefix and pkg are used only to annotate errors.
func (b *baseChecker) get(url, errPrefix, pkg string) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "dot-dep-checker/1 (github.com/version14/dot)")
	req.Header.Set("Accept", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", errPrefix, pkg, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
		return nil, fmt.Errorf("%s: %s: not found (HTTP %d)", errPrefix, pkg, resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s: HTTP %d", errPrefix, pkg, resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("%s: %s: read: %w", errPrefix, pkg, err)
	}
	return body, nil
}

// --- npm ---

type npmChecker struct{ base baseChecker }

func (c *npmChecker) Ecosystem() Ecosystem { return EcosystemNPM }

func (c *npmChecker) Check(pkg, _ string) (DepInfo, error) {
	body, err := c.base.get("https://registry.npmjs.org/"+pkg+"/latest", "npm", pkg)
	if err != nil {
		return DepInfo{}, err
	}

	var manifest struct {
		Version    string `json:"version"`
		Deprecated string `json:"deprecated"`
	}
	if err := json.Unmarshal(body, &manifest); err != nil {
		return DepInfo{}, fmt.Errorf("npm: %s: parse: %w", pkg, err)
	}

	return DepInfo{
		Latest:     manifest.Version,
		Deprecated: manifest.Deprecated != "",
		Notice:     manifest.Deprecated,
	}, nil
}

// --- Go proxy ---

type goProxyChecker struct{ base baseChecker }

func (c *goProxyChecker) Ecosystem() Ecosystem { return EcosystemGo }

func (c *goProxyChecker) Check(pkg, _ string) (DepInfo, error) {
	// Go module paths with uppercase letters need !-escaping per GOPROXY
	// protocol, but most modules are lowercase — handle the common case.
	body, err := c.base.get("https://proxy.golang.org/"+pkg+"/@latest", "goproxy", pkg)
	if err != nil {
		return DepInfo{}, err
	}

	var info struct {
		Version string `json:"Version"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return DepInfo{}, fmt.Errorf("goproxy: %s: parse: %w", pkg, err)
	}

	// Go modules do not expose a deprecation field via the proxy API;
	// deprecation is expressed via a `Deprecated:` comment in go.mod itself.
	return DepInfo{Latest: info.Version}, nil
}

// --- Cargo (crates.io) ---

type cargoChecker struct{ base baseChecker }

func (c *cargoChecker) Ecosystem() Ecosystem { return EcosystemCargo }

func (c *cargoChecker) Check(pkg, _ string) (DepInfo, error) {
	body, err := c.base.get("https://crates.io/api/v1/crates/"+pkg, "crates.io", pkg)
	if err != nil {
		return DepInfo{}, err
	}

	var result struct {
		Crate struct {
			NewestVersion string `json:"newest_version"`
		} `json:"crate"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return DepInfo{}, fmt.Errorf("crates.io: %s: parse: %w", pkg, err)
	}

	return DepInfo{Latest: result.Crate.NewestVersion}, nil
}

// --- Maven Central ---

type mavenChecker struct{ base baseChecker }

func (c *mavenChecker) Ecosystem() Ecosystem { return EcosystemMaven }

// Check expects pkg in "groupId:artifactId" format (e.g. "org.springframework:spring-core").
func (c *mavenChecker) Check(pkg, _ string) (DepInfo, error) {
	parts := strings.SplitN(pkg, ":", 2)
	if len(parts) != 2 {
		return DepInfo{}, fmt.Errorf("maven: %q must be in groupId:artifactId format", pkg)
	}
	groupID, artifactID := parts[0], parts[1]

	url := fmt.Sprintf(
		"https://search.maven.org/solrsearch/select?q=g:%s+AND+a:%s&core=gav&rows=1&wt=json",
		groupID, artifactID,
	)
	body, err := c.base.get(url, "maven", pkg)
	if err != nil {
		return DepInfo{}, err
	}

	var result struct {
		Response struct {
			Docs []struct {
				Version string `json:"v"`
			} `json:"docs"`
		} `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return DepInfo{}, fmt.Errorf("maven: %s: parse: %w", pkg, err)
	}

	if len(result.Response.Docs) == 0 {
		return DepInfo{}, fmt.Errorf("maven: %s: not found", pkg)
	}

	return DepInfo{Latest: result.Response.Docs[0].Version}, nil
}
