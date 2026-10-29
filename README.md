# Field Service

Manages mini-soccer fields, time slots, and per-date schedules (availability +
price). The frontend reads schedules here before creating an order.

## Ports and dependencies

| Dependency | Address | Notes |
|---|---|---|
| This service | `http://localhost:8003` | `config.json` -> `port` |
| PostgreSQL | `localhost:5432` | database `go-field-service` (tables auto-migrated on start) |
| user-service | `http://localhost:8001` | token/role checks for admin routes |
| MinIO | `localhost:9100` | field images, bucket `field-images` (container `minio`) |

Local Postgres is the native Windows PostgreSQL 17 service on `5432`; the
`postgres` container in `docker-compose.yml` publishes no host port, so
`go run .` on the host uses the native service. MinIO runs as the `minio`
container (`docker start minio`).

## Prerequisites

- Go >= version in `go.mod` (auto-downloaded via `GOTOOLCHAIN=auto`, the default).
- PostgreSQL reachable at `localhost:5432`.
- Docker for the MinIO container.

## Setup

```bash
cp config.json.example config.json   # fill database.*, signatureKey, internalService.user.signatureKey
go mod tidy
```

`config.json` and `.env` hold real credentials and are gitignored - never commit them.

Inter-service signature headers (any client calling this service):

```
x-service-name: <caller appName>
x-request-at:   <unix seconds>
x-api-key:      sha256("{x-service-name}:{signatureKey}:{x-request-at}")
```

## Run

```bash
go run .
# or hot reload: make watch-prepare (first time), make watch
```

Start Postgres + MinIO first; on start this service migrates its tables and
ensures the `field-images` bucket exists.

## API

Public (signature headers only, no user token):

| Method | Path | Body / Query |
|---|---|---|
| GET | `/api/v1/field` | - |
| GET | `/api/v1/field/:uuid` | - |
| GET | `/api/v1/field/schedule/lists/:uuid` | `?date=YYYY-MM-DD` (field UUID) |

Authenticated (`Authorization: Bearer <token>` + signature headers):

| Method | Path | Role | Body |
|---|---|---|---|
| GET | `/api/v1/field/pagination` | admin, customer | - |
| POST | `/api/v1/field` | admin | multipart (name, location, pricePerHour, image) |
| PUT | `/api/v1/field/:uuid` | admin | - |
| DELETE | `/api/v1/field/:uuid` | admin | - |
| POST | `/api/v1/field/schedule` | admin | `{"fieldID":1,"date":"YYYY-MM-DD","timeID":"<time-uuid>"}` |
| POST | `/api/v1/field/schedule/one-month` | admin | `{"fieldID":1}` |
| PUT | `/api/v1/field/schedule/:uuid` | admin | - |
| PATCH | `/api/v1/field/schedule` | admin | - |
| DELETE | `/api/v1/field/schedule/:uuid` | admin | - |

Schedules are created by an admin (single slot via `POST /api/v1/field/schedule`
or a month of slots via `POST /api/v1/field/schedule/one-month`); the booking
flow only consumes existing rows, so a seeded schedule per date is enough for
the frontend.

## Verify locally

```bash
TS=$(date +%s); KEY=$(printf "field-service:%s:%s" "<signatureKey>" "$TS" | sha256sum | cut -d' ' -f1)

# list fields (signature headers only)
curl "http://localhost:8003/api/v1/field" \
  -H "x-service-name: field-service" -H "x-request-at: $TS" -H "x-api-key: $KEY"

# schedules for a field on a date (what the booking page calls)
curl "http://localhost:8003/api/v1/field/schedule/lists/<field-uuid>?date=$(date +%F)" \
  -H "x-service-name: field-service" -H "x-request-at: $TS" -H "x-api-key: $KEY"
```

## Docker

```bash
docker compose up -d --build --force-recreate   # image dansss/field-service, port 8003
```
