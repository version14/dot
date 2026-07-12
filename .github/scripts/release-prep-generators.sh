#!/usr/bin/env bash
#
# Rule 3 (ADR-0002): bump every Generator whose Fingerprint has moved since the
# last release, once, from the final state.
#
# Run this on main immediately before cutting a tag:
#
#     make release-prep      # then review the diff, commit, and tag
#
# Why here and not in the dependency PR
# ─────────────────────────────────────
# A Catalog Pin change alters what several Generators scaffold, so their Manifest
# Versions must move. But a version bump computed inside a PR is a *relative*
# operation — "0.8.0 → 0.9.0" — worked out against main at PR-creation time and
# applied at merge time, when main has moved. With several dependency PRs open,
# the version you end up with is a function of merge order. That is the bug that
# left generator versions wrong, and it is fatal for long-lived Migration PRs: a
# React Migration open for two months would collide with every weekly Rollup that
# touches the same manifest.
#
# Deriving the bump here, once, from the final state makes merge order provably
# irrelevant. Five PRs in any sequence produce the identical manifest.
set -euo pipefail

BASE_REF="${1:-$(git describe --tags --abbrev=0 2>/dev/null || echo "")}"
if [[ -z "$BASE_REF" ]]; then
  echo "release-prep: no previous tag found; pass a base ref explicitly." >&2
  exit 1
fi

BUMP_TYPE="${BUMP_TYPE:-patch}"   # BUMP_TYPE=minor for a behaviour-changing release
WORKTREE=$(mktemp -d)
BASE_FP=$(mktemp)
HEAD_FP=$(mktemp)
trap 'git worktree remove --force "$WORKTREE" 2>/dev/null || true; rm -f "$BASE_FP" "$HEAD_FP"' EXIT

echo "release-prep: comparing against $BASE_REF"

# Fingerprints depend on the generators compiled into the binary, so the base
# fingerprints must come from a binary built from the BASE generators. Overlay
# this checkout's fingerprinting code onto the base tree — older tags may not
# have `gen-fingerprint` at all.
git worktree add --detach "$WORKTREE" "$BASE_REF" >/dev/null 2>&1
mkdir -p "$WORKTREE/internal/fingerprint"
cp internal/fingerprint/fingerprint.go "$WORKTREE/internal/fingerprint/"
cp internal/flow/scripted.go           "$WORKTREE/internal/flow/"
cp internal/generator/executor.go      "$WORKTREE/internal/generator/"
cp internal/cli/gen_fingerprint.go     "$WORKTREE/internal/cli/"
cp internal/cli/gen_check.go           "$WORKTREE/internal/cli/"
cp internal/cli/runner.go              "$WORKTREE/internal/cli/"
cp internal/cli/command.go             "$WORKTREE/internal/cli/"

( cd "$WORKTREE" && go build -o "$WORKTREE/dot-base" ./cmd/dot )
"$WORKTREE/dot-base" gen-fingerprint --json > "$BASE_FP"

go build -o bin/dot ./cmd/dot
./bin/dot gen-fingerprint --json > "$HEAD_FP"

MOVED=$(python3 - "$BASE_FP" "$HEAD_FP" <<'PY'
import json, sys
base = json.load(open(sys.argv[1]))
head = json.load(open(sys.argv[2]))
for name in sorted(head):
    if name in base and base[name] != head[name]:
        print(name)
PY
)

if [[ -z "$MOVED" ]]; then
  echo "release-prep: no generator's output changed since $BASE_REF — nothing to bump."
  exit 0
fi

echo ""
echo "These generators scaffold something different than they did at $BASE_REF:"
while IFS= read -r name; do
  [[ -z "$name" ]] && continue

  # Already bumped by a human in this cycle? Then leave it alone — bumping again
  # would compound, which is precisely what we are here to prevent.
  head_ver=$(grep -oE 'Version:\s*"[0-9]+\.[0-9]+\.[0-9]+"' "generators/$name/manifest.go" | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
  base_ver=$(git show "$BASE_REF:generators/$name/manifest.go" 2>/dev/null \
              | grep -oE 'Version:\s*"[0-9]+\.[0-9]+\.[0-9]+"' | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || echo "")

  if [[ -n "$base_ver" && "$head_ver" != "$base_ver" ]]; then
    echo "  $name — already bumped ($base_ver → $head_ver), leaving it"
    continue
  fi
  ./bin/dot gen-bump --name "$name" --bump "$BUMP_TYPE"
done <<<"$MOVED"

echo ""
echo "release-prep: done. Review the diff, commit, then tag."
