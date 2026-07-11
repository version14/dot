#!/usr/bin/env bash
# Filters dep-report.json down to the generators touched by the current PR,
# writing pr-dep-report.json. Requires GH_TOKEN, PR_NUMBER, and
# GITHUB_REPOSITORY in the environment (all standard Actions step env).
set -euo pipefail

# Generator directories touched by this PR.
# Use --paginate + per_page=100 so large PRs (200+ files) are fully covered.
# --jq extracts one generator name per line per page; sort -u deduplicates
# across pages; jq -R/jq -s converts the text list to a JSON array.
CHANGED_GENS=$(gh api "repos/$GITHUB_REPOSITORY/pulls/$PR_NUMBER/files?per_page=100" \
  --paginate \
  --jq '[.[] | select(.filename | startswith("generators/")) | .filename | split("/")[1]] | .[]' \
  | sort -u \
  | jq -R . | jq -s .)

echo "Changed generators: $CHANGED_GENS"

if [ "$CHANGED_GENS" = "[]" ]; then
  echo "No generator files changed in this PR — producing empty filtered report."
  jq '. + {entries: []}' dep-report.json > pr-dep-report.json
else
  jq --argjson gens "$CHANGED_GENS" \
    '. + {entries: [.entries[] | select(.generator as $g | ($gens | index($g)) != null)]}' \
    dep-report.json > pr-dep-report.json
  echo "Filtered to $(jq '.entries | length' pr-dep-report.json) entries from generators: $CHANGED_GENS"
fi
