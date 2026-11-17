#!/usr/bin/env bash
# Install Operator Desk pure-shell pre-commit hooks
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
HOOKS_DIR="$ROOT/.githooks"

if [[ ! -d "$ROOT/.git" ]]; then
  echo "Not a git repository. Initializing..."
  git init
fi

git config core.hooksPath "$HOOKS_DIR"
chmod +x "$HOOKS_DIR/pre-commit"

echo "Hooks installed. core.hooksPath = $HOOKS_DIR"
echo "Pure shell only. Do not add .pre-commit-config.yaml unless that rule is retired in writing."
echo "Run a test commit to verify."
