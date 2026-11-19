# Operator Desk

Operator Pages LLC. Local EventV1 desk.

This is a scaffold that can now fail. It is not QuantumGuard. It is not `qg verify`. It is not post-quantum. It is not a seat.

## Proof

Committed fixtures, algorithm `ed25519-test`:

| File | Required result | Exit |
|---|---|---|
| `vectors/eventv1/good_ed25519_envelope.json` | PASS | 0 |
| `vectors/eventv1/bad_ed25519_envelope.json` | FAIL | 1 |
| `vectors/eventv1/sample_envelope.json` | INDETERMINATE | 2 |

```bash
./scripts/setup-hooks.sh
bash scripts/run-eventv1-verify.sh
```

The script builds `bin/eventv1-verify` (gitignored) and requires all three results. `PQC_TEST_SEED` is not read during verify. The public key is inside the fixture. The committed key is the documented 32-byte zero seed.

`eventv1 generate --out vectors/eventv1` rewrites the good and bad fixtures. If `PQC_TEST_SEED` is set to anything other than 64 hex zeros, the rewrite will not match the committed files. Do not run generate in CI.

## What is not true yet

- `tools/pqc-harness/` is a note, not a harness. No ML-DSA.
- `envelope` is not implemented. It exits 2.
- Workflow filenames still say `pqc-`. The jobs do not perform post-quantum verification and do not HSM-sign.
- A file at `/etc/host-signer/pubkey.hex` proves the file exists. It does not prove a signature.
- This repository is public. Do not put a private key in it.
- `docs/gate001/` runbooks are not a second verifier.

## Layout

```
tools/eventv1/          verify and generate
vectors/eventv1/        three fixtures
scripts/run-eventv1-verify.sh
.github/actions/eventv1-verify/
docs/CONTRIBUTING.md
docs/gate001/
```

## Stack that this tree actually pins

Go 1.26 (`go.mod`). gofmt, `go vet`, `go test` via `.githooks/pre-commit`. GitHub Actions on `ubuntu-latest`. The self-hosted workflow only checks that a pubkey file exists, then runs the same ed25519-test proof.

QuantumGuard v1 stays in [danpennic-sys/quantumguard](https://github.com/danpennic-sys/quantumguard). Do not merge these verdicts.
