package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/version14/dot/internal/deps"
)

// Info is what a registry knows about a package right now.
type Info struct {
	Latest     string
	Deprecated bool
	Notice     string
}

// Registry queries one package Ecosystem.
//
// Ecosystems differ in ways that cannot be papered over, which is why this is an
// interface rather than a switch on a filename:
//
//   - Version notation. npm writes "^5.61.3"; Go writes "v1.9.0"; Cargo writes
//     "1.0.4". Render owns that.
//   - Deprecation. npm exposes a `deprecated` field. The Go proxy does not —
//     deprecation lives in a `Deprecated:` comment in go.mod.
//   - Patchability. A Go *major* bump changes the import path
//     (github.com/x/y → github.com/x/y/v2). No amount of string rewriting
//     handles that, so for Go it can never be a Rollup, only a Migration.
//
// The old checker inferred the ecosystem from which file a generator happened to
// write, which is why `pom.xml` sat in a lookup table with no extractor behind it.
// A Pin declares its Ecosystem; nothing is guessed.
type Registry interface {
	Ecosystem() deps.Ecosystem

	// Latest reports the newest published version of pkg, and whether the
	// registry marks it deprecated.
	Latest(pkg string) (Info, error)

	// Render formats latest in the notation of pin, preserving pin's constraint
	// prefix ("^5.61.3" + "5.62.0" → "^5.62.0").
	Render(pin, latest string) string
}

// registryFor returns the Registry for an Ecosystem.
//
// Adding an ecosystem is one implementation and one line here. Cargo, Maven and
// Go are absent on purpose: no generator scaffolds a Cargo.toml, pom.xml or
// go.mod today, and the previous tool carried unreachable code for all three
// (`pom.xml` was even routed to an extractor that did not exist). They come back
// when a generator needs them, not before.
func registryFor(eco deps.Ecosystem) (Registry, bool) {
	switch eco {
	case deps.EcosystemNPM:
		return &npmRegistry{client: &http.Client{Timeout: 15 * time.Second}}, true
	default:
		return nil, false
	}
}

// --- npm ---

type npmRegistry struct{ client *http.Client }

func (r *npmRegistry) Ecosystem() deps.Ecosystem { return deps.EcosystemNPM }

func (r *npmRegistry) Latest(pkg string) (Info, error) {
	// The `latest` dist-tag is what `npm install <pkg>` resolves to, and it
	// excludes prereleases — which is exactly the version we want to compare a
	// Pin against.
	url := "https://registry.npmjs.org/" + pkg + "/latest"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Info{}, fmt.Errorf("npm: %s: %w", pkg, err)
	}
	req.Header.Set("User-Agent", "dot-dep-checker/2 (github.com/version14/dot)")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return Info{}, fmt.Errorf("npm: %s: %w", pkg, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Info{}, fmt.Errorf("npm: %s: HTTP %d", pkg, resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Info{}, fmt.Errorf("npm: %s: read: %w", pkg, err)
	}

	var manifest struct {
		Version    string `json:"version"`
		Deprecated string `json:"deprecated"`
	}
	if err := json.Unmarshal(body, &manifest); err != nil {
		return Info{}, fmt.Errorf("npm: %s: parse: %w", pkg, err)
	}

	return Info{
		Latest:     manifest.Version,
		Deprecated: manifest.Deprecated != "",
		Notice:     manifest.Deprecated,
	}, nil
}

// Render preserves the Pin's constraint prefix. The prefix is the compatibility
// promise; rewriting "^19.2.6" as a bare "20.0.0" would silently discard it.
func (r *npmRegistry) Render(pin, latest string) string {
	prefix := pin[:len(pin)-len(strings.TrimLeft(pin, "^~>=<"))]
	return prefix + latest
}
