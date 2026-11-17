# Gate 001 – PR Review Checklist (Author + Reviewer)

Copy into PR description or use alongside `REVIEWER_CHECKLIST.md`. These are two files on purpose. Diff them if they drift. Do not fold one into the other in a drive-by edit.

A seat demo is not a substitute for a real verifier.

## Author
- [ ] `gofmt -w .` run
- [ ] `go vet ./...` clean
- [ ] `go test ./... -count=1` green
- [ ] Pre-commit hooks installed and passing (`./scripts/setup-hooks.sh`)
- [ ] Sample / new envelopes under `vectors/eventv1/` pass `scripts/check-envelopes.sh`
- [ ] No secrets
- [ ] PR body follows `PR_BODY.md` template
- [ ] `go.mod` still says `go 1.26`, or this PR is the dedicated toolchain bump and nothing else

## Reviewer
- [ ] Code matches stated intent
- [ ] Tests actually exercise the changed path
- [ ] CI workflows still green (especially self-hosted signing path if touched)
- [ ] Docs updated if behavior or paths changed
- [ ] Vertical alignment: still pure Go + existing Actions surface. Python stays off the hooks.
- [ ] Proof artifact present (test log, vector, CI link)
- [ ] The PR does not call this scaffold a signature verify

## Sign-off
Reviewer: _______________  
Date: _______________  
Gate status: PASS / HOLD / REJECT
