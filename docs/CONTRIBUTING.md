# Contributing

Operator Desk only. Do not change QuantumGuard CONTRACT-v1 from this repo.

## Language

Go 1.26, as pinned in `go.mod`. Format with gofmt. `go test ./...` must pass.

## What verify is

`eventv1 verify` prints one of:

- `PASS` — exit 0. The embedded `ed25519-test` public key signed the canonical message.
- `FAIL` — exit 1. The algorithm was `ed25519-test` and the signature check failed.
- `INDETERMINATE` — exit 2. The file is not a checkable `ed25519-test` envelope.

PASS is not an operator identity. PASS is not `qg verify`. PASS is not post-quantum.

A directory exits 1 if any file is FAIL. FAIL outranks INDETERMINATE. INDETERMINATE outranks PASS.

## Canonical message

Exactly this JSON, compact, this key order, no HTML escaping, no timestamp:

```json
{"message":"...","type":"...","version":"eventv1"}
```

Unknown JSON fields, duplicate keys, and a wall-clock field are INDETERMINATE.

## Do not

- Do not make a test that cannot fail.
- Do not describe this desk as the Operator Pack.
- Do not commit a non-zero `PQC_TEST_SEED`, a private key, or a Stripe secret.
- The committed fixtures are the public 32-byte zero seed. `generate` overwrites them.
