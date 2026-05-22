#!/usr/bin/env bash
# Enforce minimum statement coverage for internal/game (see docs/balance.md).
set -euo pipefail

MIN_COVERAGE="${MIN_GAME_COVERAGE:-80}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PROFILE="${1:-$ROOT/coverage.game.out}"

cd "$ROOT"
go test ./internal/game/... -covermode=atomic -coverprofile="$PROFILE" >/dev/null

total="$(go tool cover -func="$PROFILE" | awk '/^total:/ { gsub(/%/, "", $3); print $3 }')"
if [[ -z "$total" ]]; then
  echo "check-game-coverage: could not read total from $PROFILE" >&2
  exit 1
fi

echo "internal/game coverage: ${total}% (minimum ${MIN_COVERAGE}%)"
awk -v got="$total" -v min="$MIN_COVERAGE" 'BEGIN { exit !(got + 0 >= min + 0) }' || {
  echo "check-game-coverage: coverage ${total}% is below ${MIN_COVERAGE}%" >&2
  exit 1
}
