# Gate 001 – PR Review Checklist (Author + Reviewer)

Copy into PR description or use alongside REVIEWER_CHECKLIST.md.

## Author
- [ ] `gofmt -w .` run
- [ ] `go vet ./...` clean
- [ ] `go test ./... -count=1` green
- [ ] Pre-commit hooks installed and passing (`./scripts/setup-hooks.sh`)
- [ ] Sample / new envelopes under `vectors/eventv1/` are valid
- [ ] No secrets
- [ ] PR body follows `PR_BODY.md` template

## Reviewer
- [ ] Code matches stated intent
- [ ] Tests actually exercise the changed path
- [ ] CI workflows still green (especially self-hosted signing path if touched)
- [ ] Docs updated if behavior or paths changed
- [ ] Vertical alignment: still pure Go + existing Actions surface
- [ ] Proof artifact present (test log, vector, CI link)

## Sign-off
Reviewer: _______________  
Date: _______________  
Gate status: PASS / HOLD / REJECT
