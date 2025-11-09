package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_ValidConfig(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configYAML := `
server:
  port: 8080
  host: "localhost"
  log_level: "info"
  log_format: "json"
  read_timeout: 30s
  write_timeout: 30s
  shutdown_timeout: 10s

pulumicost:
  core_path: "/usr/local/bin/pulumicost"
  plugin_dir: "~/.pulumicost/plugins"
  spec_version: "0.1.0"
  batch_size: 100

plugins:
  timeout: 30s
  max_concurrent: 10
  health_check_interval: 60s
  retry_attempts: 3
  retry_delay: 5s

mcp:
  enable_streaming: true
  max_message_size: 10485760
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
    path: "/metrics"
  tracing:
    enabled: false
    endpoint: "http://jaeger:14268/api/traces"
    sample_rate: 0.1
  logging:
    format: "json"
    level: "info"
    output: "stdout"

security:
  tls:
    enabled: false
  auth:
    enabled: false
    method: "token"
    token_header: "X-API-Token"

pulumi:
  backend_url: "https://api.pulumi.com"

features:
  experimental: false
  forecasting: true
  anomaly_detection: true
  recommendations: true

rate_limiting:
  enabled: false
  requests_per_minute: 60
  burst: 10

cors:
  enabled: false
  allowed_origins:
    - "http://localhost:3000"
  allowed_methods:
    - "GET"
    - "POST"
  allowed_headers:
    - "Content-Type"
  exposed_headers:
    - "X-Request-ID"
  max_age: "3600"

development:
  debug: false
  profiling: false
  hot_reload: false
`

	err := os.WriteFile(configPath, []byte(configYAML), 0644)
	require.NoError(t, err)

	// Load config
	cfg, err := Load(configPath)
	require.NoError(t, err)

	// Assert values
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "info", cfg.Server.LogLevel)
	assert.Equal(t, "/usr/local/bin/pulumicost", cfg.PulumiCost.CorePath)
	assert.Equal(t, 100, cfg.PulumiCost.BatchSize)
	assert.Equal(t, 10, cfg.Plugins.MaxConcurrent)
	assert.Equal(t, true, cfg.MCP.EnableStreaming)
	assert.Equal(t, int64(10485760), cfg.MCP.MaxMessageSize)
	assert.Equal(t, true, cfg.Features.Forecasting)
}

func TestLoad_InvalidPath(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read config file")
}

func TestLoad_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	invalidYAML := `
server:
  port: invalid_port
  host: localhost
`

	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	_, err = Load(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse config file")
}

func TestLoad_EnvOverrides(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configYAML := `
server:
  port: 8080
  host: "localhost"
  log_level: "info"

pulumicost:
  core_path: "/usr/local/bin/pulumicost"
  plugin_dir: "~/.pulumicost/plugins"
  batch_size: 100

pulumi:
  backend_url: "https://api.pulumi.com"
`

	err := os.WriteFile(configPath, []byte(configYAML), 0644)
	require.NoError(t, err)

	// Set environment variables
	require.NoError(t, os.Setenv("PULUMI_ACCESS_TOKEN", "pul-test-token"))
	require.NoError(t, os.Setenv("MCP_HOST", "0.0.0.0"))
	require.NoError(t, os.Setenv("MCP_PORT", "9090"))
	require.NoError(t, os.Setenv("MCP_LOG_LEVEL", "debug"))
	require.NoError(t, os.Setenv("PULUMICOST_CORE_PATH", "/custom/path/pulumicost"))
	require.NoError(t, os.Setenv("PULUMICOST_PLUGIN_DIR", "/custom/plugins"))
	defer func() {
		os.Unsetenv("PULUMI_ACCESS_TOKEN")
		os.Unsetenv("MCP_HOST")
		os.Unsetenv("MCP_PORT")
		os.Unsetenv("MCP_LOG_LEVEL")
		os.Unsetenv("PULUMICOST_CORE_PATH")
		os.Unsetenv("PULUMICOST_PLUGIN_DIR")
	}()

	// Load config
	cfg, err := Load(configPath)
	require.NoError(t, err)

	// Assert environment overrides
	assert.Equal(t, "pul-test-token", cfg.Pulumi.AccessToken)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Server.LogLevel)
	assert.Equal(t, "/custom/path/pulumicost", cfg.PulumiCost.CorePath)
	assert.Equal(t, "/custom/plugins", cfg.PulumiCost.PluginDir)
}

func TestValidate_Valid(t *testing.T) {
	cfg := Default()
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestValidate_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"port too low", 0},
		{"port too high", 70000},
		{"negative port", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Default()
			cfg.Server.Port = tt.port
			err := cfg.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid server port")
		})
	}
}

func TestValidate_InvalidLogLevel(t *testing.T) {
	cfg := Default()
	cfg.Server.LogLevel = "invalid"
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid log level")
}

func TestValidate_MissingCorePath(t *testing.T) {
	cfg := Default()
	cfg.PulumiCost.CorePath = ""
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "core_path is required")
}

func TestValidate_MissingPluginDir(t *testing.T) {
	cfg := Default()
	cfg.PulumiCost.PluginDir = ""
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plugin_dir is required")
}

func TestValidate_InvalidBatchSize(t *testing.T) {
	cfg := Default()
	cfg.PulumiCost.BatchSize = 0
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "batch_size must be at least 1")
}

func TestValidate_InvalidMaxConcurrent(t *testing.T) {
	cfg := Default()
	cfg.Plugins.MaxConcurrent = 0
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max_concurrent must be at least 1")
}

func TestValidate_NegativeRetryAttempts(t *testing.T) {
	cfg := Default()
	cfg.Plugins.RetryAttempts = -1
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retry_attempts cannot be negative")
}

func TestValidate_InvalidMaxMessageSize(t *testing.T) {
	cfg := Default()
	cfg.MCP.MaxMessageSize = 512
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max_message_size must be at least 1024 bytes")
}

func TestValidate_InvalidMetricsPort(t *testing.T) {
	cfg := Default()
	cfg.Observability.Metrics.Enabled = true
	cfg.Observability.Metrics.Port = 0
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid metrics port")
}

func TestValidate_TLSMissingCert(t *testing.T) {
	cfg := Default()
	cfg.Security.TLS.Enabled = true
	cfg.Security.TLS.CertFile = ""
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cert_file is required")
}

func TestValidate_TLSMissingKey(t *testing.T) {
	cfg := Default()
	cfg.Security.TLS.Enabled = true
	cfg.Security.TLS.CertFile = "/path/to/cert.pem"
	cfg.Security.TLS.KeyFile = ""
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key_file is required")
}

func TestValidate_InvalidAuthMethod(t *testing.T) {
	cfg := Default()
	cfg.Security.Auth.Enabled = true
	cfg.Security.Auth.Method = "invalid"
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid auth method")
}

func TestValidate_ValidAuthMethods(t *testing.T) {
	validMethods := []string{"token", "oauth", "mtls"}

	for _, method := range validMethods {
		t.Run(method, func(t *testing.T) {
			cfg := Default()
			cfg.Security.Auth.Enabled = true
			cfg.Security.Auth.Method = method
			err := cfg.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestDefault(t *testing.T) {
	cfg := Default()

	// Assert default values
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "info", cfg.Server.LogLevel)
	assert.Equal(t, "json", cfg.Server.LogFormat)
	assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, 10*time.Second, cfg.Server.ShutdownTimeout)

	assert.Equal(t, "/usr/local/bin/pulumicost", cfg.PulumiCost.CorePath)
	assert.Equal(t, "~/.pulumicost/plugins", cfg.PulumiCost.PluginDir)
	assert.Equal(t, "0.1.0", cfg.PulumiCost.SpecVersion)
	assert.Equal(t, 100, cfg.PulumiCost.BatchSize)

	assert.Equal(t, 30*time.Second, cfg.Plugins.Timeout)
	assert.Equal(t, 10, cfg.Plugins.MaxConcurrent)
	assert.Equal(t, 60*time.Second, cfg.Plugins.HealthCheckInterval)
	assert.Equal(t, 3, cfg.Plugins.RetryAttempts)
	assert.Equal(t, 5*time.Second, cfg.Plugins.RetryDelay)

	assert.True(t, cfg.MCP.EnableStreaming)
	assert.Equal(t, int64(10*1024*1024), cfg.MCP.MaxMessageSize)
	assert.Equal(t, 30*time.Second, cfg.MCP.ConnectionTimeout)

	assert.True(t, cfg.Cache.Enabled)
	assert.Equal(t, 5*time.Minute, cfg.Cache.TTL.PluginMetadata)
	assert.Equal(t, 30*time.Second, cfg.Cache.TTL.CostData)
	assert.Equal(t, 1*time.Minute, cfg.Cache.TTL.PulumiState)

	assert.True(t, cfg.Observability.Metrics.Enabled)
	assert.Equal(t, 9090, cfg.Observability.Metrics.Port)
	assert.Equal(t, "/metrics", cfg.Observability.Metrics.Path)

	assert.False(t, cfg.Observability.Tracing.Enabled)
	assert.Equal(t, 0.1, cfg.Observability.Tracing.SampleRate)

	assert.False(t, cfg.Security.TLS.Enabled)
	assert.False(t, cfg.Security.Auth.Enabled)
	assert.Equal(t, "token", cfg.Security.Auth.Method)
	assert.Equal(t, "X-API-Token", cfg.Security.Auth.TokenHeader)

	assert.Equal(t, "https://api.pulumi.com", cfg.Pulumi.BackendURL)

	assert.False(t, cfg.Features.Experimental)
	assert.True(t, cfg.Features.Forecasting)
	assert.True(t, cfg.Features.AnomalyDetection)
	assert.True(t, cfg.Features.Recommendations)

	assert.False(t, cfg.RateLimiting.Enabled)
	assert.Equal(t, 60, cfg.RateLimiting.RequestsPerMinute)
	assert.Equal(t, 10, cfg.RateLimiting.Burst)

	assert.False(t, cfg.CORS.Enabled)
	assert.Contains(t, cfg.CORS.AllowedOrigins, "http://localhost:3000")
	assert.Contains(t, cfg.CORS.AllowedMethods, "POST")
	assert.Contains(t, cfg.CORS.AllowedHeaders, "Content-Type")

	assert.False(t, cfg.Development.Debug)
	assert.False(t, cfg.Development.Profiling)
	assert.False(t, cfg.Development.HotReload)

	// Ensure default config passes validation
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestLoad_ValidationFailure(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Config with invalid values
	configYAML := `
server:
  port: 0
  host: "localhost"

pulumicost:
  core_path: "/usr/local/bin/pulumicost"
  plugin_dir: "~/.pulumicost/plugins"
  batch_size: 100
`

	err := os.WriteFile(configPath, []byte(configYAML), 0644)
	require.NoError(t, err)

	_, err = Load(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validate config")
}

func TestLoad_TimeoutParsing(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configYAML := `
server:
  port: 8080
  host: "localhost"
  read_timeout: 45s
  write_timeout: 1m
  shutdown_timeout: 15s

pulumicost:
  core_path: "/usr/local/bin/pulumicost"
  plugin_dir: "~/.pulumicost/plugins"
  batch_size: 100
`

	err := os.WriteFile(configPath, []byte(configYAML), 0644)
	require.NoError(t, err)

	cfg, err := Load(configPath)
	require.NoError(t, err)

	assert.Equal(t, 45*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 1*time.Minute, cfg.Server.WriteTimeout)
	assert.Equal(t, 15*time.Second, cfg.Server.ShutdownTimeout)
}
