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

#### Authentication

Create OIDC credentials in your identity provider.
For example if you're using [Authentik](https://goauthentik.io/) you can find the instructions [here](https://docs.goauthentik.io/add-secure-apps/providers/oauth2/create-oauth2-provider/).

The callback URL is `<custom_domain>/api/auth/callback/openid-connect`.

#### Data

- `DB_DIR` - Location of your database content.
- `LOG_DIR` - Location of the log directory.
