# ScopePay

AI-powered freelance project payments. Base skeleton: Go/Gin/GORM backend + Vite React TypeScript web app.

## Layout

```
backend/                       Go API (Gin + GORM)
  cmd/server/main.go           entrypoint
  internal/config/             env loading
  internal/db/                 gorm open + auto-migrate (sqlite/postgres)
  internal/models/             User, Project, Milestone, Payment, WebhookEvent
  internal/router/             routes (/health) + CORS
web/                           Vite + React + TypeScript SPA
  src/main.tsx                 router (react-router)
  src/pages/Home.tsx           calls GET /health, shows result
  src/api/client.ts            typed fetch client
  vite.config.ts               dev proxy /api -> localhost:8080
```

## Backend

Requires Go 1.24+ and a C compiler (cgo, for the SQLite driver).

```sh
cd backend
go run ./cmd/server
```

Defaults: `PORT=8080`, `DB_DRIVER=sqlite`, `DB_DSN=scope_pay.db`. See `../.env.example` for all vars.

```sh
curl localhost:8080/health
# {"db":"sqlite","status":"ok"}
```

Tables are created on first run via AutoMigrate.

`CGO_ENABLED=0` will break SQLite (`go-sqlite3` is a C binding). Keep it `1` (Go's default when `gcc` is present).

### Postgres

```sh
DB_DRIVER=postgres \
DB_DSN="host=localhost user=scopepay password=scopepay dbname=scopepay port=5432 sslmode=disable" \
go run ./cmd/server
```

## Web

```sh
cd web
npm install
npm run dev          # http://localhost:5173
```

Dev proxy forwards `/api/*` to `http://localhost:8080`, so no CORS or URL config needed locally. For a deployed build, set `VITE_API_URL` to the API origin at build time.

```sh
npx tsc -b           # typecheck
npm run build        # production bundle in dist/
```

Backend allows browser origin `http://localhost:5173` by default; override with `CORS_ORIGINS` (comma-separated).

## Not built yet

Intentionally out of scope for this base setup: auth, AI intake, PayPal integration, proposal/milestone UI, tests.
