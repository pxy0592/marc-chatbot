# Marc Chatbot

基于 Wechaty 的微信群聊机器人管理台。项目采用前后端分离架构：

- **后端**：Go，负责 Wechaty 生命周期、群消息收发、群授权、SQLite 持久化、审计和 REST/SSE API。
- **前端**：TypeScript + React + Vite，提供绑定二维码、连接状态、群启停、消息流、主动发送和审计界面。
- **开发流程**：使用 Spec Kit / SDD，需求、研究、计划、数据模型、API 契约、任务与实现保持可追踪。

> Wechaty 中文官方文档：`https://wechaty.js.org/zh/docs/`

## SDD artifacts

本功能的 SDD 产物位于 [`specs/001-wechat-group-bot/`](specs/001-wechat-group-bot/)：

- [`spec.md`](specs/001-wechat-group-bot/spec.md)：用户故事、需求和验收标准
- [`research.md`](specs/001-wechat-group-bot/research.md)：Wechaty、Puppet Service、存储与鉴权决策
- [`plan.md`](specs/001-wechat-group-bot/plan.md)：技术方案和项目结构
- [`data-model.md`](specs/001-wechat-group-bot/data-model.md)：实体、约束和状态机
- [`contracts/openapi.yaml`](specs/001-wechat-group-bot/contracts/openapi.yaml)：HTTP/SSE 契约
- [`tasks.md`](specs/001-wechat-group-bot/tasks.md)：依赖有序实施任务
- [`quickstart.md`](specs/001-wechat-group-bot/quickstart.md)：mock 与真实 Wechaty 验证流程

## Architecture

```text
React admin console
        │ REST + authenticated SSE
        ▼
Go HTTP API ── Application services ── SQLite
                         │
                         ▼
                  BotGateway port
                    ├─ mock
                    └─ go-wechaty v0.4.12
                           │
                           ▼
                  Wechaty Puppet Service
```

核心业务依赖 `BotGateway` 抽象，不直接暴露 Wechaty SDK 类型。默认 `mock` 模式可在无微信账号和 Puppet token 的情况下运行测试和完整演示。

## Local development

### Prerequisites

- Go 1.26.x
- Node.js 22+
- npm 10+

### Install frontend dependencies

```bash
cd frontend
npm install
```

### Run both applications in mock mode

```bash
ADMIN_TOKEN=dev-admin-token BOT_DRIVER=mock ./scripts/dev.sh
```

打开 `http://localhost:5173`，使用 `dev-admin-token` 登录。管理台提供 mock 登录、离线和消息注入操作。

也可以分别启动：

```bash
cd backend
ADMIN_TOKEN=dev-admin-token BOT_DRIVER=mock go run ./cmd/server

cd frontend
npm run dev
```

## Real Wechaty mode

Go Wechaty 需要 Puppet Service/provider。先阅读：

- `https://wechaty.js.org/docs/puppet-services/diy/`
- `https://wechaty.js.org/docs/puppet-providers/`

然后仅通过环境变量注入真实配置：

```bash
cd backend
BOT_DRIVER=wechaty \
ADMIN_TOKEN='<high-entropy-admin-token>' \
WECHATY_PUPPET_SERVICE_TOKEN='<provider-token>' \
go run ./cmd/server
```

前端点击“开始绑定”后，后端才启动实际 Wechaty/Puppet 连接并接收扫码事件。

## Tests

```bash
cd backend
go test ./...

cd ../frontend
npm test -- --run
npm run lint
npm run build
```

GitHub Actions 工作流位于 `.github/workflows/ci.yml`。它会在提交到 `main`、针对 `main` 的 Pull Request 以及手动触发时，并行执行：

- 后端模块校验、格式检查、`go vet`、竞态检测测试和服务端构建。
- 前端锁文件安装、ESLint、Vitest 和生产构建。

## Scope and safety boundaries

- MVP 支持单机器人、少量管理员和微信群文本消息。
- 机器人加入微信群由微信客户端中的授权用户手动完成；系统不会绕过平台规则自动拉人或入群。
- 微信个人账号自动化能力受 Puppet provider、账号状态和平台规则影响，mock 验证不代表真实微信环境已验证。
- `ADMIN_TOKEN` 和 `WECHATY_PUPPET_SERVICE_TOKEN` 不会持久化到数据库或返回给前端；禁止提交 `.env`。
- 已发现群默认停用，只有管理员明确启用后消息才进入有效业务流。
