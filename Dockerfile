# First stage: Get Golang image from DockerHub.
FROM golang:1.26.4-alpine3.23 AS backend-builder

# Set our working directory for this stage.
WORKDIR /app

# Copy all of our files.
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./migrate.go

FROM node:26.3.0-alpine3.23 AS base-frontend
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

# corepack being fucked as usual
RUN npm i -g corepack@latest
RUN corepack enable
COPY ./ui /frontend/ui
WORKDIR /frontend

FROM base-frontend AS frontend-prod-deps
WORKDIR /frontend/ui
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM base-frontend AS frontend-build

WORKDIR /frontend/ui

ARG APP_VERSION
ENV VITE_APP_VERSION=$APP_VERSION

ENV CI=true
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

# Last stage: discard everything except our executables.
FROM alpine:3.23 AS prod

# Set our next working directory.
WORKDIR /app

# Copy our executable and our built React application.
COPY --from=backend-builder /app/server .
COPY --from=backend-builder /app/migrate .
COPY --from=frontend-build /frontend/public ./public
COPY ./config ./config

ENV APP_ENV=production

# Declare entrypoints and activation commands.
EXPOSE 3000
ENTRYPOINT ["/bin/sh", "-c", "./migrate && ./server"]
