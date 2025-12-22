# Horizontal Scaling Validation (SC-010)

## Requirement

System must scale horizontally to serve 500+ concurrent users with P95 latency <3.5 seconds.

## Validation Strategy

### Phase 1: Load Test Infrastructure Setup

```yaml
# kubernetes/load-test-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pulumicost-mcp
spec:
  replicas: 3  # Start with 3 replicas
  selector:
    matchLabels:
      app: pulumicost-mcp
  template:
    metadata:
      labels:
        app: pulumicost-mcp
    spec:
      containers:
      - name: mcp-server
        image: pulumicost-mcp:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        env:
        - name: MCP_LOG_LEVEL
          value: "info"
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: pulumicost-mcp-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: pulumicost-mcp
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Phase 2: Load Generation

```go
// test/scalability/horizontal_scaling_test.go
package scalability

import (
    "context"
    "sync"
    "sync/atomic"
    "testing"
    "time"
)

func TestHorizontalScaling500Users(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping scaling test in short mode")
    }

    const (
        targetUsers       = 500
        rampUpDuration    = 5 * time.Minute
        sustainedDuration = 10 * time.Minute
        requestsPerUser   = 20
    )

    // Calculate users to add per second during ramp-up
    rampRate := float64(targetUsers) / rampUpDuration.Seconds()

    var (
        activeUsers   uint64
        totalRequests uint64
        errorCount    uint64
        latencies     []time.Duration
        latencyMutex  sync.Mutex
        wg            sync.WaitGroup
        stopChan      = make(chan struct{})
    )

    t.Logf("Starting horizontal scaling test:")
    t.Logf("  Ramp up: %d users over %v", targetUsers, rampUpDuration)
    t.Logf("  Sustain: %v at %d users", sustainedDuration, targetUsers)

    startTime := time.Now()

    // Ramp-up phase
    go func() {
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                currentUsers := atomic.LoadUint64(&activeUsers)
                if currentUsers >= targetUsers {
                    return
                }

                // Add users based on ramp rate
                usersToAdd := int(rampRate)
                for i := 0; i < usersToAdd; i++ {
                    wg.Add(1)
                    go simulateUser(&wg, &totalRequests, &errorCount, &latencies, &latencyMutex, stopChan, requestsPerUser)
                    atomic.AddUint64(&activeUsers, 1)
                }

            case <-stopChan:
                return
            }
        }
    }()

    // Wait for ramp-up to complete
    time.Sleep(rampUpDuration)

    t.Logf("Ramp-up complete: %d active users", atomic.LoadUint64(&activeUsers))
    t.Logf("Sustaining load for %v...", sustainedDuration)

    // Sustained load phase
    time.Sleep(sustainedDuration)

    // Shutdown
    close(stopChan)
    wg.Wait()

    totalDuration := time.Since(startTime)

    // Calculate statistics
    reqCount := atomic.LoadUint64(&totalRequests)
    errCount := atomic.LoadUint64(&errorCount)
    successRate := float64(reqCount-errCount) / float64(reqCount) * 100
    throughput := float64(reqCount) / totalDuration.Seconds()

    latencyMutex.Lock()
    stats := calculateLatencyStats(latencies)
    latencyMutex.Unlock()

    t.Logf("\n=== Horizontal Scaling Results ===")
    t.Logf("Total Duration:     %v", totalDuration)
    t.Logf("Peak Users:         %d", atomic.LoadUint64(&activeUsers))
    t.Logf("Total Requests:     %d", reqCount)
    t.Logf("Successful:         %d (%.1f%%)", reqCount-errCount, successRate)
    t.Logf("Failed:             %d", errCount)
    t.Logf("Throughput:         %.1f req/s", throughput)
    t.Logf("\nLatency Statistics:")
    t.Logf("  P50: %v", stats.P50)
    t.Logf("  P90: %v", stats.P90)
    t.Logf("  P95: %v", stats.P95)
    t.Logf("  P99: %v", stats.P99)

    // SC-010: Validate requirements
    if successRate < 99.0 {
        t.Errorf("Success rate %.1f%% below 99%% threshold (SC-010)", successRate)
    } else {
        t.Logf("✓ Success rate %.1f%% meets requirement", successRate)
    }

    maxP95 := 3500 * time.Millisecond
    if stats.P95 > maxP95 {
        t.Errorf("P95 latency %v exceeds %v at 500 users (SC-010)", stats.P95, maxP95)
    } else {
        t.Logf("✓ P95 latency %v meets requirement at 500 users", stats.P95)
    }
}

func simulateUser(wg *sync.WaitGroup, reqCount, errCount *uint64, latencies *[]time.Duration, mutex *sync.Mutex, stopChan chan struct{}, requests int) {
    defer wg.Done()

    for i := 0; i < requests; i++ {
        select {
        case <-stopChan:
            return
        default:
            start := time.Now()
            err := makeRequest()
            latency := time.Since(start)

            atomic.AddUint64(reqCount, 1)
            if err != nil {
                atomic.AddUint64(errCount, 1)
            }

            mutex.Lock()
            *latencies = append(*latencies, latency)
            mutex.Unlock()

            // Think time between requests
            time.Sleep(100 * time.Millisecond)
        }
    }
}
```

### Phase 3: Monitoring During Test

```bash
#!/bin/bash
# scripts/monitor-scaling-test.sh

echo "=== Monitoring Horizontal Scaling Test ==="

while true; do
    # Pod count
    echo -n "Pods: "
    kubectl get pods -l app=pulumicost-mcp --no-headers | wc -l

    # CPU/Memory usage
    echo "Resource Usage:"
    kubectl top pods -l app=pulumicost-mcp

    # Request metrics
    echo "Request Metrics:"
    kubectl exec -it $(kubectl get pod -l app=pulumicost-mcp -o jsonpath='{.items[0].metadata.name}') -- \
        curl -s localhost:8080/metrics | grep pulumicost_requests_total

    echo "---"
    sleep 10
done
```

### Phase 4: Results Analysis

```python
# scripts/analyze-scaling-results.py
import matplotlib.pyplot as plt
import pandas as pd

def analyze_scaling_test(metrics_file):
    """Analyze horizontal scaling test results."""

    # Load metrics
    df = pd.read_csv(metrics_file)

    # Create visualizations
    fig, axes = plt.subplots(2, 2, figsize=(15, 10))

    # 1. Latency over time
    axes[0, 0].plot(df['timestamp'], df['p95_latency_ms'])
    axes[0, 0].axhline(y=3500, color='r', linestyle='--', label='3.5s threshold')
    axes[0, 0].set_title('P95 Latency Over Time')
    axes[0, 0].set_xlabel('Time')
    axes[0, 0].set_ylabel('Latency (ms)')
    axes[0, 0].legend()

    # 2. Throughput vs pod count
    axes[0, 1].scatter(df['pod_count'], df['throughput_rps'])
    axes[0, 1].set_title('Throughput vs Pod Count')
    axes[0, 1].set_xlabel('Number of Pods')
    axes[0, 1].set_ylabel('Throughput (req/s)')

    # 3. Error rate over time
    axes[1, 0].plot(df['timestamp'], df['error_rate_percent'])
    axes[1, 0].axhline(y=1.0, color='r', linestyle='--', label='1% threshold')
    axes[1, 0].set_title('Error Rate Over Time')
    axes[1, 0].set_xlabel('Time')
    axes[1, 0].set_ylabel('Error Rate (%)')
    axes[1, 0].legend()

    # 4. Active users over time
    axes[1, 1].plot(df['timestamp'], df['active_users'])
    axes[1, 1].set_title('Active Users Over Time')
    axes[1, 1].set_xlabel('Time')
    axes[1, 1].set_ylabel('Active Users')

    plt.tight_layout()
    plt.savefig('scaling-test-results.png', dpi=300)
    print("Results saved to scaling-test-results.png")

    # Calculate key metrics
    print("\n=== Scaling Test Summary ===")
    print(f"Peak Users:          {df['active_users'].max()}")
    print(f"Peak Throughput:     {df['throughput_rps'].max():.1f} req/s")
    print(f"Max P95 Latency:     {df['p95_latency_ms'].max():.0f} ms")
    print(f"Peak Pod Count:      {df['pod_count'].max()}")
    print(f"Avg Error Rate:      {df['error_rate_percent'].mean():.2f}%")

    # Validate requirements
    max_latency = df['p95_latency_ms'].max()
    avg_error_rate = df['error_rate_percent'].mean()

    if max_latency <= 3500:
        print(f"✓ P95 latency {max_latency:.0f}ms within 3.5s requirement")
    else:
        print(f"✗ P95 latency {max_latency:.0f}ms exceeds 3.5s requirement")

    if avg_error_rate <= 1.0:
        print(f"✓ Error rate {avg_error_rate:.2f}% within 1% threshold")
    else:
        print(f"✗ Error rate {avg_error_rate:.2f}% exceeds 1% threshold")

if __name__ == "__main__":
    analyze_scaling_test("scaling_metrics.csv")
```

## Acceptance Criteria

- [ ] System handles 500+ concurrent users
- [ ] P95 latency remains <3.5s under load
- [ ] Error rate stays <1%
- [ ] HPA successfully scales pods (3-10 replicas)
- [ ] Resource utilization remains healthy (<80% CPU/Memory)
- [ ] No memory leaks detected over sustained period
- [ ] Graceful degradation under extreme load

## Known Considerations

1. **Network Latency**
   - Geographic distribution affects latency
   - CDN/load balancer adds overhead
   - Inter-pod communication varies

2. **Database Connections**
   - Connection pool sizing critical
   - Plugin gRPC connections need pooling
   - May become bottleneck before pods

3. **Shared State**
   - Stateless design enables scaling
   - Metrics aggregation may need sharding
   - Tracing backend capacity

## Continuous Monitoring

```yaml
# prometheus/alerts.yml
groups:
  - name: scaling
    interval: 30s
    rules:
      - alert: HighLatencyUnderLoad
        expr: histogram_quantile(0.95, pulumicost_request_duration_seconds_bucket) > 3.5
        for: 5m
        annotations:
          summary: "P95 latency exceeds 3.5s under load"

      - alert: HPANotScaling
        expr: kube_deployment_spec_replicas{deployment="pulumicost-mcp"} == kube_deployment_status_replicas_available{deployment="pulumicost-mcp"}
        for: 10m
        annotations:
          summary: "HPA not scaling despite high load"
```

## References

- Kubernetes HPA: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/
- Load Testing Best Practices: https://grafana.com/docs/k6/latest/testing-guides/
