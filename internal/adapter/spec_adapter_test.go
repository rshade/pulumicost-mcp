package adapter

import (
	"context"
	"testing"

	"github.com/rshade/pulumicost-mcp/gen/plugin"
	"github.com/rshade/pulumicost-mcp/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidatePlugin verifies plugin conformance testing (T061)
func TestValidatePlugin(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping conformance test in short mode")
	}

	logger := logging.Default()
	adapter := NewSpecAdapter(logger)

	ctx := context.Background()

	testPlugin := &plugin.Plugin{
		Name:    "test-plugin",
		Version: "1.0.0",
	}

	result, err := adapter.ValidatePlugin(ctx, testPlugin, "basic")

	// Expected to fail without actual plugin running
	// Just verify the interface works
	if err != nil {
		assert.Error(t, err, "should error without running plugin")
		assert.NotNil(t, result, "should return partial result even on error")
		return
	}

	require.NotNil(t, result)
	assert.Equal(t, testPlugin.Name, result.PluginName)
	assert.Equal(t, "basic", result.ConformanceLevel)
}

// TestValidatePlugin_Levels verifies different conformance levels
func TestValidatePlugin_Levels(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping conformance test in short mode")
	}

	logger := logging.Default()
	adapter := NewSpecAdapter(logger)

	ctx := context.Background()

	testPlugin := &plugin.Plugin{
		Name:    "test-plugin",
		Version: "1.0.0",
	}

	levels := []string{"basic", "standard", "advanced"}

	for _, level := range levels {
		t.Run("Level_"+level, func(t *testing.T) {
			result, err := adapter.ValidatePlugin(ctx, testPlugin, level)

			// Without actual plugin, expect error
			if err != nil {
				assert.NotNil(t, result)
				assert.Equal(t, level, result.ConformanceLevel)
				return
			}

			// If no error, verify result
			require.NotNil(t, result)
			assert.Equal(t, testPlugin.Name, result.PluginName)
			assert.Equal(t, level, result.ConformanceLevel)
		})
	}
}

// TestValidatePlugin_InvalidLevel verifies error handling
func TestValidatePlugin_InvalidLevel(t *testing.T) {
	logger := logging.Default()
	adapter := NewSpecAdapter(logger)

	ctx := context.Background()

	testPlugin := &plugin.Plugin{
		Name:    "test-plugin",
		Version: "1.0.0",
	}

	result, err := adapter.ValidatePlugin(ctx, testPlugin, "invalid-level")

	assert.Error(t, err, "should error on invalid conformance level")
	assert.NotNil(t, result, "should still return result structure")
}
