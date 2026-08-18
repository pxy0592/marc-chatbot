# Compose Contract

- `docker compose up --build -d` starts exactly one backend and one frontend service.
- Only frontend port is published by default.
- Backend must be healthy before frontend dependency is considered satisfied.
- `ADMIN_TOKEN` is mandatory at interpolation/start time.
- `BOT_DRIVER` defaults to `mock`.
- `WECHATY_PUPPET_SERVICE_TOKEN` is optional unless `BOT_DRIVER=wechaty`.
- SQLite data is stored in a named volume mounted at `/data`.
- Both services run with all Linux capabilities dropped and no-new-privileges.
- Both root filesystems are read-only; required temporary paths use tmpfs.
