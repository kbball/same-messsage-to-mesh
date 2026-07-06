# SAME → Mesh

![SAME → Mesh](frontend/public/logo.png)

Decodes SAME/EAS weather alerts from an RTL-SDR dongle and publishes plain-text alerts to an MQTT broker for rebroadcast over a Meshcore ham radio mesh network. Designed for ARES operators who want automated weather alert distribution across a mesh.

---

## Features

- **Live alert feed** — decoded SAME messages appear in real-time via Server-Sent Events
- **Configurable filters** — select states, counties, and event types to act on
- **SDR configuration** — set device path and NOAA Weather Radio frequency via the UI
- **Reference data** — FIPS county codes and SAME event codes refreshable from NOAA/Census
- **MQTT publishing** — Phase 2; publishes plain-text alerts (e.g. `[TOR] TORNADO WARNING - Fulton County GA - Until 14:30 (KFFC/NWS)`) to a configurable topic

---

## Requirements

- RTL-SDR USB dongle (for live decoding)
- Docker + docker-compose
- Go 1.22+ and Node 20+ (development only)

---

## Quick Start

```bash
cp .env.example .env
# Edit .env — set DB_PASSWORD at minimum
make db-up      # starts Postgres + Mosquitto
make dev        # starts backend (port 8080) + frontend (port 5173)
```

Open [http://localhost:5173](http://localhost:5173).

On first run, go to **Reference Data** and click **Refresh from NOAA** to populate FIPS county codes and event codes.

---

## Configuration

All configuration is via environment variables. Copy `.env.example` to `.env` for local development.

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `8080` | HTTP API port |
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `DB_HOST` | `localhost` | Postgres host |
| `DB_PORT` | `5432` | Postgres port |
| `DB_NAME` | `same_mesh` | Database name |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | — | **Required.** Database password |
| `DB_SSL_MODE` | `disable` | Postgres SSL mode |
| `SDR_DEVICE_PATH` | `/dev/bus/usb` | RTL-SDR USB device path |
| `SDR_FREQUENCY` | `162550000` | NOAA Weather Radio frequency in Hz |
| `MQTT_ENABLED` | `false` | Set `true` to enable MQTT publishing (Phase 2) |
| `MQTT_HOST` | `localhost` | MQTT broker host |
| `MQTT_PORT` | `1883` | MQTT broker port |
| `MQTT_PUBLISH_TOPIC` | `same/alerts` | MQTT topic for alert messages |

---

## SDR Setup

1. Plug in your RTL-SDR dongle.
2. In `docker-compose.yml`, uncomment the `devices:` block under the `app` service and adjust the path if needed:
   ```yaml
   devices:
     - /dev/bus/usb:/dev/bus/usb
   ```
3. Set `SDR_DEVICE_PATH` and `SDR_FREQUENCY` via the **SDR Config** tab in the UI (or in `.env`).
4. NOAA Weather Radio frequencies: 162.400 / 162.425 / 162.450 / 162.475 / 162.500 / 162.525 / **162.550 MHz** (most common).

---

## Architecture

Hexagonal (Ports & Adapters) — dependency direction is always inward:

```
adapter → application → domain
```

| Layer | Path | Responsibility |
|---|---|---|
| Domain | `backend/internal/domain/` | Entities and port interfaces. Pure Go, no frameworks. |
| Application | `backend/internal/application/` | Business logic. Orchestrates domain via interfaces. |
| Adapter | `backend/internal/adapter/` | HTTP, Postgres, SDR, MQTT, NOAA, SSE. |

**Stack:** Go · Postgres · goose migrations · React · TypeScript · Material UI · Vite · Docker

---

## Development

```bash
make help        # list all targets
make test        # run all tests
make coverage    # tests with coverage report
make lint        # golangci-lint + ESLint
make fmt         # gofmt + Prettier
make migrate-create NAME=add_something   # scaffold a new migration
```

Pre-commit hook runs `make fmt && make lint` automatically after `make install-hooks`.

---

## Deployment

An operator-only compose file pulls the pre-built image from GHCR:

```bash
docker compose -f docker-compose.operator.yml up -d
```

Image: `ghcr.io/kbball/same-message-to-mesh:latest`

Releases are tagged via the manual GitHub Actions workflow (`.github/workflows/release.yml`).
