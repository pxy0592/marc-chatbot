# syntax=docker/dockerfile:1.12

FROM golang:1.26-alpine@sha256:3889b425f035be855a72fb4755265311293b6d414521f0a519d819df32222d83 AS build
WORKDIR /src

ARG GOPROXY=https://goproxy.cn|https://proxy.golang.org|direct
ENV GOPROXY="$GOPROXY"

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" \
    go build -trimpath -ldflags="-s -w" -o /out/marc-chatbot ./cmd/server

FROM alpine:3.22@sha256:14358309a308569c32bdc37e2e0e9694be33a9d99e68afb0f5ff33cc1f695dce
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -D -H -G app app \
    && mkdir -p /data \
    && chown app:app /data

COPY --from=build --chown=app:app /out/marc-chatbot /usr/local/bin/marc-chatbot

LABEL org.opencontainers.image.title="Marc Chatbot Backend" \
      org.opencontainers.image.description="Go backend for the Marc WeChat chatbot console"

ENV HTTP_ADDR=:8080 \
    DATABASE_PATH=/data/marc-chatbot.db \
    BOT_DRIVER=mock

USER app:app
VOLUME ["/data"]
EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/healthz >/dev/null || exit 1

ENTRYPOINT ["/usr/local/bin/marc-chatbot"]
