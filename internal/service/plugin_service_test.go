package service

import (
	"context"
	"testing"

	"github.com/rshade/pulumicost-mcp/gen/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestList tests listing all available plugins
func TestList(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.ListPayload{
		IncludeHealth: false,
	}

	result, err := service.List(ctx, payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Plugins)

	// Plugins list may be empty if no plugins installed
	// This is expected behavior - just verify the call succeeds
	t.Logf("Found %d plugins", len(result.Plugins))

	// If plugins exist, verify structure
	if len(result.Plugins) > 0 {
		firstPlugin := result.Plugins[0]
		assert.NotEmpty(t, firstPlugin.Name)
		assert.NotEmpty(t, firstPlugin.Version)
	}
}

// TestList_WithHealthCheck tests listing plugins with health status
func TestList_WithHealthCheck(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.ListPayload{
		IncludeHealth: true,
	}

	result, err := service.List(ctx, payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Plugins)

	t.Logf("Found %d plugins with health check", len(result.Plugins))

	// When health check is included and plugins exist, verify health status is populated
	if len(result.Plugins) > 0 {
		firstPlugin := result.Plugins[0]
		assert.NotNil(t, firstPlugin.HealthStatus)
		if firstPlugin.HealthStatus != nil {
			assert.NotEmpty(t, firstPlugin.HealthStatus.Status)
		}
	}
}

// TestGetInfo tests getting detailed plugin information
func TestGetInfo(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.GetInfoPayload{
		PluginName: "aws-cost-source",
	}

	result, err := service.GetInfo(ctx, payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "aws-cost-source", result.Name)
	assert.NotEmpty(t, result.Version)
	assert.NotNil(t, result.Capabilities)
	assert.True(t, result.Capabilities.SupportsActual)
	assert.Contains(t, result.Capabilities.SupportsProviders, "aws")
}

// TestGetInfo_NotFound tests getting info for non-existent plugin
func TestGetInfo_NotFound(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.GetInfoPayload{
		PluginName: "nonexistent-plugin",
	}

	result, err := service.GetInfo(ctx, payload)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

// TestGetInfo_EmptyName tests validation of plugin name
func TestGetInfo_EmptyName(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.GetInfoPayload{
		PluginName: "",
	}

	result, err := service.GetInfo(ctx, payload)

	require.Error(t, err)
	assert.Nil(t, result)
}

// TestValidate tests plugin conformance validation
func TestValidate(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.ValidatePayload{
		PluginPath:       "/path/to/plugin",
		ConformanceLevel: "basic", // Use lowercase - spec adapter expects lowercase
	}

	result, err := service.Validate(ctx, payload)

	// Validation may return error if plugin doesn't exist (expected)
	// But result should still be populated
	require.NotNil(t, result)
	assert.Equal(t, "/path/to/plugin", result.PluginName)
	assert.Equal(t, "basic", result.ConformanceLevel)

	// Error is expected since plugin doesn't actually exist
	if err == nil {
		t.Log("Validation succeeded (plugin must exist)")
	} else {
		t.Logf("Validation returned expected error: %v", err)
		assert.False(t, result.Passed, "should not pass validation without running plugin")
	}
}

// TestValidate_InvalidPath tests validation with invalid path
func TestValidate_InvalidPath(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.ValidatePayload{
		PluginPath:       "",
		ConformanceLevel: "STANDARD",
	}

	result, err := service.Validate(ctx, payload)

	require.Error(t, err)
	assert.Nil(t, result)
}

// TestValidate_InvalidConformanceLevel tests validation with invalid conformance level
func TestValidate_InvalidConformanceLevel(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.ValidatePayload{
		PluginPath:       "/path/to/plugin",
		ConformanceLevel: "INVALID",
	}

	result, err := service.Validate(ctx, payload)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid conformance level")
}

// TestHealthCheck tests plugin health check
func TestHealthCheck(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.HealthCheckPayload{
		PluginName: "aws-cost-source",
	}

	result, err := service.HealthCheck(ctx, payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Status)
	assert.Contains(t, []string{"healthy", "unhealthy", "degraded"}, result.Status)

	if result.Status == "healthy" {
		assert.NotNil(t, result.LatencyMs)
		assert.Greater(t, *result.LatencyMs, int64(0))
	}
}

// TestHealthCheck_NotFound tests health check for non-existent plugin
func TestHealthCheck_NotFound(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.HealthCheckPayload{
		PluginName: "nonexistent-plugin",
	}

	result, err := service.HealthCheck(ctx, payload)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

// TestHealthCheck_EmptyName tests validation of plugin name
func TestHealthCheck_EmptyName(t *testing.T) {
	service := NewPluginService("/tmp/plugins", nil)
	ctx := context.Background()

	payload := &plugin.HealthCheckPayload{
		PluginName: "",
	}

	result, err := service.HealthCheck(ctx, payload)

	require.Error(t, err)
	assert.Nil(t, result)
}
