# Cascade OJ: Copilot Instructions

## Repo map and boundaries
- Monorepo with backend in `cascade-oj/` and frontend in `cascade-oj-frontend/`.
- Backend runtime services: `gateway`, `public(auth)`, `user`, `judge`, `admin`.
- Service entrypoints are in `app/**/cmd/*/main.go`; DI wiring is generated in `wire_gen.go` (do not edit generated wire files directly).
- API contracts live in `api/cascade/**.proto`; HTTP routes are defined by `google.api.http` options in proto, then generated via `make api`.

## Actual data flow (important for feature work)
- User-facing reads/writes mostly go through `gateway` (`app/gateway/internal/service/*`) and then to `user` via gRPC client (`app/gateway/internal/data/data.go`).
- Judging is async: `gateway`/`user` publish to RabbitMQ queues `gojudge-submission-queue` and `gojudge-self-test-queue` (see `pkg/mq/channel.go`).
- `judge` service subscribes those queues in `app/services/judge/internal/server/mq.go`, runs go-judge, then writes results back to Redis/MySQL.
- Token header key is custom `token` (not `Authorization`); frontend request wrappers in `packages/*/src/**/request*.ts` already follow this.

## Build, run, and debug workflows
- Backend expects Unix-like shell even on Windows (use Git Bash): run commands from `cascade-oj/cascade-oj/`.
- Typical backend flow:
  1. `make init` (tooling)
  2. `docker compose -f docker-compose.dependencies.yml up -d` (or `./deploy.sh`)
  3. `make config api ent wire build`
  4. `make db-init` (after MySQL is up)
- Deploy script `deploy.sh` is the standard operations entry (supports `--build`, `--pull`, `--down`, service-scoped deploy).
- Frontend uses pnpm workspace from `cascade-oj-frontend/`: `pnpm install`, then `pnpm run competition:dev` / `admin:dev` / `login:dev`.
- Frontend package engines require Node `^20.19.0 || >=22.12.0` (see package `engines` fields).

## Project-specific conventions to preserve
- Prefer editing source-of-truth, then regenerate:
  - Ent: edit `ent/schema/*.go`, then `make ent`.
  - Proto/API: edit `api/**/*.proto`, then `make api` (and `make openapi` if API docs changed).
  - Wire: edit provider sets, then `make wire`.
- Avoid manual edits in generated files under `ent/` (non-schema), `app/**/cmd/**/wire_gen.go`, and generated API code.
- Auth guard is enforced per service HTTP server:
  - `admin` uses `auth.RoleAdmin` globally.
  - `user` protects selected methods via selector path list.
  - `public` auth service has no auth middleware.

## Integration points and config expectations
- Infrastructure defaults come from compose + `config/*/config.yaml`: MySQL, Redis, RabbitMQ, go-judge.
- Service hostnames in config assume docker network names (e.g. `cascade_mysql`, `cascade_rabbitmq`, `cascade_go_judge`).
- Gateway self-registers endpoint to `user.Register`; if networking changes, update `config/gateway/config.yaml` (`endpoint_preset`, `api_key`, `ipv4_prefix`, `port`).

{
  "builder": {
    "gc": {
      "defaultKeepStorage": "20GB",
      "enabled": true
    }
  },
  "experimental": false,
  "registry-mirrors": [
    "https://docker.m.daocloud.io",
    "https://hub.rat.dev"
  ]
}