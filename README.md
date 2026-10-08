# Operator Desk

**Operator Pages LLC · Quantomguard AI**  
CEO / Founder: Daniel H. Pennick

Local operator workstation package for EventV1 + PQC signing / verification pipelines.  
Built for speed-to-proof. Reuses existing stack pieces. MVP-ready.

## Quick Start (MVP in minutes)

```bash
# 1. Unpack
unzip operator-desk.zip && cd operator-desk

# 2. Install local hooks (zero Python dependency)
./scripts/setup-hooks.sh

# 3. Set test seed (repo secret style – never commit real values)
export PQC_TEST_SEED="0000000000000000000000000000000000000000000000000000000000000000"  # 32-byte hex example

# 4. Run local verification pipeline
./scripts/run-eventv1-verify.sh
```

## Stack Alignment (exact inventory)

| Layer              | Component                                      | Status in this desk |
|--------------------|------------------------------------------------|---------------------|
| Language           | Go 1.21–1.22                                   | go.mod pinned       |
| Format / Lint      | gofmt + go vet + go test                       | pre-commit enforced |
| Optional deep lint | golangci-lint                                  | documented          |
| CI                 | GitHub Actions (ubuntu-latest + self-hosted)   | workflows included  |
| Composite action   | .github/actions/eventv1-verify                 | present             |
| Tools              | tools/eventv1 + tools/pqc-harness              | scaffolded          |
| Vectors            | vectors/eventv1/*_envelope.json                | sample included     |
| Hooks              | .githooks/pre-commit (pure shell)              | active              |
| Secrets            | PQC_TEST_SEED (test only)                      | documented          |
| Signing            | HSM / host-signer on self-hosted runner        | runbook included    |

## Directory Map

```
operator-desk/
├── README.md                          ← you are here
├── go.mod
├── bin/                               ← built binaries land here
├── tools/
│   ├── eventv1/                       ← generator, envelope maker, verify CLI
│   └── pqc-harness/                   ← PQC test harness
├── vectors/eventv1/                   ← golden envelopes
├── .github/
│   ├── actions/eventv1-verify/        ← composite action
│   └── workflows/                     ← pqc-local-test, pqc-sign-test, etc.
├── .githooks/pre-commit               ← pure-shell enforcement
├── scripts/
│   ├── setup-hooks.sh
│   └── run-eventv1-verify.sh
└── docs/gate001/
    ├── PR_BODY.md
    ├── REVIEWER_CHECKLIST.md
    ├── OPS_SIGNER_RUNBOOK.md
    └── PR_REVIEW_CHECKLIST.md
```

## Operator Workflows

### Local proof loop
1. Edit / generate envelope
2. `go test ./...` (or let pre-commit run it)
3. `./scripts/run-eventv1-verify.sh`
4. Inspect vectors/eventv1/

### CI proof
- Push → `pqc-local-test.yml` runs on ubuntu-latest
- Signing path → `pqc-sign-test-self-hosted.yml` (requires self-hosted runner labeled `signing`)

### Production signing
See `docs/gate001/OPS_SIGNER_RUNBOOK.md`. Public key expected at `/etc/host-signer/pubkey.hex` on the signer host.

## Secrets Policy
- **Never** commit real seeds or private keys.
- `PQC_TEST_SEED` is for deterministic test vectors only.
- Local secret scan is conservative grep in pre-commit.
- Recommended CI: gitleaks or trufflehog.

## Next Big Thing Watch
This desk is Go-primary. Node.js / Python adjacency points are noted in `docs/gate001/` for future seat expansion (CLI wrappers, dashboard, automation agents).

---

**Proof sells. Sell seats.**  
Operator Desk v0.1 – ready for internal distribution.
