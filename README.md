# Operator Desk

**Operator Pages LLC · Quantomguard AI**  
CEO / Founder: Daniel H. Pennick

Local operator workstation package for EventV1 + PQC signing / verification pipelines.  
Built for speed-to-proof. Reuses existing stack pieces. MVP-ready.

This tree is a scaffold. `tools/eventv1 verify` does not check a signature. See `docs/CONTRIBUTING.md`.

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
│   ├── check-envelopes.sh
│   └── run-eventv1-verify.sh
└── docs/
    ├── CONTRIBUTING.md
    ├── ships/
    └── gate001/
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
- Hosted secret scan → gitleaks in `pre-commit-checks.yml` (ubuntu-latest only)

### Production signing
See `docs/gate001/OPS_SIGNER_RUNBOOK.md`. Public key expected at `/etc/host-signer/pubkey.hex` on the signer host.

## What this repo does not verify

`tools/eventv1 verify` prints a pass. `go test` here is a smoke stub. Neither is a QuantumGuard verification.

The frozen verifier is `danpennic-sys/quantumguard`. Do not reimplement it under `web/` in this repo. A Node seat may show an artifact this pipeline produced. It does not become the verifier.

## Secrets Policy
- **Never** commit real seeds or private keys.
- `PQC_TEST_SEED` is for deterministic test vectors only. The all-zero hex above is a placeholder. Do not replace it in git with a real seed.
- Local secret scan is conservative grep in pre-commit.
- Hosted CI runs gitleaks on ubuntu-latest. Do not put that scanner on the signing runner.

## Next Big Thing Watch
This desk is Go-primary. Node.js may sell a seat on an artifact. Python stays off the hook path and off the signer. No retarget of the primary language in this change.

---

**Proof sells. Sell seats.**  
Operator Desk – Go 1.26 stack-synced. Scaffold is not the verifier.
