# Gate 001 – Reviewer Checklist

Use this for every PR that touches EventV1, PQC harness, signing, or Operator Desk tooling.

Author-facing twin: `PR_REVIEW_CHECKLIST.md`. Do not merge the two files. If they disagree, diff them before editing either.

A seat demo is not a substitute for a real verifier. This repo's `tools/eventv1 verify` is a scaffold.

## Correctness
- [ ] Changes are covered by `go test` (or explicit reason why not)
- [ ] Envelopes remain valid JSON (`scripts/check-envelopes.sh`, which is `jq empty`)
- [ ] Deterministic behavior when `PQC_TEST_SEED` is set

## Security
- [ ] No secrets in diff
- [ ] Pre-commit secret patterns still appropriate
- [ ] Self-hosted signing path does not leak key material into logs
- [ ] gitleaks stays on ubuntu-latest. It is not added to the signing runner.

## Stack Hygiene
- [ ] gofmt clean
- [ ] go vet clean
- [ ] `go.mod` stays at 1.26 unless a dedicated toolchain PR says otherwise. The old 1.21–1.22 line is retired.
- [ ] Composite action still works end-to-end
- [ ] No new language runtime

## Operational
- [ ] Runbook (`OPS_SIGNER_RUNBOOK.md`) still accurate if signing changed
- [ ] Workflows that target `self-hosted, signing, linux` still correct

## Proof Sells
- [ ] Clear before/after proof artifact (test output, vector, screenshot of CI)
- [ ] Ready to demonstrate to a customer or internal stakeholder in < 5 minutes
- [ ] The demo is labeled scaffold or verifier. It is not both.
