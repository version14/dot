package versioning

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Cache wraps the on-disk store at ~/.dot/cache/generators/<name>/<version>/.
// It does NOT speak HTTP — fetching is the plugin/installer's job. The cache
// is a content store: write fully-resolved versions in, ask for them by name
// + version constraint later.
type Cache struct {
	Root string // absolute path to ~/.dot/cache/generators (or test override)
}

// NewCache constructs a Cache rooted at $XDG_CACHE_HOME/dot/generators (or
// ~/.dot/cache/generators on systems without XDG).
func NewCache() (*Cache, error) {
	root, err := defaultCacheRoot()
	if err != nil {
		return nil, err
	}
	return &Cache{Root: root}, nil
}

func defaultCacheRoot() (string, error) {
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, "dot", "generators"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("versioning: home dir: %w", err)
	}
	return filepath.Join(home, ".dot", "cache", "generators"), nil
}

// Has reports whether a specific version of name is present in the cache.
func (c *Cache) Has(name string, v Version) bool {
	_, err := os.Stat(c.versionDir(name, v))
	return err == nil
}

// Path returns the directory where name@v lives in the cache. The directory
// may or may not exist; callers use this to write into the cache or read
// from it after a Has check.
func (c *Cache) Path(name string, v Version) string {
	return c.versionDir(name, v)
}

// Available returns every cached version of name in descending semver order.
// Empty slice means "nothing cached".
func (c *Cache) Available(name string) ([]Version, error) {
	dir := filepath.Join(c.Root, name)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("versioning: read cache %s: %w", dir, err)
	}
	out := make([]Version, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		v, err := Parse(e.Name())
		if err != nil {
			continue // ignore stray dirs
		}
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Compare(out[j]) > 0
	})
	return out, nil
}

// Resolve picks the highest cached version of name that satisfies cnst.
// Returns an error containing the available versions if no match exists.
func (c *Cache) Resolve(name string, cnst Constraint) (Version, error) {
	versions, err := c.Available(name)
	if err != nil {
		return Version{}, err
	}
	for _, v := range versions {
		if cnst.Allows(v) {
			return v, nil
		}
	}
	if len(versions) == 0 {
		return Version{}, fmt.Errorf("versioning: %q not in cache", name)
	}
	return Version{}, fmt.Errorf("versioning: no cached %q matches %q (have: %s)",
		name, cnst, joinVersions(versions))
}

func (c *Cache) versionDir(name string, v Version) string {
	return filepath.Join(c.Root, name, v.String())
}

func joinVersions(vs []Version) string {
	out := ""
	for i, v := range vs {
		if i > 0 {
			out += ", "
		}
		out += v.String()
	}
	return out
}
