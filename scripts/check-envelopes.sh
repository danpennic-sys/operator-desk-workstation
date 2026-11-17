#!/usr/bin/env bash
# One envelope check for the hook, the composite action, and hosted CI.
# Structural only: jq empty. Not a signature check. Not the quantumguard verifier.
set -euo pipefail

DIR="${1:-vectors/eventv1}"

if ! command -v jq >/dev/null 2>&1; then
  if [[ "${REQUIRE_JQ:-}" == "1" ]]; then
    echo "ERROR: jq is required"
    exit 1
  fi
  echo "WARN: jq not found – skipping envelope JSON validation"
  exit 0
fi

shopt -s nullglob
files=("$DIR"/*_envelope.json)
if [[ ${#files[@]} -eq 0 ]]; then
  echo "[jq] no *_envelope.json under $DIR"
  exit 0
fi

for f in "${files[@]}"; do
  echo "[jq] validating $f"
  jq empty "$f"
done
