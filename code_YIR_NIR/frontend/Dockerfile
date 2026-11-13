# Сборка SPA (monorepo: shared + frontend)
FROM node:22-alpine AS builder

WORKDIR /app

RUN corepack enable && corepack prepare pnpm@9.15.0 --activate

COPY package.json pnpm-workspace.yaml pnpm-lock.yaml ./
COPY shared/package.json ./shared/
COPY frontend/package.json ./frontend/

RUN pnpm install --frozen-lockfile

COPY shared ./shared
COPY frontend ./frontend

RUN pnpm --filter @sociometry/shared build
RUN pnpm --filter frontend build

# Статика + reverse proxy к backend
FROM nginx:1.27-alpine

COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/frontend/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
