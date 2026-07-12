#!/usr/bin/env bash
#
# Turns a dep-report.json into pull requests and issues.
#
# Three vehicles, because closing Pin drift is not one kind of work (see
# CONTEXT.md and ADR-0002):
#
#   Rollup     — every update that satisfies its Pin's existing constraint, in ONE
#                batched PR. Bot-owned: rebuilt from scratch on every run. Excludes
#                any package that has an open dep PR of its own.
#   Migration  — an update that BREAKS its Pin's constraint. One PR per package.
#                Human-owned from the moment it is raised; the bot never touches
#                that branch again.
#   Deprecation— an issue. Never a PR. A deprecation is not fixed by a version
#                bump: the replacement is a different package, or none.
#
# What this script does NOT do, deliberately:
#
#   - It does not touch generators/. It patches internal/deps/npm.go and nothing
#     else. The old script ran `git add generators/`, which swept whatever was
#     dirty in the working tree into the current commit — that is how 997e489,
#     titled "bump @clerk/clerk-react", came to contain 13 files of unrelated
#     packages and no clerk changes.
#   - It does not bump manifest versions. Those are derived at release from
#     Fingerprints (`dot gen-check --bumped`), so merge order cannot corrupt them.
#     A Migration PR that sits open for months therefore rebases cleanly: it only
#     ever touches its own line of the Catalog.
#   - It does not run test-flows. The PR's own CI does that.
#
# Requires: gh, jq, dep-report.json, GH_TOKEN, GITHUB_REPOSITORY.
set -euo pipefail

REPORT="${1:-dep-report.json}"
ROLLUP_BRANCH="deps/npm/rollup"

if [[ ! -f "$REPORT" ]]; then
  echo "dep-bot: $REPORT not found" >&2
  exit 1
fi

git config user.name  "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

# Every dep branch with an open PR. A package with one is spoken for: either a
# human is doing its Migration, or a human ejected it from the Rollup because it
# broke the templates. Either way the bot leaves it alone.
# (Newline-delimited rather than an array: `mapfile` is bash 4+, and this script
# should still run on a stock macOS bash if someone wants to try it by hand.)
OPEN_BRANCHES=$(
  gh pr list --state open --limit 200 --json headRefName \
    --jq '.[] | select(.headRefName | startswith("deps/")) | .headRefName' \
    --repo "$GITHUB_REPOSITORY"
)
has_open_pr() {
  printf '%s\n' "$OPEN_BRANCHES" | grep -qxF "$1"
}

branch_for() { echo "deps/npm/${1//\//-}"; }   # @scope/pkg -> deps/npm/@scope-pkg

# ─────────────────────────────────────────────────────────────────────────────
# Rollup
# ─────────────────────────────────────────────────────────────────────────────
build_rollup() {
  local entries
  entries=$(jq -c '.entries[] | select(.kind == "rollup" and .deprecated == false)' "$REPORT")

  # Drop packages that already have their own open PR.
  local included=() pkg
  while IFS= read -r e; do
    [[ -z "$e" ]] && continue
    pkg=$(jq -r '.package' <<<"$e")
    if has_open_pr "$(branch_for "$pkg")"; then
      echo "  skip $pkg — it has its own open PR"
      continue
    fi
    included+=("$e")
  done <<<"$entries"

  if [[ ${#included[@]} -eq 0 ]]; then
    echo "Rollup: nothing to do."
    # An empty Rollup means everything in it merged. Retire a stale branch.
    if has_open_pr "$ROLLUP_BRANCH"; then
      echo "  (an open Rollup PR exists but has no remaining content — leaving it for a human)"
    fi
    return
  fi

  # The Rollup branch is bot-owned and disposable: rebuild it from main every
  # run, so it always reflects what is currently behind. Nobody should push to it.
  git checkout -B "$ROLLUP_BRANCH" origin/main

  for e in "${included[@]}"; do
    pkg=$(jq -r '.package'   <<<"$e")
    local proposed; proposed=$(jq -r '.proposed' <<<"$e")
    echo "  $pkg → $proposed"
    ./bin/dep-checker patch --package="$pkg" --version="$proposed"
  done

  # Only ever the Catalog. Never generators/.
  git add internal/deps/npm.go
  if git diff --cached --quiet; then
    echo "Rollup: patches were all no-ops — nothing to commit."
    git checkout main
    return
  fi

  local title="chore(deps): roll up ${#included[@]} dependency update(s)"
  git commit -m "$title"
  git push origin "$ROLLUP_BRANCH" --force-with-lease

  local body
  # shellcheck disable=SC2016  # backticks are markdown in the printf format
  body=$(
    printf '## Dependency Rollup\n\n'
    printf 'Every update here **satisfies its existing constraint** — the package itself promises compatibility.\n\n'
    printf '| Package | From | To |\n| --- | --- | --- |\n'
    for e in "${included[@]}"; do
      printf '| `%s` | `%s` | `%s` |\n' \
        "$(jq -r '.package' <<<"$e")" "$(jq -r '.pin' <<<"$e")" "$(jq -r '.proposed' <<<"$e")"
    done
    printf '\n> ⚠️ **This branch is rebuilt from scratch on every run — do not push to it.**\n'
    printf '> If one of these breaks the templates, eject it into its own PR:\n'
    printf '> the next Rollup will then exclude it automatically.\n\n'
    printf 'Generator manifests are **not** bumped here. They are derived at release from\n'
    printf 'Fingerprints, so merge order cannot corrupt them (ADR-0002).\n\n'
    printf -- '---\n🤖 [dep-checker](../tree/main/tools/dep-checker)\n'
  )

  if has_open_pr "$ROLLUP_BRANCH"; then
    echo "Rollup: existing PR updated by force-push."
  else
    gh pr create --title "$title" --body "$body" --base main --head "$ROLLUP_BRANCH" \
      --label dependencies --repo "$GITHUB_REPOSITORY"
  fi
  git checkout main
}

# ─────────────────────────────────────────────────────────────────────────────
# Migrations — one PR per package, human-owned once raised
# ─────────────────────────────────────────────────────────────────────────────
build_migrations() {
  local entries
  entries=$(jq -c '.entries[] | select(.kind == "migration" and .deprecated == false)' "$REPORT")
  [[ -z "$entries" ]] && { echo "Migrations: none."; return; }

  while IFS= read -r e; do
    [[ -z "$e" ]] && continue
    local pkg pin proposed branch
    pkg=$(jq -r '.package'  <<<"$e")
    pin=$(jq -r '.pin'      <<<"$e")
    proposed=$(jq -r '.proposed' <<<"$e")
    branch=$(branch_for "$pkg")

    if has_open_pr "$branch"; then
      echo "  $pkg — open PR already exists; not touching it."
      continue
    fi

    echo "  $pkg: $pin → $proposed (breaks the constraint)"
    git checkout -B "$branch" origin/main
    ./bin/dep-checker patch --package="$pkg" --version="$proposed"
    git add internal/deps/npm.go
    if git diff --cached --quiet; then
      git checkout main
      continue
    fi

    git commit -m "chore(deps)!: migrate ${pkg} to ${proposed}"
    git push origin "$branch" --force-with-lease

    # shellcheck disable=SC2016  # backticks are markdown in the printf format
    gh pr create \
      --title "chore(deps)!: migrate \`${pkg}\` ${pin} → ${proposed}" \
      --body "$(printf '## Migration: `%s`\n\n`%s` → `%s`\n\n**This breaks the existing constraint**, so it is not a routine bump — the package is signalling a breaking change. Generator code or templates may need to change alongside it.\n\nCI will tell you: if test-flows passes as-is, this is just a version bump after all. If it fails, that failure *is* the work.\n\nThis branch is yours now — the bot will not touch it again. Leaving this PR open tells the bot you have this in hand, and keeps `%s` out of the Rollup.\n\n---\n🤖 [dep-checker](../tree/main/tools/dep-checker)\n' "$pkg" "$pin" "$proposed" "$pkg")" \
      --base main --head "$branch" \
      --label dependencies \
      --repo "$GITHUB_REPOSITORY"
    git checkout main
  done <<<"$entries"
}

# ─────────────────────────────────────────────────────────────────────────────
# Deprecations — issue only, never a PR
# ─────────────────────────────────────────────────────────────────────────────
build_deprecation_issues() {
  local entries
  entries=$(jq -c '.entries[] | select(.deprecated)' "$REPORT")
  [[ -z "$entries" ]] && { echo "Deprecations: none."; return; }

  while IFS= read -r e; do
    [[ -z "$e" ]] && continue
    local pkg notice title existing num state
    pkg=$(jq -r '.package' <<<"$e")
    notice=$(jq -r '.deprecation_notice // "(no notice given)"' <<<"$e")

    # The dedup key is one package, so it is stable. The old script keyed on the
    # issue *title*, which contained a variable package list — which is why the
    # same @types/bcryptjs deprecation got filed twice under different titles
    # (#102 and #138).
    title="deprecated: ${pkg} (npm)"

    existing=$(gh issue list --search "\"$title\" in:title" --state all \
      --json number,state --jq '.[0] // empty' --repo "$GITHUB_REPOSITORY")

    if [[ -n "$existing" ]]; then
      num=$(jq -r '.number' <<<"$existing")
      state=$(jq -r '.state' <<<"$existing")
      if [[ "$state" == "CLOSED" ]]; then
        gh issue reopen "$num" --repo "$GITHUB_REPOSITORY"
        gh issue comment "$num" --repo "$GITHUB_REPOSITORY" \
          --body "Still deprecated as of $(date -u +%Y-%m-%d). Reopening."
      fi
      echo "  $pkg — issue #$num already tracks this."
      continue
    fi

    echo "  $pkg — filing a deprecation issue."
    # shellcheck disable=SC2016  # backticks are markdown in the printf format
    gh issue create --title "$title" --repo "$GITHUB_REPOSITORY" \
      --label "dependencies,deprecated" \
      --body "$(printf '## `%s` is deprecated\n\n> %s\n\n### Why there is no PR\n\nA deprecation is **never** fixed by a version bump. The remedy is a *different package*, or none at all — a rename, or a removal. Bumping to the latest deprecated version resolves nothing, so the bot will not open one.\n\nThis package is also excluded from the Rollup until it is dealt with.\n\n### What to do\n\n- Decide on the replacement (the notice above usually names it)\n- Update the Catalog (`internal/deps/npm.go`) and every generator that names `%s`\n- Close this issue\n\n---\n🤖 [dep-checker](../tree/main/tools/dep-checker)\n' "$pkg" "$notice" "$pkg")"
  done <<<"$entries"
}

echo "=== Rollup ==="
build_rollup
echo ""
echo "=== Migrations ==="
build_migrations
echo ""
echo "=== Deprecations ==="
build_deprecation_issues
