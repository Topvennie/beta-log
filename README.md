# Beta-log

A web application focussed around bouldering.

## Features

### OIDC

- Supports any OIDC provider

### Workouts

- Create custom workouts

### Boulder App Integration

- Fetch and save data from third party boulder apps used by your climbing gym. Currently supports:
  - Toplogger

## Production Deployment

### Recommended Deployment (Docker)

1. Copy `docker-compose.prod.yml` -> `docker-compose.yml`.
2. Copy `.env.prod.example` -> `.env`.
3. Fill in the `.env` file (see later section).
4. Run `docker compose up -d`.
5. The server is reachable on port **3000**.

To update:

```bash
docker compose pull
docker compose down
docker compose up -d
```

### Manual Setup (Advanced)

The container image is published as: `ghcr.io/topvennie/beta-log`.

Required additional services:

- **Postgres**

Configuration variables can be overridden via environment variables.
Keys use uppercase and replace `.` with `_` (e.g. `db.host` -> `DB_HOST`)

Search the codebase for `config.Get` to view available configuration settings.
You will probably need to change a couple to use your own services.

### Environment Variables

TODO
