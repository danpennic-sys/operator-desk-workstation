#!/usr/bin/env bash
# Desk proof. Three committed fixtures. Not qg verify. Not post-quantum.
# Verify uses the public key inside each file. PQC_TEST_SEED is not read.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "=== EventV1 desk proof: ed25519-test, not QuantumGuard, not PQC ==="

mkdir -p bin
go build -o bin/eventv1-verify ./tools/eventv1/

expect() {
  local want="$1" file="$2" out rc first
  set +e
  out="$(./bin/eventv1-verify verify --vector "$file")"
  rc=$?
  set -e
  printf '%s\n' "$out"
  first="${out%%$'\n'*}"
  case "$want" in
    PASS)
      [ "$rc" -eq 0 ] || { echo "expected PASS exit 0, got $rc"; exit 1; }
      ;;
    FAIL)
      [ "$rc" -eq 1 ] || { echo "expected FAIL exit 1, got $rc"; exit 1; }
      ;;
    INDETERMINATE)
      [ "$rc" -eq 2 ] || { echo "expected INDETERMINATE exit 2, got $rc"; exit 1; }
      ;;
    *)
      echo "unknown expected verdict $want"
      exit 1
      ;;
  esac
  case "$first" in
    "$want"*) ;;
    *)
      echo "expected first line to start with $want"
      exit 1
      ;;
  esac
}

expect PASS vectors/eventv1/good_ed25519_envelope.json
expect FAIL vectors/eventv1/bad_ed25519_envelope.json
expect INDETERMINATE vectors/eventv1/sample_envelope.json

echo "[go test]"
go test ./... -count=1

echo "=== PASS, FAIL, and INDETERMINATE each observed once ==="
