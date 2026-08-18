# Tasks: 微信群聊机器人管理

**Input**: Design documents from `specs/001-wechat-group-bot/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: 后端服务和 HTTP 合同必须有自动化测试；前端关键交互必须有组件测试；mock 模式用于无真实微信凭证的验收。

**Organization**: 任务按用户故事组织，每个故事均包含可独立验证的纵向切片。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可在不同文件上并行执行
- **[Story]**: 对应 spec.md 的用户故事
- 每个任务包含明确文件路径

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 初始化前后端工程、开发命令和安全忽略规则。

- [X] T001 创建 `backend/`、`frontend/`、`scripts/` 目录结构并扩展根 `.gitignore`
- [X] T002 [P] 初始化 Go 模块和依赖声明于 `backend/go.mod`
- [X] T003 [P] 初始化 React + TypeScript + Vite 工程于 `frontend/package.json`、`frontend/tsconfig*.json`、`frontend/vite.config.ts`
- [X] T004 [P] 创建后端与前端配置示例 `backend/.env.example`、`frontend/.env.example`
- [X] T005 [P] 创建本地联合启动脚本 `scripts/dev.sh`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成所有用户故事共用的领域模型、数据库、配置、鉴权、事件和 API 骨架。

**⚠️ CRITICAL**: 本阶段完成前不开始用户故事功能。

- [X] T006 [P] 定义领域实体和状态枚举于 `backend/internal/domain/models.go`
- [X] T007 [P] 实现环境配置与校验于 `backend/internal/config/config.go`
- [X] T008 实现 SQLite 迁移和 repository 于 `backend/internal/store/sqlite.go`
- [X] T009 [P] 定义 `BotGateway` 端口和事件类型于 `backend/internal/bot/gateway.go`
- [X] T010 [P] 实现内存事件总线于 `backend/internal/service/events.go`
- [X] T011 实现应用服务基础结构和审计辅助方法于 `backend/internal/service/service.go`
- [X] T012 [P] 实现 JSON 错误、Bearer 鉴权、CORS 和 request logging 于 `backend/internal/api/middleware.go`
- [X] T013 建立 HTTP router、health endpoint 和依赖装配于 `backend/internal/api/router.go`、`backend/cmd/server/main.go`
- [X] T014 [P] 创建前端 API 类型与鉴权 fetch client 于 `frontend/src/api/types.ts`、`frontend/src/api/client.ts`
- [X] T015 [P] 创建前端基础主题和页面框架于 `frontend/src/index.css`、`frontend/src/App.tsx`

**Checkpoint**: 基础设施可启动，health 可访问，业务 API 已受鉴权保护。

---

## Phase 3: User Story 1 - 绑定并保持机器人在线 (Priority: P1) 🎯 MVP

**Goal**: 管理员可发起绑定、查看二维码/状态、取消绑定，并获得准确连接状态。

**Independent Test**: mock 模式下完成开始绑定、模拟登录、模拟离线和重新绑定，状态与审计均正确。

### Tests for User Story 1

- [X] T016 [P] [US1] 编写 bot 状态和绑定服务测试于 `backend/internal/service/binding_test.go`
- [X] T017 [P] [US1] 编写 bot/binding HTTP 合同测试于 `backend/internal/api/binding_test.go`
- [X] T018 [P] [US1] 编写绑定面板组件测试于 `frontend/src/components/BindingPanel.test.tsx`

### Implementation for User Story 1

- [X] T019 [US1] 实现 mock bot 生命周期、二维码和事件注入于 `backend/internal/bot/mock.go`
- [X] T020 [US1] 实现绑定与状态应用服务于 `backend/internal/service/binding.go`
- [X] T021 [US1] 实现 `/api/v1/bot` 和 binding endpoints 于 `backend/internal/api/binding.go`
- [X] T022 [US1] 实现 SSE 事件 endpoint 于 `backend/internal/api/events.go`
- [X] T023 [US1] 实现前端登录、状态和绑定面板于 `frontend/src/components/AuthGate.tsx`、`frontend/src/components/BindingPanel.tsx`
- [X] T024 [US1] 将绑定状态、轮询和 SSE 刷新接入 `frontend/src/App.tsx`、`frontend/src/hooks/useBotEvents.ts`

**Checkpoint**: 无真实 Wechaty token 时也可独立演示完整绑定状态流。

---

## Phase 4: User Story 2 - 查看并启用目标群聊 (Priority: P1)

**Goal**: 发现机器人群聊并逐群启用或停用消息接入。

**Independent Test**: mock 登录后看到默认群，启用和停用可跨服务重启持久化。

### Tests for User Story 2

- [X] T025 [P] [US2] 编写群 repository 与服务测试于 `backend/internal/service/groups_test.go`
- [X] T026 [P] [US2] 编写群 API 合同测试于 `backend/internal/api/groups_test.go`
- [X] T027 [P] [US2] 编写群列表组件测试于 `frontend/src/components/GroupPanel.test.tsx`

### Implementation for User Story 2

- [X] T028 [US2] 实现群发现、同步和启停服务于 `backend/internal/service/groups.go`
- [X] T029 [US2] 实现 `/api/v1/groups` endpoints 于 `backend/internal/api/groups.go`
- [X] T030 [US2] 实现群列表、状态标识和启停交互于 `frontend/src/components/GroupPanel.tsx`

**Checkpoint**: 群发现和授权边界可独立验证。

---

## Phase 5: User Story 3 - 接收和查看群消息 (Priority: P1)

**Goal**: 接收已启用群的文本消息，去重、标记自身消息并在管理台查询。

**Independent Test**: 向启用群注入两次同一 provider message ID，仅出现一条记录；自身消息不触发循环。

### Tests for User Story 3

- [X] T031 [P] [US3] 编写入站消息去重和过滤测试于 `backend/internal/service/messages_test.go`
- [X] T032 [P] [US3] 编写消息查询 API 测试于 `backend/internal/api/messages_test.go`
- [X] T033 [P] [US3] 编写消息列表组件测试于 `frontend/src/components/MessagePanel.test.tsx`

### Implementation for User Story 3

- [X] T034 [US3] 实现 bot 事件消费和入站消息服务于 `backend/internal/service/messages.go`
- [X] T035 [US3] 实现 `/api/v1/messages` endpoint 和 mock 消息开发 endpoint 于 `backend/internal/api/messages.go`
- [X] T036 [US3] 实现消息筛选和列表 UI 于 `frontend/src/components/MessagePanel.tsx`

**Checkpoint**: 入站文本消息在 5 秒目标内可见并可查询。

---

## Phase 6: User Story 4 - 从管理界面向群里发消息 (Priority: P2)

**Goal**: 对在线且已启用群进行幂等文本发送并呈现最终结果。

**Independent Test**: 相同 Idempotency-Key 重复请求只调用一次 gateway；离线和停用群发送明确失败。

### Tests for User Story 4

- [X] T037 [P] [US4] 编写主动发送幂等和失败状态测试于 `backend/internal/service/outbound_test.go`
- [X] T038 [P] [US4] 编写群发送 API 合同测试于 `backend/internal/api/outbound_test.go`
- [X] T039 [P] [US4] 编写发送表单组件测试于 `frontend/src/components/SendMessageForm.test.tsx`

### Implementation for User Story 4

- [X] T040 [US4] 实现主动发送状态机和幂等服务于 `backend/internal/service/outbound.go`
- [X] T041 [US4] 实现 `/api/v1/groups/{groupId}/messages` endpoint 于 `backend/internal/api/outbound.go`
- [X] T042 [US4] 实现群消息发送表单和结果反馈于 `frontend/src/components/SendMessageForm.tsx`

**Checkpoint**: 管理台和目标群形成可靠文本消息闭环。

---

## Phase 7: User Story 5 - 审计机器人操作 (Priority: P3)

**Goal**: 查询关键绑定、群启停和主动发送审计事件，不泄露秘密。

**Independent Test**: 执行三类操作后均有审计记录，搜索数据库和 API 响应不包含配置 token。

### Tests for User Story 5

- [X] T043 [P] [US5] 编写审计查询和敏感信息排除测试于 `backend/internal/service/audit_test.go`
- [X] T044 [P] [US5] 编写审计 API 合同测试于 `backend/internal/api/audit_test.go`

### Implementation for User Story 5

- [X] T045 [US5] 实现审计查询服务和 `/api/v1/audit-events` endpoint 于 `backend/internal/service/audit.go`、`backend/internal/api/audit.go`
- [X] T046 [US5] 实现审计面板于 `frontend/src/components/AuditPanel.tsx`

**Checkpoint**: 关键操作具有可查询且非敏感的审计链。

---

## Phase 8: Real Wechaty Adapter & Polish

**Purpose**: 接入真实 provider，完成安全、文档、构建和端到端验证。

- [X] T047 实现 `go-wechaty v0.4.12` 真实适配器于 `backend/internal/bot/wechaty.go`
- [X] T048 将 `BOT_DRIVER=mock|wechaty` 选择接入 `backend/cmd/server/main.go`
- [X] T049 [P] 添加后端集成测试工具和 API 测试夹具于 `backend/internal/testutil/testapp.go`
- [X] T050 [P] 完善前端响应式视觉、空状态和错误状态于 `frontend/src/index.css`
- [X] T051 [P] 更新根 `README.md`，链接 SDD 产物、架构、启动和真实 Wechaty 风险说明
- [X] T052 运行 `gofmt`、`go test ./...`、`npm test -- --run`、`npm run build` 并修复问题
- [X] T053 按 `specs/001-wechat-group-bot/quickstart.md` 执行 mock 模式验证并记录结果

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 Setup**: 无依赖。
- **Phase 2 Foundational**: 依赖 Phase 1，阻塞所有用户故事。
- **Phase 3 US1**: 依赖 Phase 2；提供 bot 生命周期。
- **Phase 4 US2**: 依赖 Phase 2；完整在线群发现验收依赖 US1。
- **Phase 5 US3**: 依赖 US2 的群启用边界。
- **Phase 6 US4**: 依赖 US1 在线状态和 US2 群启用边界。
- **Phase 7 US5**: 审计基础由 Phase 2 提供，最终验收依赖 US1/US2/US4 事件。
- **Phase 8**: 依赖所有目标故事完成。

### User Story Dependencies

```text
Foundation
├── US1 Binding/Status ─┬─> US4 Outbound
│                      └─> real adapter validation
└── US2 Group Scope ────┬─> US3 Inbound
                        └─> US4 Outbound
US1 + US2 + US4 ──────────> US5 audit acceptance
```

### Within Each User Story

1. 先写测试并确认在缺少实现时失败。
2. 领域/repository 能力先于 service。
3. service 先于 HTTP handler。
4. API 合同稳定后接前端组件。
5. 在 checkpoint 运行该故事相关测试。

### Parallel Opportunities

- T002-T005 可并行。
- T006、T007、T009、T010 可在不同文件并行。
- 每个故事的 Go service test、API contract test、React component test 可并行编写。
- US1 完成后，US2 的前端与 US1 的视觉优化可并行。
- T049-T051 可并行。

## Implementation Strategy

### MVP First

1. 完成 Setup + Foundation。
2. 完成 US1 绑定/状态。
3. 完成 US2 群授权。
4. 完成 US3 入站消息。
5. 使用 mock 模式演示最小可用接收链路。

### Incremental Delivery

- Increment 1: mock 绑定与在线状态。
- Increment 2: 群发现和启用。
- Increment 3: 入站消息可见。
- Increment 4: 主动发送闭环。
- Increment 5: 审计与真实 Wechaty adapter。

## Notes

- 所有秘密只通过环境变量配置，禁止写入 fixture、截图或文档示例真实值。
- `go-wechaty` 的 provider 能力和账号稳定性必须在真实环境单独验证；mock 测试不能被表述为真实微信验证。
- 每完成一个任务，将其勾选为 `[X]`。
