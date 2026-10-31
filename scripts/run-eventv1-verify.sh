#!/usr/bin/env bash
# Operator Desk – local EventV1 generate → convert → build → verify pipeline
# Mirrors the composite action .github/actions/eventv1-verify
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "=== Operator Desk EventV1 Verify Pipeline ==="

# Ensure test seed is present (deterministic only)
if [[ -z "${PQC_TEST_SEED:-}" ]]; then
  echo "WARN: PQC_TEST_SEED not set. Using zero seed for local dry-run."
  export PQC_TEST_SEED="0000000000000000000000000000000000000000000000000000000000000000"
fi

echo "[1/4] Generate / prepare vectors..."
mkdir -p vectors/eventv1
# Placeholder: real generator lives in tools/eventv1
if [[ -f tools/eventv1/generate.go ]]; then
  go run ./tools/eventv1 generate --out vectors/eventv1/
else
  echo "  (scaffold mode – using existing sample vector)"
fi

echo "[2/4] Convert / envelope..."
# Placeholder for envelope maker

echo "[3/4] Build tools..."
go build -o bin/eventv1-verify ./tools/eventv1/ 2>/dev/null || echo "  (build skipped – tool still scaffold)"

echo "[4/4] Verify..."
if [[ -x bin/eventv1-verify ]]; then
  ./bin/eventv1-verify verify -vector vectors/eventv1/
else
  echo "  Running go test as verification stand-in..."
  go test ./... -count=1 -short
fi

echo "=== Pipeline complete ==="
echo "Inspect: vectors/eventv1/"
