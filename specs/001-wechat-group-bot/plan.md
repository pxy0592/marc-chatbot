# Implementation Plan: 微信群聊机器人管理

**Branch**: `001-wechat-group-bot` | **Date**: 2026-08-18 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/001-wechat-group-bot/spec.md`

## Summary

构建一个前后端分离的微信机器人管理应用：Go 后端负责管理员鉴权、机器人生命周期、Wechaty Puppet Service 接入、群列表、群消息收发、SQLite 持久化、审计和事件流；TypeScript + React 前端提供绑定二维码、连接状态、群启停、消息列表、主动发送及审计视图。后端将 Wechaty 隔离在 `BotGateway` 适配器之后，同时提供 mock 适配器，使本地开发、自动化测试和无真实微信凭证的演示仍可完整运行。

## Technical Context

**Language/Version**: Go 1.26.x；TypeScript 5.x；Node.js 22 LTS

**Primary Dependencies**: Go `chi/v5`、`modernc.org/sqlite`、`go-wechaty v0.4.12`；React 19、Vite、QR code renderer、Vitest

**Storage**: 单节点 SQLite 数据库；通过迁移初始化 bot、group、message、outbound request、audit event 表

**Testing**: Go 标准 `testing` + `httptest`；前端 Vitest + Testing Library；端到端快速验证使用 mock bot 模式

**Target Platform**: Linux 容器或本地 Linux/macOS 开发环境；现代桌面浏览器

**Project Type**: 前后端分离 Web 应用 + 外部消息平台适配器

**Performance Goals**: 管理 API 在本地数据规模下 p95 小于 300ms；事件到达后 5 秒内可见；单实例支持至少 100 个群、30 天内 100,000 条文本消息

**Constraints**: 真实微信连接依赖可用的 Wechaty Puppet Service token/provider；敏感令牌仅通过环境变量注入；文本消息为 MVP；服务重启后业务状态可恢复；不将微信个人账号自动化能力误表述为微信官方开放能力

**Scale/Scope**: 单管理员域、单机器人实例、内部部署；一个管理页面；7 组核心 REST 接口 + 1 个 SSE 事件流

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

当前 `.specify/memory/constitution.md` 仍为未填写模板，没有已批准的项目原则可形成强制门禁。本计划采用以下保守门禁，均通过：

- **需求先行**: `spec.md` 已完成并通过 requirements checklist。
- **可替换外部适配器**: Wechaty 接入位于端口/适配器边界之后，核心业务不直接依赖外部 SDK 对象。
- **测试可独立运行**: mock bot 模式不需要真实微信账号、Puppet token 或公网服务。
- **最小权限与秘密隔离**: 管理端接口需要 bearer token；Puppet token 和管理员 token 不写入代码、数据库、日志或前端构建物。
- **可观察失败**: 连接、发送和审计状态必须持久化并通过接口可查询。
- **YAGNI**: MVP 不实现多租户、多机器人、富媒体、自动拉人入群或规则引擎。

**Post-design re-check**: 数据模型、接口契约和项目结构均维持单机器人、文本消息和适配器边界，没有新增门禁违反项。

## Project Structure

### Documentation (this feature)

```text
specs/001-wechat-group-bot/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── openapi.yaml
├── checklists/
│   └── requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
backend/
├── cmd/server/main.go
├── internal/
│   ├── api/             # HTTP handlers, auth, SSE and response mapping
│   ├── bot/             # BotGateway contract, mock adapter, Wechaty adapter
│   ├── config/          # environment configuration
│   ├── domain/          # entities and business status enums
│   ├── service/         # application use cases
│   └── store/           # SQLite repository and migrations
├── tests/
├── go.mod
└── .env.example

frontend/
├── src/
│   ├── api/             # typed fetch client and API DTOs
│   ├── components/      # status, binding, group, message and audit panels
│   ├── hooks/           # polling/SSE lifecycle hooks
│   ├── App.tsx
│   └── main.tsx
├── tests/
├── package.json
└── .env.example

scripts/
└── dev.sh
```

**Structure Decision**: 使用 `backend/` 和 `frontend/` 两个独立应用目录，REST/SSE 契约作为边界。后端内部使用端口/适配器分层，把第三方 Wechaty SDK 和业务状态机分开；前端保持单页管理台，避免在 MVP 引入不必要的路由和状态管理框架。

## Complexity Tracking

无需要例外说明的复杂度违反项。
