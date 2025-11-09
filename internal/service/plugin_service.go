package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rshade/pulumicost-mcp/gen/plugin"
)

// PluginService implements the plugin.Service interface
type PluginService struct {
	// In a real implementation, this would have dependencies like:
	// - pluginDirectory string
	// - grpcClientFactory func(address string) (*grpc.ClientConn, error)
	// For now, we use mock data since pulumicost-core isn't ready
}

// NewPluginService creates a new Plugin Service instance
func NewPluginService(pluginDir interface{}, logger interface{}) *PluginService {
	return &PluginService{}
}

// List returns all available cost source plugins
func (s *PluginService) List(ctx context.Context, payload *plugin.ListPayload) (*plugin.ListResult, error) {
	// Mock plugin data
	plugins := []*plugin.Plugin{
		{
			Name:        "aws-cost-source",
			Version:     "v1.0.0",
			Description: stringPtr("AWS Cost and Usage Report data source"),
			Capabilities: &plugin.PluginCapabilities{
				SupportsProjected: false,
				SupportsActual:    true,
				SupportsProviders: []string{"aws"},
			},
		},
		{
			Name:        "azure-cost-source",
			Version:     "v1.0.0",
			Description: stringPtr("Azure Cost Management data source"),
			Capabilities: &plugin.PluginCapabilities{
				SupportsProjected: false,
				SupportsActual:    true,
				SupportsProviders: []string{"azure"},
			},
		},
		{
			Name:        "infracost-plugin",
			Version:     "v0.10.0",
			Description: stringPtr("Infracost-based projected cost estimation"),
			Capabilities: &plugin.PluginCapabilities{
				SupportsProjected: true,
				SupportsActual:    false,
				SupportsProviders: []string{"aws", "azure", "gcp"},
			},
		},
	}

	// If health check is requested, populate health status
	if payload.IncludeHealth {
		for i := range plugins {
			latency := int64(15 + i*5) // Mock latency
			plugins[i].HealthStatus = &plugin.HealthStatus{
				Status:    "healthy",
				LastCheck: stringPtr(time.Now().Format(time.RFC3339)),
				LatencyMs: &latency,
			}
		}
	}

	return &plugin.ListResult{
		Plugins: plugins,
	}, nil
}

// GetInfo returns detailed information about a specific plugin
func (s *PluginService) GetInfo(ctx context.Context, payload *plugin.GetInfoPayload) (*plugin.GetInfoResult, error) {
	// Validate plugin name
	if payload.PluginName == "" {
		return nil, fmt.Errorf("plugin name cannot be empty")
	}

	// Mock plugin database
	mockPlugins := map[string]*plugin.GetInfoResult{
		"aws-cost-source": {
			Name:        "aws-cost-source",
			Version:     "v1.0.0",
			Description: stringPtr("AWS Cost and Usage Report data source"),
			Capabilities: &plugin.PluginCapabilities{
				SupportsProjected: false,
				SupportsActual:    true,
				SupportsProviders: []string{"aws"},
			},
			HealthStatus: &plugin.HealthStatus{
				Status:    "healthy",
				LastCheck: stringPtr(time.Now().Add(-5 * time.Minute).Format(time.RFC3339)),
				LatencyMs: int64Ptr(12),
			},
			GrpcAddress: stringPtr("localhost:50051"),
			Configuration: map[string]any{
				"region":       "us-east-1",
				"bucket":       "my-cur-bucket",
				"report_name":  "cost-usage-report",
				"poll_interval": "1h",
			},
		},
		"azure-cost-source": {
			Name:        "azure-cost-source",
			Version:     "v1.0.0",
			Description: stringPtr("Azure Cost Management data source"),
			Capabilities: &plugin.PluginCapabilities{
				SupportsProjected: false,
				SupportsActual:    true,
				SupportsProviders: []string{"azure"},
			},
			HealthStatus: &plugin.HealthStatus{
				Status:    "healthy",
				LastCheck: stringPtr(time.Now().Add(-3 * time.Minute).Format(time.RFC3339)),
				LatencyMs: int64Ptr(18),
			},
			GrpcAddress: stringPtr("localhost:50052"),
			Configuration: map[string]any{
				"subscription_id": "12345678-1234-1234-1234-123456789012",
				"poll_interval":  "30m",
			},
		},
		"infracost-plugin": {
			Name:        "infracost-plugin",
			Version:     "v0.10.0",
			Description: stringPtr("Infracost-based projected cost estimation"),
			Capabilities: &plugin.PluginCapabilities{
				SupportsProjected: true,
				SupportsActual:    false,
				SupportsProviders: []string{"aws", "azure", "gcp"},
			},
			HealthStatus: &plugin.HealthStatus{
				Status:    "healthy",
				LastCheck: stringPtr(time.Now().Add(-1 * time.Minute).Format(time.RFC3339)),
				LatencyMs: int64Ptr(25),
			},
			GrpcAddress: stringPtr("localhost:50053"),
			Configuration: map[string]any{
				"api_key":       "***",
				"cache_enabled": true,
				"currency":      "USD",
			},
		},
	}

	result, exists := mockPlugins[payload.PluginName]
	if !exists {
		return nil, &plugin.NotFoundError{
			Message:  fmt.Sprintf("plugin '%s' not found", payload.PluginName),
			Resource: &payload.PluginName,
		}
	}

	return result, nil
}

// Validate runs conformance tests on a plugin
func (s *PluginService) Validate(ctx context.Context, payload *plugin.ValidatePayload) (*plugin.PluginValidationReport, error) {
	// Validate inputs
	if payload.PluginPath == "" {
		return nil, fmt.Errorf("plugin path cannot be empty")
	}

	// Validate conformance level
	validLevels := map[string]bool{"BASIC": true, "STANDARD": true, "FULL": true}
	if !validLevels[payload.ConformanceLevel] {
		return nil, fmt.Errorf("invalid conformance level: %s (must be BASIC, STANDARD, or FULL)", payload.ConformanceLevel)
	}

	// Mock validation tests
	testResults := []*plugin.ValidationTest{
		{
			Name:       "gRPC Interface",
			Passed:     true,
			DurationMs: int64Ptr(15),
		},
		{
			Name:       "Health Check Endpoint",
			Passed:     true,
			DurationMs: int64Ptr(10),
		},
		{
			Name:       "Cost Query Basic",
			Passed:     true,
			DurationMs: int64Ptr(25),
		},
	}

	// Add more tests for higher conformance levels
	if payload.ConformanceLevel == "STANDARD" || payload.ConformanceLevel == "FULL" {
		testResults = append(testResults, &plugin.ValidationTest{
			Name:       "Resource Filtering",
			Passed:     true,
			DurationMs: int64Ptr(20),
		})
		testResults = append(testResults, &plugin.ValidationTest{
			Name:       "Time Range Queries",
			Passed:     true,
			DurationMs: int64Ptr(18),
		})
	}

	if payload.ConformanceLevel == "FULL" {
		testResults = append(testResults, &plugin.ValidationTest{
			Name:       "Granularity Support",
			Passed:     true,
			DurationMs: int64Ptr(22),
		})
		testResults = append(testResults, &plugin.ValidationTest{
			Name:       "Tag-based Queries",
			Passed:     true,
			DurationMs: int64Ptr(19),
		})
	}

	// All tests passed
	allPassed := true
	for _, test := range testResults {
		if !test.Passed {
			allPassed = false
			break
		}
	}

	return &plugin.PluginValidationReport{
		PluginName:       payload.PluginPath,
		ConformanceLevel: payload.ConformanceLevel,
		Passed:           allPassed,
		TestResults:      testResults,
		Timestamp:        stringPtr(time.Now().Format(time.RFC3339)),
	}, nil
}

// HealthCheck checks plugin health and connectivity
func (s *PluginService) HealthCheck(ctx context.Context, payload *plugin.HealthCheckPayload) (*plugin.HealthStatus, error) {
	// Validate plugin name
	if payload.PluginName == "" {
		return nil, fmt.Errorf("plugin name cannot be empty")
	}

	// Mock plugin health status database
	mockHealthStatus := map[string]*plugin.HealthStatus{
		"aws-cost-source": {
			Status:    "healthy",
			LastCheck: stringPtr(time.Now().Format(time.RFC3339)),
			LatencyMs: int64Ptr(12),
		},
		"azure-cost-source": {
			Status:    "healthy",
			LastCheck: stringPtr(time.Now().Format(time.RFC3339)),
			LatencyMs: int64Ptr(18),
		},
		"infracost-plugin": {
			Status:    "healthy",
			LastCheck: stringPtr(time.Now().Format(time.RFC3339)),
			LatencyMs: int64Ptr(25),
		},
		"slow-plugin": {
			Status:       "degraded",
			LastCheck:    stringPtr(time.Now().Format(time.RFC3339)),
			LatencyMs:    int64Ptr(500),
			ErrorMessage: stringPtr("High latency detected"),
		},
	}

	status, exists := mockHealthStatus[payload.PluginName]
	if !exists {
		return nil, &plugin.NotFoundError{
			Message:  fmt.Sprintf("plugin '%s' not found", payload.PluginName),
			Resource: &payload.PluginName,
		}
	}

	return status, nil
}

// Helper functions
func int64Ptr(i int64) *int64 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
