# futureAppointmentScheduler

A small Go API for trainer scheduling. It supports:
- listing appointments by trainer
- creating a new appointment
- finding trainer availability in a time range

The app uses PostgreSQL, runs migration/seed on startup, and exposes HTTP on port `8080` by default.

## Run the Project

### Docker Compose

```bash
docker compose up --build
```

Stop services:

```bash
docker compose down
```

Stop services and remove DB volume data:

```bash
docker compose down -v
```

### Makefile shortcuts

```bash
make up
make down
make down-volumes
make test
```

## Endpoints

- `GET /appointments?trainer_id=<id>`
- `POST /appointments`
- `GET /availability?trainer_id=<id>&starts_at=<rfc3339>&ends_at=<rfc3339>`

## Sample curl commands

### 1) Get appointments for a trainer

```bash
curl "http://localhost:8080/appointments?trainer_id=1"
```

### 2) Create a valid appointment (expected `201`)

```bash
curl -X POST "http://localhost:8080/appointments" \
  -H "Content-Type: application/json" \
  -d '{
    "trainer_id": 1,
    "user_id": 42,
    "starts_at": "2019-01-24T11:00:00-08:00",
    "ends_at": "2019-01-24T11:30:00-08:00"
  }'
```

### 3) Create an overlapping appointment (expected `409`)

This reuses an existing seeded slot for trainer `1`.

```bash
curl -X POST "http://localhost:8080/appointments" \
  -H "Content-Type: application/json" \
  -d '{
    "trainer_id": 1,
    "user_id": 99,
    "starts_at": "2019-01-24T09:00:00-08:00",
    "ends_at": "2019-01-24T09:30:00-08:00"
  }'
```

### 4) Get trainer availability in a range

```bash
curl "http://localhost:8080/availability?trainer_id=1&starts_at=2019-01-24T08:00:00-08:00&ends_at=2019-01-24T17:00:00-08:00"
```

### 5) Method not allowed checks (expected `405`)

```bash
curl -X PUT "http://localhost:8080/appointments"
curl -X POST "http://localhost:8080/availability"
```

### 6) Missing/invalid parameter examples (expected `400`)

```bash
curl "http://localhost:8080/appointments"
curl "http://localhost:8080/availability?trainer_id=1&starts_at=bad&ends_at=bad"
```

## Overlapping appointment challenge (solved)

The key challenge is concurrent booking: two requests could otherwise insert overlapping slots at the same time.

This project solves that in PostgreSQL with an exclusion constraint in `migrations/001_init.sql`:
- constraint name: `appointments_no_overlap`
- scope: same `trainer_id`
- rule: time ranges cannot overlap (`tstzrange(starts_at, ends_at, '[)')`)

When Postgres rejects an overlapping insert, the repository maps that DB error to `ErrAppointmentOverlap`, and the API returns `409 Conflict`.

## Project layout

```text
futureAppointmentScheduler/
|-- cmd/
|   `-- api/
|       `-- main.go
|-- internal/
|   |-- appointments/
|   |   |-- handler.go
|   |   |-- service.go
|   |   |-- repository.go
|   |   `-- models.go
|   `-- db/
|       |-- postgres.go
|       `-- seed.go
|-- migrations/
|   `-- 001_init.sql
|-- data/
|   `-- appointments.json
|-- docker-compose.yml
|-- Dockerfile
`-- Makefile
```
