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
| Language           | **Go 1.26** (Operator Pages stack standard)    | go.mod pinned       |
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
├── README.md
├── go.mod
├── bin/
├── tools/
│   ├── eventv1/
│   └── pqc-harness/
├── vectors/eventv1/
├── .github/
│   ├── actions/eventv1-verify/
│   └── workflows/
├── .githooks/pre-commit
├── scripts/
│   ├── setup-hooks.sh
│   └── run-eventv1-verify.sh
└── docs/gate001/
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
This desk is Go-primary. Node.js 26+ / Python 3.12+ adjacency points are noted in `docs/gate001/` for future seat expansion.

---

**Proof sells. Sell seats.**  
Operator Desk – Go 1.26 stack-synced.
