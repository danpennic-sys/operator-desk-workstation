# Gate 001 – Reviewer Checklist

Use this for every PR that touches EventV1, PQC harness, signing, or Operator Desk tooling.

## Correctness
- [ ] Changes are covered by `go test` (or explicit reason why not)
- [ ] Envelopes remain valid JSON (jq empty)
- [ ] Deterministic behavior when `PQC_TEST_SEED` is set

## Security
- [ ] No secrets in diff
- [ ] Pre-commit secret patterns still appropriate
- [ ] Self-hosted signing path does not leak key material into logs

## Stack Hygiene
- [ ] gofmt clean
- [ ] go vet clean
- [ ] Go version stays in 1.21–1.22 range
- [ ] Composite action still works end-to-end

## Operational
- [ ] Runbook (`OPS_SIGNER_RUNBOOK.md`) still accurate if signing changed
- [ ] Workflows that target `self-hosted, signing, linux` still correct

## Proof Sells
- [ ] Clear before/after proof artifact (test output, vector, screenshot of CI)
- [ ] Ready to demonstrate to a customer or internal stakeholder in < 5 minutes
