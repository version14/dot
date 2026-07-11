#!/usr/bin/env bash
# Opens (or updates) one PR per outdated/deprecated dependency group found in
# dep-report.json, and tracks deprecations with an issue. Runs from the repo
# root with dep-report.json present (see dep-scan-core action) and requires
# GH_TOKEN + GITHUB_REPOSITORY in the environment.
set -euo pipefail

NEEDS_WORK=$(jq -r '
  .entries[] | select(.outdated or .deprecated)
  | .ecosystem + "|" + .package + "|" + .latest
' dep-report.json | sort -u)

if [[ -z "$NEEDS_WORK" ]]; then
  echo "All template dependencies are up to date."
  exit 0
fi

git config user.name  "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

# Helper: map @types/<name> -> <name>
base_package() {
  local pkg="$1"
  if [[ "$pkg" == @types/* ]]; then
    echo "${pkg#@types/}"
  else
    echo "$pkg"
  fi
}

update_rank() {
  local update_type="$1"
  case "$update_type" in
    major) echo 3 ;;
    minor) echo 2 ;;
    patch) echo 1 ;;
    *) echo 0 ;;
  esac
}

max_update_type() {
  local entries="${1-}"
  if [[ -z "$entries" ]]; then
    entries=$(cat)
  fi
  local max=0
  local max_type="unknown"
  while IFS= read -r t; do
    local rank
    rank=$(update_rank "$t")
    if [[ "$rank" -gt "$max" ]]; then
      max="$rank"
      max_type="$t"
    fi
  done <<< "$entries"
  echo "$max_type"
}

# Track which logical groups have already been processed.
declare -A GROUP_SEEN
MINOR_PATCH_GROUPED=false

# Process one (ecosystem, package, latest) tuple at a time, but
# combine @types/<name> with <name> when both exist for npm.
while IFS='|' read -r ECOSYSTEM PACKAGE LATEST; do
  echo ""
  echo "=== $ECOSYSTEM/$PACKAGE (latest: $LATEST) ==="

  base="$(base_package "$PACKAGE")"
  if [[ "$ECOSYSTEM" = "npm" ]]; then
    GROUP_ENTRIES=$(jq -c \
      --arg eco "$ECOSYSTEM" \
      --arg pkg "$PACKAGE" \
      --arg base "$base" \
      '.entries[]
       | select(.ecosystem == $eco and (
           .package == $pkg or
           .package == ("@types/" + $base) or
           .package == $base
         ))' \
      dep-report.json)
    ENTRY_UPDATE_TYPE=$(max_update_type "$(echo "$GROUP_ENTRIES" | jq -r '.update_type')")
  else
    ENTRY_UPDATE_TYPE=$(max_update_type "$(jq -r \
      --arg eco "$ECOSYSTEM" \
      --arg pkg "$PACKAGE" \
      '.entries[] | select(.ecosystem == $eco and .package == $pkg) | .update_type' \
      dep-report.json)")
  fi

  if [[ "$ECOSYSTEM" = "npm" ]] && { [[ "$ENTRY_UPDATE_TYPE" = "minor" ]] || [[ "$ENTRY_UPDATE_TYPE" = "patch" ]]; }; then
    GROUP_KEY="npm|minor-patch"
  else
    GROUP_KEY="${ECOSYSTEM}|${base}|${ENTRY_UPDATE_TYPE}"
  fi
  if [[ "${GROUP_SEEN[$GROUP_KEY]+_}" ]]; then
    echo "  Group ${GROUP_KEY} already processed — skipping."
    continue
  fi
  GROUP_SEEN[$GROUP_KEY]=1

  BRANCH=""
  if [[ "$ECOSYSTEM" = "npm" ]]; then
    if [[ "$ENTRY_UPDATE_TYPE" = "minor" ]] || [[ "$ENTRY_UPDATE_TYPE" = "patch" ]]; then
      if [[ "$MINOR_PATCH_GROUPED" = "true" ]]; then
        echo "  Minor/patch npm group already processed — skipping."
        continue
      fi
      MINOR_PATCH_GROUPED=true
      ENTRIES=$(jq -c '
        .entries[]
        | select(.ecosystem == "npm" and (.update_type == "minor" or .update_type == "patch"))
      ' dep-report.json)
      BRANCH="deps/npm/minor-and-patch"
    else
      ENTRIES="$GROUP_ENTRIES"
    fi
  else
    ENTRIES=$(jq -c \
      --arg eco "$ECOSYSTEM" \
      --arg pkg "$PACKAGE" \
      '.entries[] | select(.ecosystem == $eco and .package == $pkg)' \
      dep-report.json)
  fi

  if [[ -z "$ENTRIES" ]]; then
    echo "  No entries found for group $GROUP_KEY — skipping."
    continue
  fi

  if [[ -z "$BRANCH" ]]; then
    if [[ "$ECOSYSTEM" = "npm" ]]; then
      BRANCH="deps/${ECOSYSTEM}/${base}"
    else
      BRANCH="deps/${ECOSYSTEM}/${PACKAGE}"
    fi
  fi

  # Idempotency: skip if an open PR already targets this branch.
  OPEN_PRS=$(gh pr list \
    --head "$BRANCH" \
    --state open \
    --json number \
    --jq 'length' \
    --repo "$GITHUB_REPOSITORY")
  if [[ "$OPEN_PRS" -gt 0 ]]; then
    echo "  Open PR already exists on $BRANCH — skipping."
    continue
  fi

  IS_DEPRECATED=$(echo "$ENTRIES" | jq -r 'select(.deprecated) | .deprecated' | head -1)
  NOTICE=$(echo "$ENTRIES" | jq -r 'select(.deprecated) | .deprecation_notice // ""' | head -1)

  # Create a fresh branch from main.
  git checkout -B "$BRANCH" origin/main

  # Patch every generator that declares this package group.
  # --skip-manifest-bump prevents each patch call from incrementing
  # the version independently; we do one bump per generator below.
  PATCH_FAILED=false
  declare -A GEN_MAX_RANK
  declare -A GEN_MAX_TYPE
  while IFS= read -r entry; do
    GENERATOR=$(echo "$entry" | jq -r '.generator')
    ENTRY_PACKAGE=$(echo "$entry" | jq -r '.package')
    CURRENT=$(echo  "$entry" | jq -r '.current')
    ENTRY_LATEST=$(echo "$entry" | jq -r '.latest')
    ENTRY_UPDATE_TYPE=$(echo "$entry" | jq -r '.update_type')

    # Track the highest-severity update type seen for each generator.
    cur_rank="${GEN_MAX_RANK[$GENERATOR]:-0}"
    new_rank=$(update_rank "$ENTRY_UPDATE_TYPE")
    if [[ "$new_rank" -gt "$cur_rank" ]]; then
      GEN_MAX_RANK[$GENERATOR]="$new_rank"
      GEN_MAX_TYPE[$GENERATOR]="$ENTRY_UPDATE_TYPE"
    fi

    echo "  Patching $GENERATOR: $ENTRY_PACKAGE $CURRENT → ^$ENTRY_LATEST"
    if ./bin/dep-checker patch \
        --generator="$GENERATOR" \
        --package="$ENTRY_PACKAGE" \
        --current="$CURRENT" \
        --latest="$ENTRY_LATEST" \
        --skip-manifest-bump > /dev/null; then
      echo "  OK: $GENERATOR patched"
    else
      echo "  ERROR: patch failed for $GENERATOR" >&2
      PATCH_FAILED=true
      break
    fi
  done <<< "$ENTRIES"

  # Bump each generator's manifest exactly once, using the max
  # update type observed for that generator in this branch.
  if [[ "$PATCH_FAILED" != "true" ]]; then
    for generator in "${!GEN_MAX_TYPE[@]}"; do
      bump_type="${GEN_MAX_TYPE[$generator]}"
      echo "  Bumping manifest for $generator ($bump_type)"
      if ! ./bin/dep-checker bump-manifest \
          --generator="$generator" \
          --type="$bump_type"; then
        echo "  WARN: could not bump manifest for $generator"
      fi
    done
  fi

  if [[ "$PATCH_FAILED" = "true" ]]; then
    git checkout main
    continue
  fi

  # Commit, push, open PR.
  CHANGED_GENS=$(echo "$ENTRIES" | jq -r '.generator' | sed 's/^/- `/' | sed 's/$/ `/')
  PACKAGE_LABEL=$(echo "$ENTRIES" | jq -r '.package' | sort -u | paste -sd ", " -)
  git add generators/
  git commit -m "chore(deps): bump ${PACKAGE_LABEL} in templates"
  git push origin "$BRANCH" --force-with-lease

  PR_BODY=$(printf \
    '## Dependency bump\n\nBumps `%s` (%s) in template generators.\n\n### Affected generators\n\n%s\n\n---\n🤖 Generated by [dep-checker](../tree/main/tools/dep-checker)' \
    "$PACKAGE_LABEL" "$ECOSYSTEM" "$CHANGED_GENS")

  gh pr create \
    --title  "chore(deps): bump ${PACKAGE_LABEL} in templates" \
    --body   "$PR_BODY" \
    --base   main \
    --head   "$BRANCH" \
    --label  "dependencies" \
    --repo   "$GITHUB_REPOSITORY"

  # For deprecated packages, also open/reopen a tracking issue.
  if [[ "$IS_DEPRECATED" = "true" ]]; then
    ISSUE_TITLE="deprecated: ${PACKAGE_LABEL} (${ECOSYSTEM}) used in templates"

    EXISTING=$(gh issue list \
      --search  "\"$ISSUE_TITLE\" in:title" \
      --state   all \
      --json    number,state \
      --jq      '.[0] // empty' \
      --repo    "$GITHUB_REPOSITORY")

    if [[ -n "$EXISTING" ]]; then
      ISSUE_NUM=$(echo   "$EXISTING" | jq -r '.number')
      ISSUE_STATE=$(echo "$EXISTING" | jq -r '.state')
      if [[ "$ISSUE_STATE" = "CLOSED" ]]; then
        gh issue reopen "$ISSUE_NUM" --repo "$GITHUB_REPOSITORY"
        gh issue comment "$ISSUE_NUM" \
          --body "Package is still deprecated as of $(date -u +%Y-%m-%d). Reopening for review." \
          --repo "$GITHUB_REPOSITORY"
      fi
    else
      ISSUE_BODY=$(printf \
        '## Deprecated package\n\n`%s` (%s) is marked deprecated on its registry.\n\n**Deprecation notice:** %s\n\n### Affected generators\n\n%s\n\n### What to do\n\n- Check whether there is a recommended replacement package\n- Update all affected generators to use the replacement\n- Close this issue once resolved\n\nA version-bump PR has been opened — review whether it resolves the deprecation or whether a package rename is needed.' \
        "$PACKAGE_LABEL" "$ECOSYSTEM" "$NOTICE" "$CHANGED_GENS")

      gh issue create \
        --title  "$ISSUE_TITLE" \
        --body   "$ISSUE_BODY" \
        --label  "dependencies,deprecated" \
        --repo   "$GITHUB_REPOSITORY"
    fi
  fi

  git checkout main

done <<< "$NEEDS_WORK"
