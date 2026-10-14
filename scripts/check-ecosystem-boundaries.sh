#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# Mux is host-neutral: it must not depend on any Decoder package.
# Shared ecosystem vocabulary lives in decoder/internal/contracts, which
# hosts vendor rather than import from here.
FORBIDDEN_DECODER='github\.com/GrayCodeAI/decoder(/|")'

if command -v rg >/dev/null 2>&1; then
  violations="$(rg -n "$FORBIDDEN_DECODER" --glob '*.go' . || true)"
else
  violations="$(grep -rn --include='*.go' -E "$FORBIDDEN_DECODER" . || true)"
fi

if [[ -n "${violations}" ]]; then
  echo "forbidden Decoder host imports found:"
  echo "${violations}"
  echo
  echo "mux must use local contracts, never the Decoder product module"
  exit 1
fi

echo "ecosystem boundary guard passed"