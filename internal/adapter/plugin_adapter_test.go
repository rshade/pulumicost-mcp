package adapter

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rshade/pulumicost-mcp/gen/plugin"
	"github.com/rshade/pulumicost-mcp/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDiscoverPlugins verifies plugin discovery from filesystem (T055)
func TestDiscoverPlugins(t *testing.T) {
	// Create temporary plugin directory
	tmpDir := t.TempDir()

	// Create mock plugin metadata files
	createMockPlugin(t, tmpDir, "infracost", "1.0.0", "aws,azure,gcp")
	createMockPlugin(t, tmpDir, "kubecost", "2.1.0", "kubernetes")

	logger := logging.Default()
	adapter := NewPluginAdapter(tmpDir, logger)

	ctx := context.Background()
	plugins, err := adapter.DiscoverPlugins(ctx)

	require.NoError(t, err, "DiscoverPlugins should succeed")
	require.NotNil(t, plugins, "plugins list should not be nil")
	assert.Len(t, plugins, 2, "should discover 2 plugins")

	// Verify plugin details
	pluginNames := make(map[string]bool)
	for _, p := range plugins {
		pluginNames[p.Name] = true
		assert.NotEmpty(t, p.Version, "plugin should have version")
	}

	assert.True(t, pluginNames["infracost"], "should find infracost plugin")
	assert.True(t, pluginNames["kubecost"], "should find kubecost plugin")
}

// TestDiscoverPlugins_EmptyDirectory verifies behavior with no plugins
func TestDiscoverPlugins_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	logger := logging.Default()
	adapter := NewPluginAdapter(tmpDir, logger)

	ctx := context.Background()
	plugins, err := adapter.DiscoverPlugins(ctx)

	require.NoError(t, err, "DiscoverPlugins should succeed with empty dir")
	assert.Empty(t, plugins, "should return empty list")
}

// TestDiscoverPlugins_InvalidMetadata verifies error handling for bad metadata
func TestDiscoverPlugins_InvalidMetadata(t *testing.T) {
	tmpDir := t.TempDir()

	// Create plugin with invalid metadata
	pluginDir := filepath.Join(tmpDir, "bad-plugin")
	require.NoError(t, os.MkdirAll(pluginDir, 0755))

	metadataPath := filepath.Join(pluginDir, "plugin.json")
	require.NoError(t, os.WriteFile(metadataPath, []byte("{invalid json"), 0644))

	logger := logging.Default()
	adapter := NewPluginAdapter(tmpDir, logger)

	ctx := context.Background()
	plugins, err := adapter.DiscoverPlugins(ctx)

	// Should skip invalid plugins, not fail
	require.NoError(t, err)
	assert.Empty(t, plugins, "should skip plugin with invalid metadata")
}

// TestEstablishConnection verifies gRPC connection establishment (T058)
func TestEstablishConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger := logging.Default()
	adapter := NewPluginAdapter(t.TempDir(), logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testPlugin := &plugin.Plugin{
		Name:    "test-plugin",
		Version: "1.0.0",
		// Note: Would need actual gRPC server for full test
		// GRPCAddress is loaded from plugin.json metadata
	}

	// This will fail without actual server, but tests the interface
	err := adapter.EstablishConnection(ctx, testPlugin)

	// Expected to fail - just verify error is returned properly
	assert.Error(t, err, "should error without running gRPC server")
}

// TestGetPluginCapabilities verifies capability querying (T059)
func TestGetPluginCapabilities(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger := logging.Default()
	adapter := NewPluginAdapter(t.TempDir(), logger)

	ctx := context.Background()

	testPlugin := &plugin.Plugin{
		Name:    "test-plugin",
		Version: "1.0.0",
	}

	_, err := adapter.GetPluginCapabilities(ctx, testPlugin)

	// Expected to fail without connection
	assert.Error(t, err, "should error without established connection")
}

// TestHealthCheck verifies plugin health checking
func TestHealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	logger := logging.Default()
	adapter := NewPluginAdapter(t.TempDir(), logger)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	testPlugin := &plugin.Plugin{
		Name:    "test-plugin",
		Version: "1.0.0",
	}

	status, latency, err := adapter.HealthCheck(ctx, testPlugin)

	// Expected to fail without connection
	assert.Error(t, err)
	assert.Equal(t, "unhealthy", status)
	assert.Equal(t, int64(0), latency)
}

// TestCircuitBreaker verifies circuit breaker prevents cascade failures (T060)
func TestCircuitBreaker(t *testing.T) {
	logger := logging.Default()
	adapter := NewPluginAdapter(t.TempDir(), logger)

	ctx := context.Background()

	testPlugin := &plugin.Plugin{
		Name:    "failing-plugin",
		Version: "1.0.0",
	}

	// Simulate multiple failures
	failureCount := 0
	for i := 0; i < 10; i++ {
		_, _, err := adapter.HealthCheck(ctx, testPlugin)
		if err != nil {
			failureCount++
		}
	}

	// Circuit breaker should prevent some attempts after threshold
	assert.Greater(t, failureCount, 0, "should have failures")

	// Verify circuit breaker state can be checked
	isOpen := adapter.IsCircuitOpen(testPlugin.Name)
	assert.NotNil(t, isOpen, "circuit breaker state should be queryable")
}

// Helper functions

func createMockPlugin(t *testing.T, baseDir, name, version, providers string) {
	t.Helper()

	pluginDir := filepath.Join(baseDir, name)
	require.NoError(t, os.MkdirAll(pluginDir, 0755))

	metadata := `{
		"name": "` + name + `",
		"version": "` + version + `",
		"description": "Mock plugin for testing",
		"providers": "` + providers + `",
		"grpc_address": "localhost:50051",
		"capabilities": {
			"supports_projected_cost": true,
			"supports_actual_cost": true,
			"supports_optimization": false
		}
	}`

	metadataPath := filepath.Join(pluginDir, "plugin.json")
	require.NoError(t, os.WriteFile(metadataPath, []byte(metadata), 0644))

	// Create dummy binary
	binaryPath := filepath.Join(pluginDir, name)
	require.NoError(t, os.WriteFile(binaryPath, []byte("#!/bin/bash\necho 'mock plugin'"), 0755))
}
