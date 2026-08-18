# Implementation Plan: 容器化部署与镜像交付

**Branch**: `002-containerized-deployment` | **Date**: 2026-08-18 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/002-containerized-deployment/spec.md`

## Summary

为 Go 后端和 React 前端分别创建多阶段、非 root 的 Docker 镜像；使用 Compose 提供单入口、后端健康依赖、SQLite 持久化和外部秘密配置。前端运行时使用 Nginx 提供静态文件并把 `/api/`、`/healthz` 和 SSE 请求转发到内部后端服务。新增独立镜像工作流：Pull Request 只构建验证，主分支和手动运行构建并发布两个镜像到 GitHub Container Registry。

## Technical Context

**Language/Version**: Go 1.26；Node.js 22；Dockerfile frontend syntax；Compose Specification

**Primary Dependencies**: Docker Engine/BuildKit、Docker Compose v2、Alpine Linux、Nginx、GitHub Actions Docker 官方 actions

**Storage**: Compose named volume 挂载到后端 `/data`，SQLite 路径 `/data/marc-chatbot.db`

**Testing**: 既有 Go/Vitest 检查；`docker build`；`docker compose config`；容器非 root/health/HTTP smoke test；actionlint

**Target Platform**: Linux AMD64 容器；GitHub-hosted Ubuntu runner

**Project Type**: 前后端 Web 应用的容器化部署与制品发布

**Performance Goals**: 本地完整栈 60 秒内健康；CI 两个镜像 15 分钟内完成；静态前端代理 SSE 不缓冲

**Constraints**: 镜像不包含秘密；运行时非 root；PR 不推送；主分支发布 GHCR；SQLite 单后端副本；TLS 外置

**Scale/Scope**: 2 个镜像、1 个 Compose stack、现有 CI 中的镜像验证与发布 Jobs、Linux AMD64

## Constitution Check

当前 constitution 仍是模板。本功能采用并通过以下门禁：

- SDD 规格和质量清单先于实现。
- Dockerfile 使用多阶段构建和最小运行时内容。
- 最终进程非 root，Compose 删除 capabilities 并启用 no-new-privileges。
- 秘密只从运行环境注入，不使用 Docker build args 传入。
- PR 构建只验证；只有非 PR 发布 Job 获得 packages write。
- 外部 GitHub Actions 固定到完整提交 SHA。
- 真实运行与静态配置分开报告；只有实际 image/compose smoke test 通过才声明容器验证完成。

Post-design check: Dockerfile、Compose、反向代理和 Workflow 合同均满足门禁，无例外。

## Project Structure

### Documentation

```text
specs/002-containerized-deployment/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── compose.md
│   └── image-tags.md
├── checklists/requirements.md
└── tasks.md
```

### Source and deployment files

```text
backend/
├── Dockerfile
└── .dockerignore
frontend/
├── Dockerfile
├── .dockerignore
└── nginx.conf
docker-compose.yml
.env.docker.example
.github/workflows/ci.yml
README.md
```

**Structure Decision**: 每个应用以自身目录作为独立构建上下文，避免把另一应用或仓库元数据发送给 BuildKit。Compose 位于仓库根目录并仅公开前端端口；镜像工作流独立于现有代码 CI。
