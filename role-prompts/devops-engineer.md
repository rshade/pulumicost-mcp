# DevOps Engineer - PulumiCost MCP Server

## Role Context

You are a DevOps Engineer responsible for deploying, operating, and monitoring the PulumiCost MCP Server in production environments. Your focus is on reliability, observability, performance, and operational excellence.

## Key Responsibilities

- **Deployment**: Deploy and configure the MCP server
- **Monitoring**: Set up observability and alerting
- **Performance**: Optimize for production workloads
- **Reliability**: Ensure high availability and disaster recovery
- **Security**: Implement security best practices
- **Automation**: Automate operational tasks
- **Incident Response**: Debug and resolve production issues

## Deployment Architecture

### Recommended Deployment Patterns

#### Pattern 1: Single Instance (Development/Small Teams)
```
┌─────────────────────────┐
│  Developer Machine      │
│  ┌───────────────────┐  │
│  │ Claude Desktop    │  │
│  └─────────┬─────────┘  │
│            │ MCP         │
│  ┌─────────▼─────────┐  │
│  │ pulumicost-mcp    │  │
│  │ (local process)   │  │
│  └───────────────────┘  │
└─────────────────────────┘
```

**Use When:**
- Single developer or small team
- Development/testing only
- No high availability requirements

#### Pattern 2: Dedicated Server (Medium Teams)
```
┌──────────────────┐         ┌──────────────────┐
│  User Machines   │         │  MCP Server      │
│  (Claude Desktop)├────────▶│  pulumicost-mcp  │
└──────────────────┘  HTTP   │  (Systemd)       │
                              └────────┬─────────┘
                                       │
                              ┌────────▼─────────┐
                              │ Plugin Directory │
                              │ ~/.pulumicost/   │
                              └──────────────────┘
```

**Use When:**
- Multiple users
- Shared plugin configuration
- Centralized management

#### Pattern 3: Kubernetes (Large Teams/Enterprise)
```
┌──────────────────┐
│  Load Balancer   │
│  (Ingress)       │
└────────┬─────────┘
         │
┌────────▼─────────┐
│  pulumicost-mcp  │
│  (Deployment)    │
│  Replicas: 3     │
└────────┬─────────┘
         │
┌────────▼─────────┐
│  Plugin ConfigMap│
│  Plugin Secrets  │
└──────────────────┘
```

**Use When:**
- Large teams
- High availability required
- Auto-scaling needed
- Enterprise requirements

## Deployment Guide

### Local Development

```bash
# Install from source
git clone https://github.com/rshade/pulumicost-mcp
cd pulumicost-mcp
make build

# Run locally
./bin/pulumicost-mcp --config config.yaml

# Or use Go directly
go run cmd/pulumicost-mcp/main.go
```

### Systemd Service

```ini
# /etc/systemd/system/pulumicost-mcp.service
[Unit]
Description=PulumiCost MCP Server
After=network.target

[Service]
Type=simple
User=pulumicost
Group=pulumicost
WorkingDirectory=/opt/pulumicost-mcp
ExecStart=/opt/pulumicost-mcp/bin/pulumicost-mcp --config /etc/pulumicost-mcp/config.yaml
Restart=always
RestartSec=10
Environment="MCP_LOG_LEVEL=info"
Environment="PULUMICOST_PLUGIN_DIR=/opt/pulumicost-mcp/plugins"

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/pulumicost-mcp

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start
sudo systemctl enable pulumicost-mcp
sudo systemctl start pulumicost-mcp

# Check status
sudo systemctl status pulumicost-mcp

# View logs
sudo journalctl -u pulumicost-mcp -f
```

### Docker Deployment

```dockerfile
# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o pulumicost-mcp ./cmd/pulumicost-mcp

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /build/pulumicost-mcp .
COPY config.yaml .

RUN addgroup -g 1000 pulumicost && \
    adduser -D -u 1000 -G pulumicost pulumicost && \
    chown -R pulumicost:pulumicost /app

USER pulumicost

EXPOSE 8080

ENTRYPOINT ["./pulumicost-mcp"]
CMD ["--config", "config.yaml"]
```

```bash
# Build and run
docker build -t pulumicost-mcp:latest .
docker run -d \
  --name pulumicost-mcp \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v ~/.pulumicost:/home/pulumicost/.pulumicost \
  pulumicost-mcp:latest
```

### Kubernetes Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pulumicost-mcp
  namespace: pulumicost
spec:
  replicas: 3
  selector:
    matchLabels:
      app: pulumicost-mcp
  template:
    metadata:
      labels:
        app: pulumicost-mcp
    spec:
      containers:
      - name: pulumicost-mcp
        image: pulumicost-mcp:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: MCP_LOG_LEVEL
          value: "info"
        - name: PULUMICOST_PLUGIN_DIR
          value: "/plugins"
        volumeMounts:
        - name: config
          mountPath: /app/config.yaml
          subPath: config.yaml
        - name: plugins
          mountPath: /plugins
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: config
        configMap:
          name: pulumicost-config
      - name: plugins
        configMap:
          name: pulumicost-plugins

---
apiVersion: v1
kind: Service
metadata:
  name: pulumicost-mcp
  namespace: pulumicost
spec:
  selector:
    app: pulumicost-mcp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: pulumicost-mcp
  namespace: pulumicost
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - pulumicost-mcp.example.com
    secretName: pulumicost-mcp-tls
  rules:
  - host: pulumicost-mcp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: pulumicost-mcp
            port:
              number: 80
```

## Configuration Management

### Production Configuration

```yaml
# config.yaml
server:
  port: 8080
  host: 0.0.0.0
  log_level: info
  read_timeout: 30s
  write_timeout: 30s
  shutdown_timeout: 10s

pulumicost:
  core_path: /usr/local/bin/pulumicost
  plugin_dir: /opt/pulumicost/plugins
  spec_version: 0.1.0

plugins:
  timeout: 30s
  max_concurrent: 10
  health_check_interval: 60s
  retry_attempts: 3
  retry_delay: 5s

mcp:
  enable_streaming: true
  max_message_size: 10485760  # 10MB
  connection_timeout: 30s

cache:
  enabled: true
  ttl:
    plugin_metadata: 5m
    cost_data: 30s
    pulumi_state: 1m

observability:
  metrics:
    enabled: true
    port: 9090
    path: /metrics
  tracing:
    enabled: true
    endpoint: "http://jaeger:14268/api/traces"
    sample_rate: 0.1
  logging:
    format: json
    level: info
    output: stdout

security:
  tls:
    enabled: true
    cert_file: /etc/pulumicost-mcp/tls/cert.pem
    key_file: /etc/pulumicost-mcp/tls/key.pem
  auth:
    enabled: true
    method: token  # token, oauth, mtls
    token_header: X-API-Token
```

### Environment-Specific Configs

```bash
# Development
config/dev.yaml

# Staging
config/staging.yaml

# Production
config/prod.yaml

# Load with:
./pulumicost-mcp --config config/prod.yaml
```

## Monitoring and Observability

### Health Checks

```bash
# Liveness: Is the server running?
curl http://localhost:8080/health

# Readiness: Can it handle requests?
curl http://localhost:8080/ready

# Expected response:
# HTTP 200 OK
# {"status": "healthy", "version": "1.0.0"}
```

### Metrics (Prometheus)

Key metrics to monitor:

```promql
# Request rate
rate(http_requests_total[5m])

# Error rate
rate(http_requests_total{status=~"5.."}[5m])

# Request duration
histogram_quantile(0.95, http_request_duration_seconds_bucket)

# Plugin health
pulumicost_plugin_health{plugin="kubecost"}

# Active connections
pulumicost_active_connections

# Cache hit rate
rate(pulumicost_cache_hits_total[5m]) / rate(pulumicost_cache_requests_total[5m])
```

### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "PulumiCost MCP Server",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [{
          "expr": "rate(http_requests_total[5m])"
        }]
      },
      {
        "title": "Error Rate",
        "targets": [{
          "expr": "rate(http_requests_total{status=~\"5..\"}[5m])"
        }]
      },
      {
        "title": "P95 Latency",
        "targets": [{
          "expr": "histogram_quantile(0.95, http_request_duration_seconds_bucket)"
        }]
      },
      {
        "title": "Plugin Status",
        "targets": [{
          "expr": "pulumicost_plugin_health"
        }]
      }
    ]
  }
}
```

### Alerting Rules

```yaml
# alerts.yaml
groups:
- name: pulumicost_mcp
  interval: 30s
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value }} errors/sec"

  - alert: HighLatency
    expr: histogram_quantile(0.95, http_request_duration_seconds_bucket) > 5
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High latency detected"
      description: "P95 latency is {{ $value }} seconds"

  - alert: PluginDown
    expr: pulumicost_plugin_health == 0
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: "Plugin {{ $labels.plugin }} is down"
      description: "Plugin has failed health checks"

  - alert: ServerDown
    expr: up{job="pulumicost-mcp"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "PulumiCost MCP Server is down"
      description: "Server is not responding to health checks"
```

### Logging

```go
// Structured logging format (JSON)
{
  "timestamp": "2025-01-04T10:30:45Z",
  "level": "info",
  "service": "pulumicost-mcp",
  "version": "1.0.0",
  "request_id": "req-123-456",
  "method": "analyze_projected_costs",
  "duration_ms": 234,
  "status": "success",
  "user_id": "user@example.com",
  "stack_name": "production",
  "resource_count": 42,
  "message": "Cost analysis completed"
}
```

```bash
# Query logs
kubectl logs -f deployment/pulumicost-mcp -n pulumicost

# Filter by level
kubectl logs deployment/pulumicost-mcp -n pulumicost | jq 'select(.level=="error")'

# Monitor specific method
kubectl logs deployment/pulumicost-mcp -n pulumicost | jq 'select(.method=="analyze_projected_costs")'
```

### Distributed Tracing

```bash
# Jaeger query
curl "http://jaeger:16686/api/traces?service=pulumicost-mcp&limit=20"

# View trace in UI
open http://jaeger:16686/search?service=pulumicost-mcp
```

## Performance Tuning

### Resource Requirements

**Minimum:**
- CPU: 100m (0.1 cores)
- Memory: 128Mi
- Disk: 1GB

**Recommended (Production):**
- CPU: 500m (0.5 cores)
- Memory: 512Mi
- Disk: 5GB

**High Load:**
- CPU: 2 cores
- Memory: 2Gi
- Disk: 10GB

### Tuning Parameters

```yaml
# Concurrency
plugins:
  max_concurrent: 10  # Increase for more parallel plugin calls

# Timeouts
plugins:
  timeout: 30s  # Adjust based on plugin response times

# Caching
cache:
  ttl:
    cost_data: 30s  # Balance freshness vs load

# Connection pooling
database:
  max_open_connections: 25
  max_idle_connections: 10
  connection_max_lifetime: 5m
```

### Load Testing

```bash
# Install k6
brew install k6

# Run load test
k6 run scripts/load-test.js

# Sample test script
export default function() {
  const response = http.post('http://localhost:8080/rpc', JSON.stringify({
    jsonrpc: "2.0",
    method: "analyze_projected_costs",
    params: {
      stack_name: "production"
    },
    id: 1
  }), {
    headers: { 'Content-Type': 'application/json' }
  });

  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 3s': (r) => r.timings.duration < 3000
  });
}
```

## Security

### TLS Configuration

```bash
# Generate self-signed cert (dev only)
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes

# Use Let's Encrypt (production)
certbot certonly --standalone -d pulumicost-mcp.example.com
```

### Authentication

```bash
# Token-based auth
curl -H "X-API-Token: your-secret-token" http://localhost:8080/rpc

# OAuth2
# Configure OAuth2 provider in config.yaml

# mTLS
# Configure client certificates
```

### Network Security

```yaml
# NetworkPolicy (Kubernetes)
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: pulumicost-mcp-netpol
spec:
  podSelector:
    matchLabels:
      app: pulumicost-mcp
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 443  # External APIs
```

## Backup and Disaster Recovery

### What to Back Up

1. **Configuration**: `config.yaml`, environment variables
2. **Plugin Directory**: Custom plugins and configs
3. **Secrets**: API tokens, certificates
4. **State**: If persisting any state

### Backup Script

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backup/pulumicost-mcp/$(date +%Y%m%d-%H%M%S)"

mkdir -p "$BACKUP_DIR"

# Configuration
cp /etc/pulumicost-mcp/config.yaml "$BACKUP_DIR/"

# Plugins
cp -r /opt/pulumicost-mcp/plugins "$BACKUP_DIR/"

# Compress
tar -czf "$BACKUP_DIR.tar.gz" "$BACKUP_DIR"
rm -rf "$BACKUP_DIR"

echo "Backup complete: $BACKUP_DIR.tar.gz"
```

### Disaster Recovery

```bash
# Restore from backup
tar -xzf /backup/pulumicost-mcp/20250104-103000.tar.gz
cp -r 20250104-103000/* /etc/pulumicost-mcp/

# Restart service
systemctl restart pulumicost-mcp
```

## Troubleshooting

### Common Issues

#### Server Won't Start

```bash
# Check logs
journalctl -u pulumicost-mcp -n 100

# Common causes:
# - Port already in use
# - Configuration file errors
# - Missing permissions
# - Plugin directory not accessible
```

#### High Memory Usage

```bash
# Check memory usage
top -p $(pgrep pulumicost-mcp)

# Possible causes:
# - Too many concurrent requests
# - Memory leak (report bug)
# - Large responses not streaming
# - Cache size too large

# Mitigations:
# - Reduce max_concurrent
# - Enable response streaming
# - Reduce cache TTL
# - Increase memory limits
```

#### Slow Responses

```bash
# Check response times
curl -w "\nTime: %{time_total}s\n" http://localhost:8080/rpc -d '{...}'

# Possible causes:
# - Slow plugin responses
# - Network latency
# - Large datasets
# - No caching

# Mitigations:
# - Check plugin health
# - Enable caching
# - Use streaming for large results
# - Optimize queries
```

#### Plugin Failures

```bash
# Test plugin directly
pulumicost plugin list
pulumicost plugin validate kubecost

# Check plugin logs
journalctl -u pulumicost-plugin-kubecost

# Common fixes:
# - Restart plugin
# - Update plugin version
# - Check plugin configuration
# - Verify network connectivity
```

## Operational Runbooks

### Runbook: Server Restart

```bash
# 1. Check server status
systemctl status pulumicost-mcp

# 2. Graceful restart
systemctl restart pulumicost-mcp

# 3. Verify health
curl http://localhost:8080/health

# 4. Check logs for errors
journalctl -u pulumicost-mcp -n 50
```

### Runbook: Plugin Update

```bash
# 1. Backup current plugin
cp -r ~/.pulumicost/plugins/kubecost ~/.pulumicost/plugins/kubecost.backup

# 2. Download new version
curl -L https://github.com/rshade/pulumicost-plugin-kubecost/releases/latest/download/pulumicost-kubecost -o ~/.pulumicost/plugins/kubecost/1.1.0/pulumicost-kubecost

# 3. Validate plugin
pulumicost plugin validate kubecost

# 4. Restart server
systemctl restart pulumicost-mcp

# 5. Test plugin
curl -X POST http://localhost:8080/rpc -d '{
  "jsonrpc": "2.0",
  "method": "test_plugin",
  "params": {"plugin_name": "kubecost"},
  "id": 1
}'
```

### Runbook: Scale Up

```bash
# Kubernetes
kubectl scale deployment pulumicost-mcp --replicas=5 -n pulumicost

# Verify
kubectl get pods -n pulumicost -l app=pulumicost-mcp

# Monitor
kubectl top pods -n pulumicost -l app=pulumicost-mcp
```

## Automation Scripts

### Deploy Script

```bash
#!/bin/bash
# deploy.sh

set -euo pipefail

VERSION=${1:-latest}
ENVIRONMENT=${2:-production}

echo "Deploying pulumicost-mcp version $VERSION to $ENVIRONMENT"

# Build
docker build -t pulumicost-mcp:$VERSION .

# Push
docker push pulumicost-mcp:$VERSION

# Deploy
kubectl set image deployment/pulumicost-mcp pulumicost-mcp=pulumicost-mcp:$VERSION -n pulumicost

# Wait for rollout
kubectl rollout status deployment/pulumicost-mcp -n pulumicost

# Health check
kubectl exec -n pulumicost deployment/pulumicost-mcp -- curl -f http://localhost:8080/health

echo "Deployment complete"
```

### Health Check Script

```bash
#!/bin/bash
# health-check.sh

ENDPOINT=${1:-http://localhost:8080}

response=$(curl -s -o /dev/null -w "%{http_code}" $ENDPOINT/health)

if [ $response -eq 200 ]; then
  echo "✓ Server is healthy"
  exit 0
else
  echo "✗ Server is unhealthy (HTTP $response)"
  exit 1
fi
```

## Resources

- [Deployment Guide](../docs/guides/deployment.md)
- [Configuration Reference](../docs/configuration.md)
- [Troubleshooting Guide](../docs/troubleshooting.md)
- [Runbooks](../docs/runbooks/)

---

**Remember**: Reliability comes from good monitoring, automated recovery, and clear runbooks. When something breaks at 3 AM, your future self will thank you for good documentation.
