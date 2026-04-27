// Package versioning implements the small subset of semver that DOT needs:
// parsing X.Y.Z(-prerelease), comparing two versions, and matching a version
// against a constraint string like "^1.2.3", "~0.3.4", ">=2.0.0", or "1.2.3".
//
// We intentionally avoid pulling in golang.org/x/mod or hashicorp/go-version —
// the constraint vocabulary here is small and stable, and a hand-rolled
// implementation keeps the dependency surface light.
package versioning

import (
	"fmt"
	"strconv"
	"strings"
)

// Version is a parsed X.Y.Z(-prerelease) semver triple. Build metadata
// (after `+`) is dropped at parse time per semver §10.
type Version struct {
	Major, Minor, Patch int
	Prerelease          string // empty when none
	Raw                 string // original input for round-tripping
}

// Parse reads a semver string. Leading "v" is allowed (e.g. "v1.2.3").
// Returns an error if the structure isn't <major>.<minor>.<patch>[-pre].
func Parse(s string) (Version, error) {
	original := s
	s = strings.TrimPrefix(strings.TrimSpace(s), "v")
	if s == "" {
		return Version{}, fmt.Errorf("versioning: empty version")
	}

	if i := strings.IndexByte(s, '+'); i >= 0 {
		s = s[:i]
	}

	pre := ""
	if i := strings.IndexByte(s, '-'); i >= 0 {
		pre = s[i+1:]
		s = s[:i]
	}

	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("versioning: %q is not <major>.<minor>.<patch>", original)
	}
	maj, err := strconv.Atoi(parts[0])
	if err != nil || maj < 0 {
		return Version{}, fmt.Errorf("versioning: invalid major in %q", original)
	}
	min, err := strconv.Atoi(parts[1])
	if err != nil || min < 0 {
		return Version{}, fmt.Errorf("versioning: invalid minor in %q", original)
	}
	pat, err := strconv.Atoi(parts[2])
	if err != nil || pat < 0 {
		return Version{}, fmt.Errorf("versioning: invalid patch in %q", original)
	}
	return Version{Major: maj, Minor: min, Patch: pat, Prerelease: pre, Raw: original}, nil
}

// MustParse panics on error. Use only with literal known-good input.
func MustParse(s string) Version {
	v, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

// String renders the canonical semver form (no build metadata, ever).
func (v Version) String() string {
	if v.Prerelease == "" {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.Prerelease)
}

// Compare returns -1 / 0 / +1.
//
// Stable versions outrank prereleases of the same X.Y.Z (semver §11.4):
// 1.0.0 > 1.0.0-rc1.
func (v Version) Compare(other Version) int {
	if c := cmpInt(v.Major, other.Major); c != 0 {
		return c
	}
	if c := cmpInt(v.Minor, other.Minor); c != 0 {
		return c
	}
	if c := cmpInt(v.Patch, other.Patch); c != 0 {
		return c
	}
	switch {
	case v.Prerelease == "" && other.Prerelease == "":
		return 0
	case v.Prerelease == "":
		return 1
	case other.Prerelease == "":
		return -1
	default:
		return strings.Compare(v.Prerelease, other.Prerelease)
	}
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// ── Constraints ───────────────────────────────────────────────────────────

// Constraint describes a version range. Supported syntaxes:
//
//	"1.2.3"     — exact match
//	"=1.2.3"    — exact match (explicit)
//	">=1.2.3"   — at least
//	">1.2.3"    — strictly greater
//	"<=1.2.3"   — at most
//	"<1.2.3"    — strictly less
//	"~1.2.3"    — any 1.2.x with patch ≥ 3
//	"^1.2.3"    — any 1.x.x with version ≥ 1.2.3 (major-locked)
//	"^0.2.3"    — any 0.2.x with patch ≥ 3 (when major=0, minor is locked)
//	""          — accept any version
type Constraint struct {
	raw string
	op  string  // "", "=", ">=", ">", "<=", "<", "~", "^"
	ver Version // anchor (omitted when raw == "")
}

// ParseConstraint parses a single constraint expression.
// Empty string is a valid wildcard match (Allows always returns true).
func ParseConstraint(s string) (Constraint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Constraint{raw: ""}, nil
	}

	op := ""
	rest := s
	for _, candidate := range []string{">=", "<=", "^", "~", ">", "<", "="} {
		if strings.HasPrefix(rest, candidate) {
			op = candidate
			rest = strings.TrimSpace(strings.TrimPrefix(rest, candidate))
			break
		}
	}

	v, err := Parse(rest)
	if err != nil {
		return Constraint{}, fmt.Errorf("versioning: constraint %q: %w", s, err)
	}
	return Constraint{raw: s, op: op, ver: v}, nil
}

// MustParseConstraint panics on error. Use only with literal known-good input.
func MustParseConstraint(s string) Constraint {
	c, err := ParseConstraint(s)
	if err != nil {
		panic(err)
	}
	return c
}

// String returns the original constraint text.
func (c Constraint) String() string { return c.raw }

// Allows reports whether v satisfies this constraint.
func (c Constraint) Allows(v Version) bool {
	if c.raw == "" {
		return true
	}
	switch c.op {
	case "", "=":
		return v.Compare(c.ver) == 0
	case ">":
		return v.Compare(c.ver) > 0
	case ">=":
		return v.Compare(c.ver) >= 0
	case "<":
		return v.Compare(c.ver) < 0
	case "<=":
		return v.Compare(c.ver) <= 0
	case "~":
		if v.Compare(c.ver) < 0 {
			return false
		}
		return v.Major == c.ver.Major && v.Minor == c.ver.Minor
	case "^":
		if v.Compare(c.ver) < 0 {
			return false
		}
		if c.ver.Major > 0 {
			return v.Major == c.ver.Major
		}
		return v.Major == 0 && v.Minor == c.ver.Minor
	}
	return false
}
