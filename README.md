# otel-demo

A Go-based Note App demonstrating a production-ready observability stack using **Prometheus**, **Loki**, and **OpenTelemetry**. Built with Gin, GORM, and PostgreSQL — containerized with Docker and deployable to Kubernetes.

---

## Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Web Framework | Gin |
| Database | PostgreSQL (via GORM) |
| Metrics | Prometheus |
| Logging | Logrus → Loki (HTTP push) |
| Tracing | OpenTelemetry (OTLP) |
| Container | Docker (multi-stage build) |
| Orchestration | Kubernetes |

---

## Project Structure

```
otel-demo/
├── controllers/        # Gin HTTP handlers (CRUD for notes)
├── models/             # GORM models
├── repository/         # DB layer
├── templates/          # HTML templates
├── k8s/                # Kubernetes manifests
├── main.go             # App entrypoint, middleware, routing
├── loki.go             # Custom Logrus → Loki HTTP hook
├── Dockerfile          # Multi-stage build
├── docker-compose.yaml # Local dev setup
└── Makefile            # Build, run, deploy commands
```

---

## Observability

### Metrics (Prometheus)
Exposed at `/metrics`. Two custom metrics are instrumented:

- `note_app_requests_total` — HTTP request counter labeled by `method`, `path`, `status`
- `note_app_request_duration_seconds` — Histogram of response latency labeled by `method`, `path`

### Logs (Loki)
Structured JSON logs via `logrus` are pushed directly to Loki using a custom HTTP hook in `loki.go` — no `promtail` or `loki-client-go` dependency required.

Each log entry includes `method`, `path`, `status`, and `duration` fields.

Configure the Loki endpoint via environment variable:
```bash
LOKI_URL=http://loki:3100/loki/api/v1/push
```

### Tracing (OpenTelemetry)
Traces are exported via OTLP HTTP to a configured endpoint (Tempo, Jaeger, etc.).

Configure via standard OTEL env vars:
```bash
OTEL_EXPORTER_OTLP_ENDPOINT=http://tempo:4318
OTEL_SERVICE_NAME=note-app
```

---

## Getting Started

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- kubectl (for Kubernetes deployment)

### Local Development

```bash
# Clone the repo
git clone https://github.com/omegazyadav/otel-demo.git
cd otel-demo

# Copy and configure env
cp .env.example .env

# Run locally
make run

# Or with Docker Compose
docker compose up
```

App runs at `http://localhost:8080`
Metrics at `http://localhost:8080/metrics`

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_DSN` | — | PostgreSQL DSN (required) |
| `LOKI_URL` | `http://loki:3100/loki/api/v1/push` | Loki push endpoint |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | — | OTLP trace exporter endpoint |
| `OTEL_SERVICE_NAME` | `note-app` | Service name in traces |
| `ENV` | — | Environment label attached to Loki logs |

---

## Makefile Commands

```bash
make build          # Build Go binary
make run            # Build and run locally
make docker-build   # Build Docker image
make docker-push    # Build and push to Docker Hub
make docker-run     # Run in Docker container
make deploy         # kubectl apply -f k8s/
make down           # kubectl delete -f k8s/
make clean          # Remove build artifacts
make help           # Show all commands
```

---

## Kubernetes Deployment

```bash
# Create DB secret (use SOPS in production — see Secrets section)
kubectl create secret generic note-app-db \
  --from-literal=dsn="host=<host> user=<user> password=<pass> dbname=notes_db port=5432 sslmode=disable"

# Deploy
make deploy

# Verify
kubectl get pods
kubectl logs -l app=note-app
```

---

## Secrets Management

Secrets are managed with [SOPS](https://github.com/getsops/sops) + AWS KMS so encrypted secrets can be safely committed to git.

```bash
# Install SOPS
brew install sops

# Encrypt secret
sops --encrypt secret.yaml > secret.enc.yaml

# Deploy (decrypt on the fly)
sops --decrypt secret.enc.yaml | kubectl apply -f -

# Edit encrypted secret
sops secret.enc.yaml
```

Configure your KMS key in `.sops.yaml`:
```yaml
creation_rules:
  - path_regex: .*secret.*\.yaml$
    kms: arn:aws:kms:<region>:<account-id>:key/<key-id>
```

> Never commit `secret.yaml`. Add it to `.gitignore`.

---

## API Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/` | List all notes |
| `POST` | `/notes` | Create a note |
| `GET` | `/notes/edit/:id` | Edit note form |
| `POST` | `/notes/update/:id` | Update a note |
| `POST` | `/notes/delete/:id` | Delete a note |
| `GET` | `/metrics` | Prometheus metrics |

---

## Roadmap

- [ ] Add OpenTelemetry tracing to DB layer via `otelgorm`
- [ ] Propagate trace context through controller → repository
- [ ] Add trace ID to Loki log entries for trace-log correlation
- [ ] Add full observability stack to `docker-compose.yaml` (Loki, Prometheus, Grafana, Tempo)
- [ ] Add External Secrets operator integration for k8s secret management

---

## License

MIT
