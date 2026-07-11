#!/usr/bin/env bash
# Posts (or updates) the dependency-report comment on the current PR from
# pr-dep-report.json. Requires GH_TOKEN, PR_NUMBER, and GITHUB_REPOSITORY in
# the environment (all standard Actions step env).
set -euo pipefail

ENTRY_COUNT=$(jq '.entries | length' pr-dep-report.json)
if [[ "$ENTRY_COUNT" -eq 0 ]]; then
  echo "No tracked dependencies in the changed generators — skipping dep comment."
  exit 0
fi

./bin/dep-checker report --input=pr-dep-report.json --output=pr-dep-report.md

MARKER="<!-- dep-checker-report -->"
BODY="${MARKER}
$(cat pr-dep-report.md)"

# Upsert: update existing marker comment, or create a new one.
EXISTING_ID=$(gh api "repos/$GITHUB_REPOSITORY/issues/$PR_NUMBER/comments" \
  --jq '[.[] | select(.body | startswith("<!-- dep-checker-report -->"))] | first | .id // empty')

if [[ -n "$EXISTING_ID" ]]; then
  gh api "repos/$GITHUB_REPOSITORY/issues/comments/$EXISTING_ID" \
    --method PATCH \
    --field body="$BODY"
  echo "Updated dep comment $EXISTING_ID on PR #$PR_NUMBER."
else
  gh pr comment "$PR_NUMBER" --body "$BODY" --repo "$GITHUB_REPOSITORY"
  echo "Created dep comment on PR #$PR_NUMBER."
fi
