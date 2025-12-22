package performance

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rshade/pulumicost-mcp/gen/cost"
	"github.com/rshade/pulumicost-mcp/internal/adapter"
	"github.com/rshade/pulumicost-mcp/internal/logging"
	"github.com/rshade/pulumicost-mcp/internal/service"
)

// TestConcurrentRequests verifies 50 concurrent requests without degradation (SC-002)
func TestConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	const (
		concurrentUsers = 50
		requestsPerUser = 10
		totalRequests   = concurrentUsers * requestsPerUser
	)

	logger := logging.Default()
	adapter := adapter.NewPulumiCostAdapter("/usr/local/bin/pulumicost-core")
	svc := service.NewCostService(adapter, logger)

	pulumiJSON := generate100ResourceStack()

	var (
		successCount uint64
		errorCount   uint64
		wg           sync.WaitGroup
		latencies    = make([]time.Duration, 0, totalRequests)
		latencyMutex sync.Mutex
	)

	startTime := time.Now()

	t.Logf("Starting load test: %d concurrent users, %d requests each = %d total",
		concurrentUsers, requestsPerUser, totalRequests)

	// Spawn concurrent users
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			for j := 0; j < requestsPerUser; j++ {
				reqStart := time.Now()

				_, err := svc.AnalyzeProjected(context.Background(), &cost.AnalyzeProjectedPayload{
					PulumiJSON: pulumiJSON,
				})

				latency := time.Since(reqStart)

				latencyMutex.Lock()
				latencies = append(latencies, latency)
				latencyMutex.Unlock()

				if err != nil {
					atomic.AddUint64(&errorCount, 1)
					t.Logf("User %d request %d failed: %v", userID, j, err)
				} else {
					atomic.AddUint64(&successCount, 1)
				}
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	totalDuration := time.Since(startTime)

	// Calculate statistics
	stats := calculateLatencyStats(latencies)
	successRate := float64(successCount) / float64(totalRequests) * 100
	throughput := float64(totalRequests) / totalDuration.Seconds()

	t.Logf("\n=== Load Test Results ===")
	t.Logf("Total Requests:  %d", totalRequests)
	t.Logf("Successful:      %d (%.1f%%)", successCount, successRate)
	t.Logf("Failed:          %d", errorCount)
	t.Logf("Total Duration:  %v", totalDuration)
	t.Logf("Throughput:      %.1f req/s", throughput)
	t.Logf("\nLatency Statistics:")
	t.Logf("  Min:  %v", stats.Min)
	t.Logf("  P50:  %v", stats.P50)
	t.Logf("  P90:  %v", stats.P90)
	t.Logf("  P95:  %v", stats.P95)
	t.Logf("  P99:  %v", stats.P99)
	t.Logf("  Max:  %v", stats.Max)

	// SC-002: 50 concurrent requests without degradation
	// Define "degradation" as:
	// 1. Success rate < 99%
	// 2. P95 latency > 3s (compared to SC-001's 2.9s baseline)

	if successRate < 99.0 {
		t.Errorf("Success rate %.1f%% is below 99%% threshold (SC-002)", successRate)
	} else {
		t.Logf("✓ Success rate %.1f%% meets requirement", successRate)
	}

	maxP95 := 3 * time.Second
	if stats.P95 > maxP95 {
		t.Errorf("P95 latency %v exceeds %v threshold under load (SC-002)", stats.P95, maxP95)
	} else {
		t.Logf("✓ P95 latency %v meets requirement under concurrent load", stats.P95)
	}
}

// TestSustainedLoad tests sustained load over a longer period
func TestSustainedLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping sustained load test in short mode")
	}

	const (
		duration      = 60 * time.Second
		concurrentReq = 10
	)

	logger := logging.Default()
	adapter := adapter.NewPulumiCostAdapter("/usr/local/bin/pulumicost-core")
	svc := service.NewCostService(adapter, logger)

	pulumiJSON := generate100ResourceStack()

	var (
		requestCount uint64
		errorCount   uint64
		wg           sync.WaitGroup
		stopChan     = make(chan struct{})
	)

	t.Logf("Running sustained load test for %v with %d concurrent requests",
		duration, concurrentReq)

	startTime := time.Now()

	// Spawn workers
	for i := 0; i < concurrentReq; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-stopChan:
					return
				default:
					_, err := svc.AnalyzeProjected(context.Background(), &cost.AnalyzeProjectedPayload{
						PulumiJSON: pulumiJSON,
					})

					atomic.AddUint64(&requestCount, 1)
					if err != nil {
						atomic.AddUint64(&errorCount, 1)
					}
				}
			}
		}()
	}

	// Run for specified duration
	time.Sleep(duration)
	close(stopChan)
	wg.Wait()

	totalDuration := time.Since(startTime)
	reqCount := atomic.LoadUint64(&requestCount)
	errCount := atomic.LoadUint64(&errorCount)
	successRate := float64(reqCount-errCount) / float64(reqCount) * 100
	throughput := float64(reqCount) / totalDuration.Seconds()

	t.Logf("\n=== Sustained Load Results ===")
	t.Logf("Duration:        %v", totalDuration)
	t.Logf("Total Requests:  %d", reqCount)
	t.Logf("Successful:      %d (%.1f%%)", reqCount-errCount, successRate)
	t.Logf("Failed:          %d", errCount)
	t.Logf("Throughput:      %.1f req/s", throughput)

	if successRate < 95.0 {
		t.Errorf("Success rate %.1f%% dropped below 95%% during sustained load", successRate)
	} else {
		t.Logf("✓ Success rate %.1f%% maintained under sustained load", successRate)
	}
}

// TestMemoryLeaks runs a long test to detect memory leaks
func TestMemoryLeaks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory leak test in short mode")
	}

	logger := logging.Default()
	adapter := adapter.NewPulumiCostAdapter("/usr/local/bin/pulumicost-core")
	svc := service.NewCostService(adapter, logger)

	pulumiJSON := generate100ResourceStack()

	const iterations = 1000

	t.Logf("Running %d iterations to detect memory leaks...", iterations)

	for i := 0; i < iterations; i++ {
		_, err := svc.AnalyzeProjected(context.Background(), &cost.AnalyzeProjectedPayload{
			PulumiJSON: pulumiJSON,
		})

		if err != nil {
			t.Fatalf("Iteration %d failed: %v", i, err)
		}

		if (i+1)%100 == 0 {
			t.Logf("Completed %d/%d iterations", i+1, iterations)
		}
	}

	t.Logf("✓ Completed %d iterations without errors", iterations)
	t.Log("Note: Run with -memprofile to analyze memory usage")
}

// BenchmarkConcurrency benchmarks concurrent request handling
func BenchmarkConcurrency(b *testing.B) {
	logger := logging.Default()
	adapter := adapter.NewPulumiCostAdapter("/usr/local/bin/pulumicost-core")
	svc := service.NewCostService(adapter, logger)

	pulumiJSON := generate100ResourceStack()

	benchmarks := []struct {
		name        string
		concurrency int
	}{
		{"1_concurrent", 1},
		{"10_concurrent", 10},
		{"50_concurrent", 50},
		{"100_concurrent", 100},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.SetParallelism(bm.concurrency)
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := svc.AnalyzeProjected(context.Background(), &cost.AnalyzeProjectedPayload{
						PulumiJSON: pulumiJSON,
					})

					if err != nil {
						b.Fatalf("Benchmark failed: %v", err)
					}
				}
			})
		})
	}
}
