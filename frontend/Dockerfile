# syntax=docker/dockerfile:1.12

FROM node:22-alpine@sha256:c610fcdfb1d5b4740dd70c284ed3cb16bb857e0f7166196e36a5501df7a3aa32 AS build
WORKDIR /app

COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

COPY . .
ARG VITE_API_BASE_URL=
ENV VITE_API_BASE_URL="$VITE_API_BASE_URL"
RUN npm run build

FROM nginx:1.29-alpine@sha256:5616878291a2eed594aee8db4dade5878cf7edcb475e59193904b198d9b830de
RUN rm -f /etc/nginx/conf.d/default.conf \
    && mkdir -p /tmp/client_temp /tmp/proxy_temp /tmp/fastcgi_temp /tmp/uwsgi_temp /tmp/scgi_temp \
    && chown -R nginx:nginx /tmp /usr/share/nginx/html

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=build --chown=nginx:nginx /app/dist /usr/share/nginx/html

LABEL org.opencontainers.image.title="Marc Chatbot Frontend" \
      org.opencontainers.image.description="React and Nginx frontend for the Marc WeChat chatbot console"

USER nginx:nginx
EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/ >/dev/null || exit 1

ENTRYPOINT []
CMD ["nginx", "-g", "daemon off;"]
