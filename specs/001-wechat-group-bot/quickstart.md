# Quickstart Validation: 微信群聊机器人管理

## Prerequisites

- Go 1.26.x
- Node.js 22 LTS and npm
- 本地可写目录用于 SQLite
- 真实模式额外需要有效的 Wechaty Puppet Service token/provider

## 1. Mock mode end-to-end

### Start backend

```bash
cd backend
cp .env.example .env
BOT_DRIVER=mock ADMIN_TOKEN=dev-admin-token go run ./cmd/server
```

Expected:

- `GET http://localhost:8080/healthz` returns `{"status":"ok"}`.
- SQLite 数据库和表自动创建。
- 业务接口未带 `Authorization: Bearer dev-admin-token` 时返回 401。

### Start frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Open the Vite URL shown in the terminal. Enter `dev-admin-token` in the login panel.

### Validate primary flows

1. Click **开始绑定**. A mock QR code and `awaiting_scan` status appear.
2. Use **模拟登录** development action. Bot status becomes `online`.
3. Confirm mock groups appear and are disabled by default.
4. Enable one group.
5. Inject a mock text message into that group; it appears once in the message list.
6. Send a text message from the group panel; status becomes `succeeded` and an outbound message record appears.
7. Re-submit using the same idempotency key through the API; no second group message is created.
8. Disable the group and verify a new injected message is stored as ignored or excluded from the active message stream according to the contract.
9. Open audit history and verify binding, group toggle, and send events are present without secrets.

## 2. Automated checks

```bash
cd backend
go test ./...

golangci-lint run ./... # optional when installed

cd ../frontend
npm test -- --run
npm run build
```

## 3. Real Wechaty mode

Review the official Wechaty documentation and select a Puppet provider compatible with the target account and deployment:

- `https://wechaty.js.org/zh/docs/`
- `https://wechaty.js.org/docs/puppet-services/diy/`
- `https://wechaty.js.org/docs/puppet-providers/`

Configure environment variables without committing values:

```bash
cd backend
BOT_DRIVER=wechaty \
ADMIN_TOKEN='<high-entropy-admin-token>' \
WECHATY_PUPPET_SERVICE_TOKEN='<provider-token>' \
go run ./cmd/server
```

Expected validation:

1. Start binding and scan the locally rendered QR code in the admin UI.
2. Bot transitions to `online`; account display name is visible but no credential is returned.
3. Manually add the bound account to a test group using WeChat.
4. Refresh/sync groups and enable the test group.
5. Send a text message from another group member and verify it appears within 5 seconds.
6. Send a text message from the admin UI and verify it arrives in the target group.
7. Stop or disconnect the provider and verify the UI does not continue to report `online`.

## 4. Security checks

- Search application output and SQLite data for the real admin/Puppet tokens; neither value should appear.
- Verify CORS only allows configured origins.
- Verify empty, oversized, offline, unavailable-group and disabled-group send attempts fail clearly.
- Verify SSE requires authentication and reconnects without disclosing the token in server logs.

## Validation record — 2026-08-18

- Backend unit/contract tests: PASS (`go test ./...`).
- Frontend component tests: PASS (4/4, `npm test -- --run`).
- Frontend lint: PASS (`npm run lint`).
- Production build: PASS (`npm run build`).
- Mock API quickstart: PASS — health, 401 protection, binding QR state, mock login, two discovered groups, group enable, inbound injection, idempotent outbound send, two message directions, four required audit actions, and token non-disclosure were verified.
- Browser validation: PASS — authenticated login, binding, mock login, group enable, inbound injection, outbound send and audit rendering were exercised in Chromium. A toggle overlay pointer-event issue discovered during browser testing was fixed.
- Real Wechaty/Puppet environment: NOT RUN — requires a valid provider token and an authorized WeChat test account; the adapter is compile-validated only.
