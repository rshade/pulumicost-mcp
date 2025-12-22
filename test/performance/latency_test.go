package performance

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/rshade/pulumicost-mcp/gen/cost"
	"github.com/rshade/pulumicost-mcp/internal/adapter"
	"github.com/rshade/pulumicost-mcp/internal/logging"
	"github.com/rshade/pulumicost-mcp/internal/service"
)

// TestP95Latency verifies that P95 latency is <3s for 100-resource stacks (SC-001)
func TestP95Latency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Create test infrastructure with 100 resources
	pulumiJSON := generate100ResourceStack()

	// Setup service
	logger := logging.Default()
	adapter := adapter.NewPulumiCostAdapter("/usr/local/bin/pulumicost-core")
	svc := service.NewCostService(adapter, logger)

	// Run 100 iterations to get meaningful P95
	const iterations = 100
	latencies := make([]time.Duration, iterations)

	t.Logf("Running %d iterations with 100-resource stack...", iterations)

	for i := 0; i < iterations; i++ {
		start := time.Now()

		_, err := svc.AnalyzeProjected(context.Background(), &cost.AnalyzeProjectedPayload{
			PulumiJSON: pulumiJSON,
		})

		latencies[i] = time.Since(start)

		if err != nil {
			t.Fatalf("Iteration %d failed: %v", i, err)
		}

		// Log progress every 10 iterations
		if (i+1)%10 == 0 {
			t.Logf("Completed %d/%d iterations", i+1, iterations)
		}
	}

	// Calculate statistics
	stats := calculateLatencyStats(latencies)

	t.Logf("\n=== Latency Statistics (100 resources, %d iterations) ===", iterations)
	t.Logf("Min:    %v", stats.Min)
	t.Logf("P50:    %v", stats.P50)
	t.Logf("P90:    %v", stats.P90)
	t.Logf("P95:    %v", stats.P95)
	t.Logf("P99:    %v", stats.P99)
	t.Logf("Max:    %v", stats.Max)
	t.Logf("Mean:   %v", stats.Mean)

	// SC-001: P95 latency must be ≤2.9 seconds for 100-resource stacks
	maxP95 := 2900 * time.Millisecond

	if stats.P95 > maxP95 {
		t.Errorf("P95 latency %v exceeds requirement of %v (SC-001)", stats.P95, maxP95)
	} else {
		t.Logf("✓ P95 latency %v meets requirement of ≤%v", stats.P95, maxP95)
	}
}

// TestLatencyByResourceCount measures how latency scales with resource count
func TestLatencyByResourceCount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	resourceCounts := []int{10, 25, 50, 100}
	logger := logging.Default()
	adapter := adapter.NewPulumiCostAdapter("/usr/local/bin/pulumicost-core")
	svc := service.NewCostService(adapter, logger)

	t.Log("\n=== Latency by Resource Count ===")

	for _, count := range resourceCounts {
		pulumiJSON := generateResourceStack(count)

		// Run 20 iterations per count
		const iterations = 20
		latencies := make([]time.Duration, iterations)

		for i := 0; i < iterations; i++ {
			start := time.Now()

			_, err := svc.AnalyzeProjected(context.Background(), &cost.AnalyzeProjectedPayload{
				PulumiJSON: pulumiJSON,
			})

			latencies[i] = time.Since(start)

			if err != nil {
				t.Fatalf("Failed with %d resources: %v", count, err)
			}
		}

		stats := calculateLatencyStats(latencies)

		t.Logf("%3d resources - P50: %6v  P95: %6v  P99: %6v",
			count, stats.P50, stats.P95, stats.P99)
	}
}

// BenchmarkAnalyzeProjected benchmarks the AnalyzeProjected method
func BenchmarkAnalyzeProjected(b *testing.B) {
	logger := logging.Default()
	adapter := adapter.NewPulumiCostAdapter("/usr/local/bin/pulumicost-core")
	svc := service.NewCostService(adapter, logger)

	benchmarks := []struct {
		name          string
		resourceCount int
	}{
		{"10_resources", 10},
		{"50_resources", 50},
		{"100_resources", 100},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			pulumiJSON := generateResourceStack(bm.resourceCount)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := svc.AnalyzeProjected(context.Background(), &cost.AnalyzeProjectedPayload{
					PulumiJSON: pulumiJSON,
				})

				if err != nil {
					b.Fatalf("Benchmark failed: %v", err)
				}
			}
		})
	}
}

// Helper functions

type LatencyStats struct {
	Min  time.Duration
	P50  time.Duration
	P90  time.Duration
	P95  time.Duration
	P99  time.Duration
	Max  time.Duration
	Mean time.Duration
}

func calculateLatencyStats(latencies []time.Duration) LatencyStats {
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	percentile := func(p float64) time.Duration {
		idx := int(float64(len(sorted)) * p)
		if idx >= len(sorted) {
			idx = len(sorted) - 1
		}
		return sorted[idx]
	}

	var sum time.Duration
	for _, lat := range sorted {
		sum += lat
	}
	mean := sum / time.Duration(len(sorted))

	return LatencyStats{
		Min:  sorted[0],
		P50:  percentile(0.50),
		P90:  percentile(0.90),
		P95:  percentile(0.95),
		P99:  percentile(0.99),
		Max:  sorted[len(sorted)-1],
		Mean: mean,
	}
}

func generate100ResourceStack() string {
	return generateResourceStack(100)
}

func generateResourceStack(count int) string {
	resources := make([]map[string]interface{}, count)

	// Generate diverse mix of resources
	for i := 0; i < count; i++ {
		resourceType := ""
		properties := make(map[string]interface{})

		switch i % 5 {
		case 0:
			// EC2 instances
			resourceType = "aws:ec2/instance:Instance"
			properties["instanceType"] = "t3.medium"
			properties["ami"] = "ami-12345678"
		case 1:
			// RDS instances
			resourceType = "aws:rds/instance:Instance"
			properties["instanceClass"] = "db.t3.small"
			properties["engine"] = "postgres"
		case 2:
			// S3 buckets
			resourceType = "aws:s3/bucket:Bucket"
			properties["acl"] = "private"
		case 3:
			// Lambda functions
			resourceType = "aws:lambda/function:Function"
			properties["runtime"] = "python3.9"
			properties["memorySize"] = 256
		case 4:
			// VPCs
			resourceType = "aws:ec2/vpc:Vpc"
			properties["cidrBlock"] = "10.0.0.0/16"
		}

		resources[i] = map[string]interface{}{
			"urn":  fmt.Sprintf("urn:pulumi:dev::myapp::%s::resource-%d", resourceType, i),
			"type": resourceType,
			"inputs": properties,
		}
	}

	stack := map[string]interface{}{
		"resources": resources,
	}

	bytes, _ := json.Marshal(stack)
	return string(bytes)
}
