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

容器化功能的 SDD 产物位于 [`specs/002-containerized-deployment/`](specs/002-containerized-deployment/)：

- [`spec.md`](specs/002-containerized-deployment/spec.md)：容器运行、秘密配置和镜像交付需求
- [`research.md`](specs/002-containerized-deployment/research.md)：多阶段构建、同源代理、GHCR 和权限决策
- [`plan.md`](specs/002-containerized-deployment/plan.md)：Docker、Compose 与 CI 实施方案
- [`contracts/`](specs/002-containerized-deployment/contracts/)：Compose 和镜像标签契约
- [`tasks.md`](specs/002-containerized-deployment/tasks.md)：容器化实施任务
- [`quickstart.md`](specs/002-containerized-deployment/quickstart.md)：构建、运行、升级和验证说明

主题系统的 SDD 产物位于 [`specs/003-theme-system/`](specs/003-theme-system/)：

- [`spec.md`](specs/003-theme-system/spec.md)：白色/黑色主题、持久化和视觉替换需求
- [`research.md`](specs/003-theme-system/research.md)：主题上下文、系统偏好和新视觉方向
- [`plan.md`](specs/003-theme-system/plan.md)：React 主题基础设施和测试方案
- [`tasks.md`](specs/003-theme-system/tasks.md)：主题重构实施任务
- [`quickstart.md`](specs/003-theme-system/quickstart.md)：自动化与浏览器验证步骤

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

## Theme

管理台提供白色和黑色两种主题。登录页右上角及登录后的顶部导航均可切换主题；主动选择会保存在浏览器中，刷新后继续使用。首次访问且没有保存值时，系统采用设备的深浅色偏好。

主题偏好只保存在浏览器，不会发送到后端，也不会影响管理员令牌、机器人状态或消息数据。

## Docker deployment

前后端均提供多阶段、非 root Docker 镜像，并通过 Compose 以单一前端端口运行。SQLite 数据保存在 named volume 中，重建容器不会删除业务数据。

### Build and start

```bash
cp .env.docker.example .env
# Edit .env and replace ADMIN_TOKEN with a high-entropy value.

docker compose build
docker compose up -d --wait
docker compose ps
```

默认管理台地址：

```text
http://localhost:8080
```

常用操作：

```bash
docker compose logs -f backend frontend
docker compose restart backend
docker compose down
```

`docker compose down` 不会删除 SQLite volume。只有明确需要删除所有本地业务数据时才运行：

```bash
docker compose down --volumes
```

### Use published GHCR images

主分支 CI 会发布：

```text
ghcr.io/pxy0592/marc-chatbot-backend:latest
ghcr.io/pxy0592/marc-chatbot-frontend:latest
```

在 `.env` 中设置后可直接拉取运行：

```dotenv
BACKEND_IMAGE=ghcr.io/pxy0592/marc-chatbot-backend:latest
FRONTEND_IMAGE=ghcr.io/pxy0592/marc-chatbot-frontend:latest
```

```bash
docker compose pull
docker compose up -d --wait
```

Pull Request 只验证两个镜像能成功构建，不会推送镜像；提交到 `main` 或手动运行 CI 时会发布 `main`、`sha-*` 和 `latest` 标签。

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
- Pull Request 镜像构建验证，以及主分支的 GHCR 镜像发布。

## Scope and safety boundaries

- MVP 支持单机器人、少量管理员和微信群文本消息。
- 机器人加入微信群由微信客户端中的授权用户手动完成；系统不会绕过平台规则自动拉人或入群。
- 微信个人账号自动化能力受 Puppet provider、账号状态和平台规则影响，mock 验证不代表真实微信环境已验证。
- `ADMIN_TOKEN` 和 `WECHATY_PUPPET_SERVICE_TOKEN` 不会持久化到数据库或返回给前端；禁止提交 `.env`。
- 已发现群默认停用，只有管理员明确启用后消息才进入有效业务流。
