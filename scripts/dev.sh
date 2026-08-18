#!/usr/bin/env bash
set -euo pipefail
trap 'kill 0' EXIT
(
  cd backend
  BOT_DRIVER="${BOT_DRIVER:-mock}" ADMIN_TOKEN="${ADMIN_TOKEN:-dev-admin-token}" go run ./cmd/server
) &
(
  cd frontend
  npm run dev
) &
wait
