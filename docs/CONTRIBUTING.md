# Contributing

Operator Pages LLC. This repository is the Operator Desk workstation scaffold.

## Verify of record

`tools/eventv1` is a scaffold. The `verify` command prints a pass. It does not check a signature. `go test` in this repo is a smoke stub (`true != true`). A green run of this scaffold is not a QuantumGuard verification.

The frozen verifier is a different repository: `danpennic-sys/quantumguard` (Ed25519, v1 frozen). This change set does not modify that repository.

A Node seat may display an artifact. It must not reimplement verify, and it must not live under `web/` in this repo. Python stays off `.githooks/` and off the signing runner. Do not add `.pre-commit-config.yaml` unless that zero-Python rule is retired in the same written change.

## Go pin

`go.mod` says `go 1.26`. Do not describe the pin as 1.21–1.22. Do not bump the toolchain in the same change as a crypto or envelope change.

## Envelopes

The only envelope check in this repo is `scripts/check-envelopes.sh` (`jq empty` on `*_envelope.json`). The hook, the composite action, and hosted CI call that script. Do not add a second schema here.

`jq` is required for a full check. The hook warns and continues if `jq` is missing. Hosted CI sets `REQUIRE_JQ=1` and fails closed.

## Secrets

`PQC_TEST_SEED` is a test-only name. Never commit a real seed or a private key. The all-zero 64-hex string checked in is a documented placeholder, not a production secret. Do not replace it in git with a real value.

Hosted secret scan is gitleaks on `ubuntu-latest` only. Do not add that scanner to the self-hosted signing workflow.

## Proof

A seat demo is not a substitute for a real verifier. Gate 001 still applies. See `docs/gate001/REVIEWER_CHECKLIST.md` and `docs/gate001/PR_REVIEW_CHECKLIST.md`.
