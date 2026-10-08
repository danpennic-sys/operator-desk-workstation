## Gate 001 – Operator Desk / EventV1 PR Body Template

### Summary
<!-- One-paragraph description of the change -->

### Proof of Correctness
- [ ] Local: `./scripts/run-eventv1-verify.sh` passes
- [ ] Pre-commit: pure-shell hooks pass
- [ ] CI: `pqc-local-test.yml` green
- [ ] Vectors: new/updated envelopes under `vectors/eventv1/` are jq-valid

### Risk & Rollback
- Risk level: Low / Medium / High
- Rollback plan:

### Secrets
- [ ] No real seeds or private keys introduced
- [ ] `PQC_TEST_SEED` used only for deterministic tests

### Vertical Alignment
- [ ] Reuses existing Go + Actions stack
- [ ] No new language runtime required for this change
- [ ] Ready for seat / MVP packaging if applicable

### Reviewer Notes
<!-- Anything the reviewer should focus on -->
