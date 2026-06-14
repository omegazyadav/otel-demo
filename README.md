# otel-demo

Observability stack using **Prometheus**, **Loki**, and **OpenTelemetry**.

---

## Stack

| Layer | Technology |
|---|---|
| Metrics | Prometheus |
| Logging | Logrus → Loki (HTTP push) |
| Tracing | OpenTelemetry (OTLP) |
| Container | Docker (multi-stage build) |
| Orchestration | Kubernetes |

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

## Kubernetes Deployment

```bash
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

## License

MIT
