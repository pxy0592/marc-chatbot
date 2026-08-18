# Research: 微信群聊机器人管理

**Date**: 2026-08-18

## 1. Wechaty Go 接入方式

**Decision**: 后端使用官方 `wechaty/go-wechaty` SDK 的稳定版本 `v0.4.12`，通过 Wechaty Puppet Service / gateway 连接实际微信能力，并把 SDK 封装为 `BotGateway`。

**Rationale**:

- 用户明确要求 Go 后端和 Wechaty。
- 官方 Go SDK 支持 scan、login、logout、message 事件，以及 room 查询和 `Room.Say` 文本发送。
- 官方 Go SDK 文档说明，多语言 SDK 需要 Puppet Service，Go 进程通过 gRPC 服务使用实际 puppet provider。
- 适配器边界可控制 SDK 版本、provider 差异、重连以及未来替换风险。

**Alternatives considered**:

- TypeScript Wechaty 独立进程 + Go 业务服务：生态更成熟，但会让“后端使用 Go”变成双后端运行时，增加部署和跨进程状态一致性复杂度。
- 直接调用非官方微信 Web 协议：维护和合规风险过高，不采用。
- 只实现 mock：无法满足真实绑定和群消息目标，仅用于本地和测试。

**Sources**:

- Wechaty 官方中文文档：`https://wechaty.js.org/zh/docs/`
- Go Wechaty 官方仓库：`https://github.com/wechaty/go-wechaty`
- Puppet Service DIY 文档：`https://wechaty.js.org/docs/puppet-services/diy/`
- Puppet provider 列表：`https://wechaty.js.org/docs/puppet-providers/`

## 2. SDK 版本策略

**Decision**: 首版固定 `go-wechaty v0.4.12`，不直接依赖 `master` 未发布提交。

**Rationale**:

- `v0.4.12` 是官方仓库可见的最新稳定 release，发布于 2024-04-01。
- 官方 `master` 在 2025-05-18 仍有提交，但未形成更新 release；生产依赖未发布提交会降低可复现性。
- 当前目标仅需要 v0.4.12 已具备的扫码、登录、消息、群查询和群发文本能力。

**Alternatives considered**:

- 固定 `master` commit：可获得更新功能，但升级风险和兼容验证成本更高。
- 浮动 latest：构建不可复现，不采用。

## 3. 连接生命周期

**Decision**: 后端进程拥有唯一 bot runtime。生命周期状态为 `unbound -> awaiting_scan -> connecting -> online -> reconnecting/offline/error`，scan/login/logout/error 事件更新持久化快照并推送 SSE。

**Rationale**:

- 单实例符合 MVP 范围。
- 明确状态机可以防止“SDK 进程仍在重连但界面显示在线”的假状态。
- SSE 足以支持单向状态与消息通知，比 WebSocket 更简单，控制操作仍走 REST。

**Alternatives considered**:

- 前端固定轮询：实现简单，但连接状态与二维码过期感知较慢；保留低频轮询作为 SSE 断线兜底。
- WebSocket：双向能力当前没有必要。

## 4. 持久化

**Decision**: 使用 SQLite WAL 模式和显式迁移，持久化群启用状态、消息、主动发送、bot 状态和审计事件。

**Rationale**:

- 单实例内部部署规模适合 SQLite，部署成本低。
- WAL 提供足够的读写并发。
- 使用 repository 接口，后续可迁移 PostgreSQL 而不改变服务层。

**Alternatives considered**:

- 仅内存：重启后无法满足状态恢复和审计要求。
- PostgreSQL：首版运维成本高于收益。

## 5. 管理员鉴权

**Decision**: MVP 使用部署时注入的高熵 bearer token；前端只在浏览器 `sessionStorage` 保存，不写入静态构建或后端数据库。

**Rationale**:

- 适合内部单管理员域，满足接口保护和快速部署。
- 所有 `/api/v1` 业务接口和 `/events` 均鉴权。
- 可在后续替换为企业 SSO，而不改变资源模型。

**Alternatives considered**:

- 自建用户名密码：需要账户、密码重置、加密与会话治理，超出 MVP。
- 无鉴权：不能接受。

## 6. 幂等与消息去重

**Decision**: 接收消息以 provider message ID 唯一去重；主动发送要求客户端提供或由服务端生成 idempotency key，并对同一 key 返回已有结果。

**Rationale**:

- 外部事件和管理请求都可能重复投递。
- 唯一约束比纯内存去重更能跨重启工作。

## 7. 前端实现

**Decision**: React + TypeScript + Vite 单页管理台，使用原生 fetch 和 EventSource 兼容方案；二维码在浏览器本地渲染，不把扫码字符串提交给第三方二维码网站。

**Rationale**:

- 页面规模有限，不需要额外全局状态框架。
- 本地 QR 渲染减少绑定数据外发。
- 响应式桌面优先布局满足内部管理需求。

## 8. 可测试性与开发模式

**Decision**: `BOT_DRIVER=mock` 为默认开发模式，支持生成模拟二维码、切换在线、注入群消息和模拟发送；`BOT_DRIVER=wechaty` 才加载真实连接配置。

**Rationale**:

- CI 和开发环境通常没有微信账号或 Puppet token。
- mock 必须经过同一 service/store/API 路径，避免测试一套、生产另一套业务逻辑。

## 9. 风险与边界

- 微信个人账号自动化能力取决于 Puppet provider、账号状态及平台规则；项目不能保证任何账号长期可用。
- “添加机器人入群”在 MVP 中表示用户通过微信客户端将已绑定账号手动加入群，管理台负责发现和启用，不实现绕过平台规则的自动入群。
- Go SDK 的稳定 release 节奏较慢，因此所有 SDK 交互必须集中在适配器目录并有 mock contract tests。
