#!/usr/bin/env bash
# Fail-closed gofmt check.
# Do not rewrite this as `gofmt -l | grep -q` under `set -o pipefail`.
# grep -q exits on the first match, the writer can die with SIGPIPE (141),
# and pipefail then lets `if !` report a dirty tree as clean.
set -euo pipefail

root="${1:-.}"
if [[ ! -d "$root" ]]; then
  echo "ERROR: gofmt root is not a directory: $root" >&2
  exit 2
fi

unformatted="$(cd "$root" && gofmt -l .)"
if [[ -n "$unformatted" ]]; then
  echo "ERROR: gofmt issues found. Run: gofmt -w ." >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

echo "[gofmt] clean"
