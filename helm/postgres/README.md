# otel-app-chart

Helm chart for the otel-demo app — a Go HTTP service with OpenTelemetry tracing, metrics, and Prometheus scraping built in.

## Add the repo

```bash
helm repo add otel-app https://rbalman.github.io/otel-app
helm repo update
```

## Install

```bash
helm install demo otel-app/otel-app-chart
```

### With Ingress enabled

```bash
helm install demo otel-app/otel-app-chart \
  --set ingress.enabled=true \
  --set ingress.className=nginx \
  --set ingress.host=demo.example.com
```

### With HPA enabled

```bash
helm install demo otel-app/otel-app-chart \
  --set hpa.enabled=true \
  --set hpa.maxReplicas=5
```

## Configuration

| Key | Default | Description |
|-----|---------|-------------|
| `replicaCount` | `1` | Number of pod replicas |
| `image.repository` | `balman/otel-demo` | Container image repository |
| `image.tag` | `v0.0.3` | Image tag |
| `image.pullPolicy` | `IfNotPresent` | Image pull policy |
| `service.port` | `8080` | HTTP port |
| `resources.requests.cpu` | `50m` | CPU request |
| `resources.requests.memory` | `64Mi` | Memory request |
| `resources.limits.cpu` | `200m` | CPU limit |
| `resources.limits.memory` | `128Mi` | Memory limit |
| `extraEnvs` | see values.yaml | Extra environment variables injected into the container |
| `podAnnotations` | Prometheus scrape annotations | Annotations added to the pod template |
| `ingress.enabled` | `false` | Enable Ingress resource |
| `ingress.className` | `""` | Ingress class (e.g. `nginx`) |
| `ingress.host` | `demo.local` | Ingress hostname |
| `ingress.path` | `/` | Ingress path |
| `ingress.pathType` | `Prefix` | Ingress path type |
| `ingress.annotations` | `{}` | Annotations on the Ingress resource |
| `ingress.tls` | `[]` | TLS configuration for the Ingress |
| `hpa.enabled` | `false` | Enable HorizontalPodAutoscaler |
| `hpa.minReplicas` | `1` | HPA minimum replicas |
| `hpa.maxReplicas` | `5` | HPA maximum replicas |
| `hpa.targetCPUUtilizationPercentage` | `70` | CPU utilisation target (%) |
| `hpa.targetMemoryUtilizationPercentage` | `80` | Memory utilisation target (%) |

## Probes

All three Kubernetes probe types are configured out of the box:

| Probe | Path | Default timing | Notes |
|-------|------|----------------|-------|
| Startup | `/health` | 15 × 2 s | 30 s budget for slow starts |
| Liveness | `/health` | every 10 s | Restarts the pod if the process hangs |
| Readiness | `/ready` | every 5 s | Removes the pod from the Service endpoint slice |

## Prometheus scraping

The pod is annotated for scraping by default:

```yaml
podAnnotations:
  prometheus.io/scrape: "true"
  prometheus.io/path: /metrics
  prometheus.io/port: "8080"
```

## Prerequisites

- Kubernetes 1.21+
- Helm 3.x
- An OpenTelemetry Collector reachable at the configured OTLP endpoint
