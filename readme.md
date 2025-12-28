# A Naive Datadog Autoscaler

> **Archived:** This repository is an archived example of a minimal Kubernetes autoscaler. It is not maintained and uses older Kubernetes APIs (`AppsV1beta1`).

A lightweight, unopinionated autoscaler that scales a specific Kubernetes Deployment based on the average value of a Datadog metric query.

It is designed to be run as a **Kubernetes CronJob** (e.g., every 5 minutes), performing a single check-and-scale operation before exiting.

## How it Works

1. **Fetches Metrics:** Queries Datadog for a specific metric over a defined period.
2. **Calculates Load:** Computes the average metric value per current replica.
   - `Load = AverageMetricValue / CurrentReplicas`
3. **Decides:** Compares `Load` against `MIN_WATERMARK` (scale down) and `MAX_WATERMARK` (scale up).
4. **Acts:** Scales the deployment up or down by **1 replica**, staying within `MIN_REPLICAS` and `MAX_REPLICAS`.

## Configuration

The application is configured entirely via Environment Variables.

### Target
| Variable | Description |
| :--- | :--- |
| `NAMESPACE` | Kubernetes namespace of the target deployment. |
| `DEPLOYMENT` | Name of the deployment to scale. |
| `LOCAL` | Set to `true` to use local `~/.kube/config`. Unset or `false` uses in-cluster service account. |

### Scaling Logic
| Variable | Description |
| :--- | :--- |
| `MIN_REPLICAS` | Minimum allowed number of replicas. |
| `MAX_REPLICAS` | Maximum allowed number of replicas. |
| `MIN_WATERMARK` | Threshold to scale **down** (if `Load < MIN`). |
| `MAX_WATERMARK` | Threshold to scale **up** (if `Load > MAX`). |

### Datadog
| Variable | Description |
| :--- | :--- |
| `DD_API_KEY` | Your Datadog API Key. |
| `DD_APP_KEY` | Your Datadog Application Key. |
| `DD_URI` | Datadog API URL (e.g., `https://api.datadoghq.com/api/v1/query`). |
| `DD_QUERY` | The metric query to execute (e.g., `avg:system.cpu.idle{service:my-service}`). |
| `DD_QUERY_PERIOD` | Time window in seconds to average metrics over (e.g., `300`). |

## Usage

### Local Development

1. **Initialize and download dependencies:**
   ```bash
   go mod tidy
   ```

2. **Build:**
   ```bash
   go build
   ```

3. **Run:**
   ```bash
   # Example configuration
   export LOCAL=true
   export NAMESPACE=default
   export DEPLOYMENT=nginx-frontend
   export MIN_REPLICAS=2
   export MAX_REPLICAS=10
   export MIN_WATERMARK=40
   export MAX_WATERMARK=80
   
   # Datadog settings
   export DD_API_KEY=...
   export DD_APP_KEY=...
   export DD_URI="https://api.datadoghq.com/api/v1/query"
   export DD_QUERY="avg:my.custom.metric{env:prod}"
   export DD_QUERY_PERIOD=300

   ./autoscaler
   ```

### Kubernetes (CronJob)

Deploy as a CronJob to run at a set interval.

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: datadog-autoscaler
spec:
  schedule: "*/5 * * * *" # Run every 5 minutes
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccountName: autoscaler-sa # Needs permission to get/update Deployments
          containers:
          - name: autoscaler
            image: your-repo/autoscaler:latest
            env:
              - name: NAMESPACE
                value: "production"
              - name: DEPLOYMENT
                value: "worker-pool"
              - name: MIN_REPLICAS
                value: "3"
              - name: MAX_REPLICAS
                value: "20"
              # ... Include all required env vars ...
          restartPolicy: OnFailure
```
