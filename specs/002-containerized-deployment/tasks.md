# Tasks: 容器化部署与镜像交付

**Input**: `specs/002-containerized-deployment/` design artifacts

## Phase 1: Container build foundation

- [X] T001 [P] Create backend build context exclusions in `backend/.dockerignore`
- [X] T002 [P] Create frontend build context exclusions in `frontend/.dockerignore`
- [X] T003 [P] Change frontend production API default to same-origin in `frontend/src/api/client.ts`

## Phase 2: User Story 1 - One-command stack (P1)

- [X] T004 [US1] Create multi-stage non-root backend image in `backend/Dockerfile`
- [X] T005 [US1] Create multi-stage non-root frontend image in `frontend/Dockerfile`
- [X] T006 [US1] Create static and reverse-proxy runtime config in `frontend/nginx.conf`
- [X] T007 [US1] Create root Compose stack in `docker-compose.yml`
- [X] T008 [US1] Create safe Docker environment template in `.env.docker.example`
- [X] T009 [US1] Validate Compose interpolation, dependency health, named volume and one-port exposure

## Phase 3: User Story 2 - Secure runtime configuration (P1)

- [X] T010 [US2] Configure backend read-only root filesystem, non-root user, dropped capabilities and writable data volume
- [X] T011 [US2] Configure frontend read-only root filesystem, non-root user, dropped capabilities and tmpfs
- [X] T012 [US2] Add image health checks for backend and frontend
- [X] T013 [US2] Verify missing and invalid runtime secrets fail without being embedded in image layers

## Phase 4: User Story 3 - Automated image delivery (P2)

- [X] T014 [P] [US3] Add PR image build validation job in `.github/workflows/ci.yml`
- [X] T015 [P] [US3] Add main/manual GHCR publication job in `.github/workflows/ci.yml`
- [X] T016 [US3] Add branch, SHA and latest image metadata and OCI labels
- [X] T017 [US3] Pin Docker actions to full immutable commit SHAs and configure BuildKit cache
- [X] T018 [US3] Enable image SBOM and provenance generation
- [X] T019 [US3] Validate Workflow with actionlint

## Phase 5: User Story 4 - Same-origin management console (P2)

- [X] T020 [US4] Proxy `/api/` and `/healthz` to backend in `frontend/nginx.conf`
- [X] T021 [US4] Disable proxy buffering and extend read timeout for SSE
- [X] T022 [US4] Configure SPA fallback and static asset caching
- [X] T023 [US4] Exercise login, health and API access through the frontend container port

## Phase 6: Documentation and final validation

- [X] T024 [P] Document Docker build, start, logs, stop, upgrade and GHCR behavior in `README.md`
- [X] T025 [P] Record container validation commands and boundaries in `specs/002-containerized-deployment/quickstart.md`
- [X] T026 Run existing backend and frontend test/build checks
- [X] T027 Build both images locally and inspect non-root runtime users
- [X] T028 Start Compose stack, wait for health, run browser/API smoke tests and verify SQLite persistence
- [X] T029 Run `docker compose down` without deleting the named volume and verify clean shutdown

## Dependencies

- T001-T003 can run in parallel.
- T004-T008 depend on build foundation.
- T010-T013 refine T004-T008.
- T014-T019 depend on Dockerfiles.
- T020-T023 depend on frontend image and Compose networking.
- T024-T029 complete after implementation.
