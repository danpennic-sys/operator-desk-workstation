# Gate 001 – Ops Signer Runbook

**Audience**: Operators running the self-hosted signing path  
**Runner labels required**: `self-hosted`, `signing`, `linux`

## Prerequisites
1. Self-hosted runner registered with the three labels above.
2. Host signer binary / HSM interface available on the machine.
3. Public key mounted (or present) at:
   ```
   /etc/host-signer/pubkey.hex
   ```
4. Repository secret `PQC_TEST_SEED` set for test runs only (never production keys).

## Normal Signing Flow
1. Trigger `pqc-sign-test-self-hosted.yml` (workflow_dispatch or push to main that touches tools/vectors).
2. Job lands on the labeled runner.
3. Composite action `.github/actions/eventv1-verify` executes:
   - generate → convert → build → verify
4. Signing step (when implemented) calls the host signer using the public key path above.
5. Resulting envelopes land in `vectors/eventv1/` and are uploaded as artifacts.

## Failure Modes & Recovery
| Symptom | Likely Cause | Action |
|---------|--------------|--------|
| Runner not picked up | Labels missing | Re-register runner with correct labels |
| pubkey.hex missing | Mount/volume issue | Restore file, restart runner service |
| Seed-related test failure | Wrong or missing `PQC_TEST_SEED` | Check repo secrets |
| go test / vet failure | Code regression | Fix in PR, do not force sign |

## Security Rules
- Never log private key material.
- Never commit production seeds.
- Local pre-commit performs conservative secret grep; CI should run gitleaks or trufflehog for depth.
- Public key only on the signer host; private material stays inside HSM / host-signer process.

## Next Big Thing Notes
Future seat expansion may add:
- Node.js / Python thin CLI wrappers around the Go binary for dashboard integration
- Webhook receiver that triggers this same pipeline
Keep the Go core authoritative.
