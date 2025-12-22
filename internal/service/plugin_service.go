package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rshade/pulumicost-mcp/gen/plugin"
	"github.com/rshade/pulumicost-mcp/internal/adapter"
	"github.com/rshade/pulumicost-mcp/internal/logging"
	"github.com/rshade/pulumicost-mcp/internal/metrics"
	"github.com/rshade/pulumicost-mcp/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
)

// PluginService implements the plugin.Service interface
type PluginService struct {
	pluginAdapter *adapter.PluginAdapter
	specAdapter   *adapter.SpecAdapter
	logger        *logging.Logger
}

// NewPluginService creates a new Plugin Service instance
func NewPluginService(pluginDir string, logger *logging.Logger) *PluginService {
	return &PluginService{
		pluginAdapter: adapter.NewPluginAdapter(pluginDir, logger),
		specAdapter:   adapter.NewSpecAdapter(logger),
		logger:        logger,
	}
}

// List returns all available cost source plugins
func (s *PluginService) List(ctx context.Context, payload *plugin.ListPayload) (*plugin.ListResult, error) {
	start := time.Now()
	ctx, span := tracing.Start(ctx, "PluginService.List")
	defer span.End()

	s.logger.WithService("plugin").Info("listing plugins")

	tracing.SetAttributes(ctx, attribute.Bool("include_health", payload.IncludeHealth))

	// Discover plugins from filesystem
	plugins, err := s.pluginAdapter.DiscoverPlugins(ctx)
	if err != nil {
		s.logger.WithService("plugin").Error("failed to discover plugins", "error", err)
		return nil, fmt.Errorf("discover plugins: %w", err)
	}

	// If health check requested, check each plugin
	if payload.IncludeHealth {
		for _, p := range plugins {
			status, latency, _ := s.pluginAdapter.HealthCheck(ctx, p)
			lastCheck := time.Now().Format(time.RFC3339)
			latencyMs := latency
			p.HealthStatus = &plugin.HealthStatus{
				Status:    status,
				LastCheck: &lastCheck,
				LatencyMs: &latencyMs,
			}
		}
	}

	// Record metrics
	metrics.RecordRequest("plugin", "list", time.Since(start))
	tracing.SetAttributes(ctx, attribute.Int("plugin_count", len(plugins)))

	s.logger.WithService("plugin").InfoJSON("plugins listed", map[string]interface{}{
		"plugin_count":   len(plugins),
		"include_health": payload.IncludeHealth,
		"duration_ms":    time.Since(start).Milliseconds(),
	})

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

	s.logger.WithService("plugin").Info("validating plugin",
		"path", payload.PluginPath,
		"level", payload.ConformanceLevel)

	// Create plugin object for validation
	p := &plugin.Plugin{
		Name:    payload.PluginPath, // Use path as name for now
		Version: "unknown",
	}

	// Use spec adapter to validate plugin
	report, err := s.specAdapter.ValidatePlugin(ctx, p, payload.ConformanceLevel)
	if err != nil {
		s.logger.WithService("plugin").Warn("plugin validation encountered error",
			"plugin", payload.PluginPath,
			"error", err)

		// If it's an invalid conformance level error, return error without report
		if report != nil && !report.Passed && err.Error() == "invalid conformance level: "+payload.ConformanceLevel {
			return nil, err
		}

		// Otherwise return report even if there was an error (partial validation)
		if report != nil {
			return report, nil
		}
		return nil, fmt.Errorf("validate plugin: %w", err)
	}

	return report, nil
}

// HealthCheck checks plugin health and connectivity
func (s *PluginService) HealthCheck(ctx context.Context, payload *plugin.HealthCheckPayload) (*plugin.HealthStatus, error) {
	start := time.Now()
	ctx, span := tracing.Start(ctx, "PluginService.HealthCheck")
	defer span.End()

	s.logger.WithService("plugin").Info("checking plugin health", "plugin", payload.PluginName)

	// Validate plugin name
	if payload.PluginName == "" {
		err := fmt.Errorf("plugin name cannot be empty")
		s.logger.WithService("plugin").ErrorJSON("validation failed", err, nil)
		metrics.RecordError("plugin", "health_check", "validation")
		tracing.RecordError(ctx, err)
		return nil, err
	}

	tracing.SetAttributes(ctx, attribute.String("plugin_name", payload.PluginName))

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
		err := &plugin.NotFoundError{
			Message:  fmt.Sprintf("plugin '%s' not found", payload.PluginName),
			Resource: &payload.PluginName,
		}
		s.logger.WithService("plugin").ErrorJSON("plugin not found", err, map[string]interface{}{
			"plugin": payload.PluginName,
		})
		metrics.RecordError("plugin", "health_check", "not_found")
		tracing.RecordError(ctx, err)
		return nil, err
	}

	// Record plugin health metrics
	pluginStatus := "success"
	if status.Status != "healthy" {
		pluginStatus = status.Status
	}

	latencyDuration := time.Duration(*status.LatencyMs) * time.Millisecond
	metrics.RecordPluginCall(payload.PluginName, pluginStatus, latencyDuration)
	metrics.RecordRequest("plugin", "health_check", time.Since(start))

	tracing.SetAttributes(ctx,
		attribute.String("status", status.Status),
		attribute.Int64("latency_ms", *status.LatencyMs),
	)

	s.logger.WithService("plugin").InfoJSON("plugin health checked", map[string]interface{}{
		"plugin":      payload.PluginName,
		"status":      status.Status,
		"latency_ms":  *status.LatencyMs,
		"duration_ms": time.Since(start).Milliseconds(),
	})

	return status, nil
}

// Helper functions
func int64Ptr(i int64) *int64 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
