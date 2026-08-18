# Quickstart: Docker deployment

## Configure

```bash
cp .env.docker.example .env
```

Set a high-entropy `ADMIN_TOKEN`. Keep `.env` untracked.

## Build and start

```bash
docker compose build
docker compose up -d
docker compose ps
```

Open:

```text
http://localhost:8080
```

## Operate

```bash
docker compose logs -f backend frontend
docker compose restart backend
docker compose down
docker compose down --volumes  # destructive: also removes SQLite data
```

## Upgrade

```bash
git pull
docker compose build --pull
docker compose up -d --remove-orphans
```

## Real Wechaty mode

Set in `.env`:

```dotenv
BOT_DRIVER=wechaty
WECHATY_PUPPET_SERVICE_TOKEN=<provider-token>
```

Then recreate the backend:

```bash
docker compose up -d --force-recreate backend frontend
```

## Validation

```bash
ADMIN_TOKEN=test-token docker compose config --quiet
docker build -t marc-chatbot-backend:test backend
docker build -t marc-chatbot-frontend:test frontend
```

Verify runtime users are non-root, Compose services become healthy, and the UI/API work through the frontend port.

## Validation record — 2026-08-18

- `docker compose config --quiet`: PASS with required runtime settings.
- Backend image build: PASS; static Go binary built with `CGO_ENABLED=0`.
- Frontend image build: PASS; TypeScript and Vite production build completed in the image stage.
- Runtime users: PASS — backend `app:app`, frontend `nginx:nginx`.
- Container hardening: PASS — read-only root filesystems, all capabilities dropped, no-new-privileges enabled.
- Health checks: PASS — both Compose services reached `healthy` within the configured wait timeout.
- Same-origin smoke test: PASS — frontend HTML, proxied `/healthz`, authenticated API and SSE events all worked through one published port.
- Persistence test: PASS — enabled group state survived backend restart and full Compose container removal/recreation while the named volume was retained.
- Missing-secret tests: PASS — Compose interpolation and direct backend container startup failed without `ADMIN_TOKEN`.
- Existing Go tests, race detector, vet, frontend lint/tests/build: PASS.
- GitHub Actions syntax: PASS with actionlint v1.7.12.
- GHCR publication: NOT RUN locally — publication requires the committed workflow to execute in GitHub Actions with repository package permissions.
