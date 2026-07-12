// Package fingerprint identifies what a Generator scaffolds, so that a Manifest
// Version moves when — and only when — the Generator's output actually changes.
//
// A Contribution is what one Generator adds to a Project when it runs: the files
// it creates and the edits it makes to files other Generators own. A Fingerprint
// is a hash of a Generator's Contributions across every scenario it is exercised
// in. Two Generators with the same Fingerprint scaffold the same thing.
//
// This answers the only question that matters when deciding whether a version
// must be bumped: *did what this Generator scaffolds actually change?*
//
//   - Reformatting a Generator's source, or fixing a comment → Fingerprint
//     unchanged. No bump demanded.
//   - Editing a template → Fingerprint changes. Bump demanded (CI rule 1).
//   - Moving a Catalog Pin the Generator names → Fingerprint changes, even
//     though no Generator source was touched. Bump demanded at release (rule 3).
//
// # Why a Contribution is a diff, not a file hash
//
// Generators merge into shared files — above all package.json. If a Contribution
// were "the hash of every file after this generator ran", then bumping the react
// Pin would move the Fingerprint of *every generator that touches package.json*,
// because they'd all see a package.json containing the new react version. Only
// two generators name react.
//
// That over-triggering is not a cosmetic flaw. It is precisely the bug this work
// exists to remove: commit d53c302 bumped 15 generator manifests for a change to
// one generator, because the old bot could not tell whose change it was.
//
// So a Contribution is the *edit itself*, computed against the state as it stood
// immediately before the generator ran:
//
//   - a file created  → its content
//   - a JSON file merged into → the key paths this generator added or changed,
//     with their values (react_app contributes `dependencies.react=^19.2.7`;
//     ark_ui contributes `dependencies.@ark-ui/react=…` and is unaffected by react)
//   - a file appended to → the appended bytes
//   - a file removed  → the removal
//
// A Contribution therefore depends only on what the generator did, never on what
// ran before it.
package fingerprint

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/version14/dot/internal/state"
)

// Collector implements generator.Observer, capturing each Generator's
// Contribution as the pipeline runs.
//
// Contributions cannot be recovered after the fact: the finished state records
// only *that* a generator touched a file, never what it put there. Bracketing
// each call is the only way to attribute the edit.
type Collector struct {
	contributions map[string][]string
	before        map[string][]byte
}

func NewCollector() *Collector {
	return &Collector{contributions: map[string][]string{}}
}

func (c *Collector) BeforeGenerate(_ string, root *state.VirtualProjectState) {
	c.before = snapshot(root)
}

func (c *Collector) AfterGenerate(name string, root *state.VirtualProjectState) {
	after := snapshot(root)
	before := c.before
	c.before = nil

	var edits []string

	for path, now := range after {
		prev := before[path] // nil when this generator created the file
		if bytes.Equal(prev, now) {
			continue
		}
		// Deliberately no create/modify distinction. Whether a generator creates
		// package.json or merges into one another generator already created is an
		// accident of topological order, not of the generator's behaviour — and if
		// it were encoded here, inserting a new generator ahead of an existing one
		// would move that generator's Fingerprint and demand a version bump for a
		// generator that did not change. A create is simply a diff against nothing.
		edits = append(edits, path+" "+digest([]byte(deltaOf(prev, now))))
	}
	for path := range before {
		if _, still := after[path]; !still {
			edits = append(edits, "remove "+path)
		}
	}

	if len(edits) == 0 {
		return
	}
	sort.Strings(edits)
	c.contributions[name] = append(c.contributions[name], digest([]byte(strings.Join(edits, "\n"))))
}

// Fingerprints returns one Fingerprint per Generator seen.
//
// A Generator's Contributions are deduplicated and sorted before hashing, so a
// Fingerprint does not depend on how many fixtures exercised the generator, nor
// on the order they ran in. Adding a fixture that produces an already-seen
// Contribution leaves the Fingerprint untouched.
func (c *Collector) Fingerprints() map[string]string {
	out := make(map[string]string, len(c.contributions))
	for name, contribs := range c.contributions {
		seen := map[string]bool{}
		uniq := make([]string, 0, len(contribs))
		for _, h := range contribs {
			if !seen[h] {
				seen[h] = true
				uniq = append(uniq, h)
			}
		}
		sort.Strings(uniq)
		out[name] = digest([]byte(strings.Join(uniq, "\n")))
	}
	return out
}

// deltaOf describes how now differs from prev, in a form that does not depend on
// prev's own content. prev is nil when the generator created the file.
func deltaOf(prev, now []byte) string {
	// JSON: report the key paths this generator added or changed. This is what
	// keeps a react Pin bump off the Fingerprint of every generator that merely
	// merges something else into the same package.json.
	var nowDoc any
	if json.Unmarshal(now, &nowDoc) == nil {
		var prevDoc any
		if len(prev) > 0 && json.Unmarshal(prev, &prevDoc) != nil {
			// It was not JSON before and is now: a wholesale replacement.
			return "rewrite\n" + string(now)
		}
		var changes []string
		diffJSON("", prevDoc, nowDoc, &changes)
		sort.Strings(changes)
		return "json\n" + strings.Join(changes, "\n")
	}

	// Append: the generator added to the end (AppendFile). The edit is the tail.
	// A created file is an append to nothing, which falls out of this naturally.
	if len(now) >= len(prev) && string(now[:len(prev)]) == string(prev) {
		return "append\n" + string(now[len(prev):])
	}

	// Anything else is a wholesale rewrite; the edit is the new content.
	return "rewrite\n" + string(now)
}

// diffJSON walks two decoded JSON values and records every leaf that was added
// or changed, as "path=value". A nil prev means the file did not exist, which is
// diffed as though it were an empty object.
func diffJSON(path string, prev, now any, out *[]string) {
	prevMap, prevIsMap := prev.(map[string]any)
	nowMap, nowIsMap := now.(map[string]any)

	if nowIsMap && prev == nil {
		prevMap, prevIsMap = map[string]any{}, true
	}

	if prevIsMap && nowIsMap {
		for key, nowVal := range nowMap {
			child := key
			if path != "" {
				child = path + "." + key
			}
			diffJSON(child, prevMap[key], nowVal, out)
		}
		for key := range prevMap {
			if _, still := nowMap[key]; !still {
				child := key
				if path != "" {
					child = path + "." + key
				}
				*out = append(*out, child+"=<removed>")
			}
		}
		return
	}

	if !sameValue(prev, now) {
		*out = append(*out, fmt.Sprintf("%s=%v", path, now))
	}
}

func sameValue(a, b any) bool {
	ab, errA := json.Marshal(a)
	bb, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(ab) == string(bb)
}

func snapshot(root *state.VirtualProjectState) map[string][]byte {
	out := make(map[string][]byte, len(root.Files))
	for path, node := range root.Files {
		out[path] = node.Content
	}
	return out
}

func digest(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])[:16]
}
