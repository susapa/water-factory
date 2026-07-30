#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

if [ ! -f .env ]; then
  echo "Error: .env not found." >&2
  echo "Copy .env.production.example to .env and fill in real values first:" >&2
  echo "  cp .env.production.example .env" >&2
  exit 1
fi

echo "==> Building images"
docker compose build

echo "==> Starting services"
docker compose up -d

echo "==> Service status"
docker compose ps

echo
echo "Deploy complete. Frontend: http://<vps-ip>/  Backend health: http://<vps-ip>:8080/health"
echo "Tail logs with: docker compose logs -f backend"
