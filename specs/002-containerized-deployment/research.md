# Research: 容器化部署与镜像交付

**Date**: 2026-08-18

## Multi-stage and runtime users

**Decision**: 两个 Dockerfile 均使用 build/runtime 分阶段；后端运行时为 Alpine 非 root 用户，前端为 Nginx 非 root 用户并监听 8080。

**Rationale**: 构建工具和开发依赖不会进入最终镜像；高端口无需网络绑定 capability；Alpine 保留健康检查所需的轻量工具。

## Frontend-to-backend routing

**Decision**: 前端默认使用同源 API 地址，Nginx 把 `/api/` 和 `/healthz` 代理到 Compose 内部服务名 `backend:8080`，并关闭 API 响应缓冲以支持 SSE。

**Rationale**: 浏览器不能解析内部容器主机名；同源入口避免生产 CORS 和编译期后端地址绑定。

## Persistence

**Decision**: 后端数据库固定为 `/data/marc-chatbot.db`，Compose named volume 挂载 `/data`。

**Rationale**: 容器可替换但业务数据必须独立持久化；SQLite 仍限制为单后端实例。

## Image workflow

**Decision**: 扩展现有 `ci.yml`。PR Job 只有 `contents: read` 并执行 build-only；主分支/手动 Job 增加 `packages: write`，登录 GHCR 并推送 backend/frontend 镜像。

**Rationale**: 分离验证与发布权限，避免 PR 获得不需要的包写权限，也避免分叉 PR 使用秘密或尝试发布。

## Tags and metadata

**Decision**: 发布 `main`、`sha-<short commit>` 和默认分支上的 `latest` 标签，并写入 OCI source/revision/title/description labels。

**Rationale**: `latest` 便于默认部署，`main` 表示分支渠道，SHA 标签提供不可歧义回滚和追踪。

## Build security

**Decision**: Docker 官方 actions 固定完整 SHA；checkout 不持久化凭据；BuildKit 使用 GitHub Actions cache；启用 SBOM 和 provenance。

**Verified action releases on 2026-08-18**:

- setup-buildx-action v4.2.0
- build-push-action v7.3.0
- login-action v4.6.0
- metadata-action v6.2.0

## Alternatives considered

- 单一镜像同时运行前后端：进程监督、日志和升级边界更复杂，不采用。
- 前端编译固定后端公网 URL：部署环境耦合且需要重新构建，不采用。
- 主分支仅构建不发布：无法形成可直接拉取的制品，不采用。
- PR 和发布共用 packages write：不符合最小权限，不采用。
