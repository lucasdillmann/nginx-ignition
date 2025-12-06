# Health checks

nginx ignition provides health check endpoints that can be used by container orchestration platforms (like Docker)
and monitoring systems to verify the application's status and availability.

## Available endpoints

### Liveness endpoint

The liveness endpoint checks if the application is running and healthy. When called, ignition checks if its internal
components are working as expected (like the connection with the database).

**Endpoint:** `GET /api/health/liveness`

**Response codes:**
- `200 OK`: Application is healthy and all components are functioning correctly
- `503 Service Unavailable`: Application is running but one or more critical components are unhealthy

**Example response (healthy):**
```json
{
  "healthy": true,
  "components": [
    {
      "name": "database",
      "healthy": true
    }
  ]
}
```

**Example response (unhealthy):**
```json
{
  "healthy": false,
  "components": [
    {
      "name": "database",
      "healthy": false,
      "reason": "lorem ipsum dolor sit amet"
    }
  ]
}
```

### Readiness endpoint

The readiness endpoint indicates whether the application is ready to accept traffic. This is useful during startup
or deployment scenarios.

**Endpoint:** `GET /api/health/readiness`

**Response code:**
- `200 OK`: Application is ready to accept requests

## Usage scenarios

### Docker Compose health checks

You can configure Docker Compose to use the liveness endpoint to monitor container health:

```yaml
services:
  nginx-ignition:
    image: nginx-ignition:latest
    ports:
      - "8090:8090"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8090/api/health/liveness"]
      interval: 30s
      timeout: 3s
      start_period: 5s
      retries: 3
```

### Kubernetes probes

For Kubernetes deployments, configure both liveness and readiness probes:

```yaml
livenessProbe:
  httpGet:
    path: /api/health/liveness
    port: 8090
  initialDelaySeconds: 10
  periodSeconds: 30

readinessProbe:
  httpGet:
    path: /api/health/readiness
    port: 8090
  initialDelaySeconds: 5
  periodSeconds: 10
```
