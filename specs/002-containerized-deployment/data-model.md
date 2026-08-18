# Deployment Model: 容器化部署与镜像交付

## BackendImage

- Build context: `backend/`
- Runtime command: `/usr/local/bin/marc-chatbot`
- Runtime user: non-root `app`
- Internal port: `8080`
- Writable path: `/data`
- Health endpoint: `/healthz`
- Required runtime setting: `ADMIN_TOKEN`

## FrontendImage

- Build context: `frontend/`
- Runtime process: Nginx non-root
- Internal port: `8080`
- Static root: `/usr/share/nginx/html`
- Proxy target: `backend:8080`
- Writable runtime paths: `/tmp` only

## ApplicationStack

```text
Browser -> frontend:8080 -> / static files
                       -> /api/* -> backend:8080
                       -> /healthz -> backend:8080
backend -> /data/marc-chatbot.db -> named volume
```

## PublishedImage

- Registry: GitHub Container Registry
- Names: `<repository>-backend`, `<repository>-frontend`
- Tags: `main`, `sha-<commit>`, `latest` on default branch
- Platform: `linux/amd64`
- Metadata: OCI source, revision, title, description; BuildKit provenance and SBOM
