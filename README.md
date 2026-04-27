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

## TODO

### 🔴 High Priority

- [ ] **Add OpenTelemetry tracing** — repo is called `otel-demo` but has no `initTracerProvider`, spans, or `otelgin` middleware
  ```go
  // go get go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin
  r.Use(otelgin.Middleware("note-app"))
  ```
- [ ] **Propagate `ctx`** from Gin handler → controller → repository so DB spans are nested under HTTP spans
- [ ] **Add `otelgorm` plugin** for automatic DB query tracing
  ```go
  db.Use(otelgorm.NewPlugin())
  ```
- [ ] **Make `loki.go` Fire() async** — currently blocks the request path if Loki is slow or unavailable
  ```go
  func (h *LokiHook) Fire(entry *logrus.Entry) error {
      go func() { /* push logic */ }()
      return nil
  }
  ```

---

### 🟡 Medium Priority

- [ ] **Add trace ID + span ID to Loki log entries** for trace-log correlation in Grafana
- [ ] **Add `level` label to Loki streams** so logs can be filtered by severity in Grafana
- [ ] **Add full observability stack to `docker-compose.yaml`** — Loki, Prometheus, Grafana, Tempo are missing
- [ ] **Add DB healthcheck in `docker-compose.yaml`** so app waits for Postgres to be ready
  ```yaml
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U postgres"]
  ```
- [ ] **Remove mixed `log` and `logrus` usage** — standardize on `logrus` throughout `main.go`
- [ ] **Add `.env.example`** — currently there is no reference file for required environment variables
- [ ] **Add `activeRequests` gauge metric** — exists in the simple demo but missing in note-app

---

### 🟢 Low Priority / Nice to Have

- [ ] **Add SOPS + `.sops.yaml`** setup for encrypted secret management (see Secrets section)
- [ ] **Add External Secrets manifest** for Kubernetes secret management via AWS SSM
- [ ] **Add Grafana dashboard JSON** for the Prometheus metrics already instrumented
- [ ] **Add sampling configuration** via `OTEL_TRACES_SAMPLER` env var for production use
  ```bash
  OTEL_TRACES_SAMPLER=parentbased_traceidratio
  OTEL_TRACES_SAMPLER_ARG=0.1
  ```
- [ ] **Write `prometheus.yml` scrape config** for local docker-compose usage
- [ ] **Add GitHub Actions CI workflow** for build + docker push on merge to main

---

## License

MIT
