// Package config provides configuration loading and validation for the PulumiCost MCP server.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete application configuration
type Config struct {
	Server        ServerConfig        `yaml:"server"`
	PulumiCost    PulumiCostConfig    `yaml:"pulumicost"`
	Plugins       PluginsConfig       `yaml:"plugins"`
	MCP           MCPConfig           `yaml:"mcp"`
	Cache         CacheConfig         `yaml:"cache"`
	Observability ObservabilityConfig `yaml:"observability"`
	Security      SecurityConfig      `yaml:"security"`
	Pulumi        PulumiConfig        `yaml:"pulumi"`
	Features      FeaturesConfig      `yaml:"features"`
	RateLimiting  RateLimitingConfig  `yaml:"rate_limiting"`
	CORS          CORSConfig          `yaml:"cors"`
	Development   DevelopmentConfig   `yaml:"development"`
}

// ServerConfig defines server listening and timeout settings
type ServerConfig struct {
	Port            int           `yaml:"port"`
	Host            string        `yaml:"host"`
	LogLevel        string        `yaml:"log_level"`
	LogFormat       string        `yaml:"log_format"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// PulumiCostConfig defines pulumicost-core integration settings
type PulumiCostConfig struct {
	CorePath    string `yaml:"core_path"`
	PluginDir   string `yaml:"plugin_dir"`
	SpecVersion string `yaml:"spec_version"`
	BatchSize   int    `yaml:"batch_size"`
}

// PluginsConfig defines plugin management settings
type PluginsConfig struct {
	Timeout             time.Duration `yaml:"timeout"`
	MaxConcurrent       int           `yaml:"max_concurrent"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
	RetryAttempts       int           `yaml:"retry_attempts"`
	RetryDelay          time.Duration `yaml:"retry_delay"`
}

// MCPConfig defines MCP protocol settings
type MCPConfig struct {
	EnableStreaming   bool          `yaml:"enable_streaming"`
	MaxMessageSize    int64         `yaml:"max_message_size"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout"`
}

// CacheConfig defines caching settings
type CacheConfig struct {
	Enabled bool     `yaml:"enabled"`
	TTL     CacheTTL `yaml:"ttl"`
}

// CacheTTL defines TTL for different cache types
type CacheTTL struct {
	PluginMetadata time.Duration `yaml:"plugin_metadata"`
	CostData       time.Duration `yaml:"cost_data"`
	PulumiState    time.Duration `yaml:"pulumi_state"`
}

// ObservabilityConfig defines metrics, tracing, and logging settings
type ObservabilityConfig struct {
	Metrics MetricsConfig `yaml:"metrics"`
	Tracing TracingConfig `yaml:"tracing"`
	Logging LoggingConfig `yaml:"logging"`
}

// MetricsConfig defines Prometheus metrics settings
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

// TracingConfig defines distributed tracing settings
type TracingConfig struct {
	Enabled    bool    `yaml:"enabled"`
	Endpoint   string  `yaml:"endpoint"`
	SampleRate float64 `yaml:"sample_rate"`
}

// LoggingConfig defines logging format and output settings
type LoggingConfig struct {
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

// SecurityConfig defines TLS and authentication settings
type SecurityConfig struct {
	TLS  TLSConfig  `yaml:"tls"`
	Auth AuthConfig `yaml:"auth"`
}

// TLSConfig defines TLS certificate settings
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// AuthConfig defines authentication settings
type AuthConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Method      string `yaml:"method"`
	TokenHeader string `yaml:"token_header"`
}

// PulumiConfig defines Pulumi-specific settings
type PulumiConfig struct {
	AccessToken  string `yaml:"access_token"`
	BackendURL   string `yaml:"backend_url"`
	DefaultStack string `yaml:"default_stack"`
}

// FeaturesConfig defines feature flags
type FeaturesConfig struct {
	Experimental      bool `yaml:"experimental"`
	Forecasting       bool `yaml:"forecasting"`
	AnomalyDetection  bool `yaml:"anomaly_detection"`
	Recommendations   bool `yaml:"recommendations"`
}

// RateLimitingConfig defines rate limiting settings
type RateLimitingConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute"`
	Burst             int  `yaml:"burst"`
}

// CORSConfig defines CORS settings
type CORSConfig struct {
	Enabled        bool     `yaml:"enabled"`
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
	ExposedHeaders []string `yaml:"exposed_headers"`
	MaxAge         string   `yaml:"max_age"`
}

// DevelopmentConfig defines development mode settings
type DevelopmentConfig struct {
	Debug      bool `yaml:"debug"`
	Profiling  bool `yaml:"profiling"`
	HotReload  bool `yaml:"hot_reload"`
}

// Load reads configuration from a YAML file and applies environment variable overrides
func Load(path string) (*Config, error) {
	// Start with default configuration
	cfg := Default()

	// If path is provided, read and merge config file
	if path != "" {
		// Read file
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}

		// Parse YAML and merge with defaults
		if unmarshalErr := yaml.Unmarshal(data, cfg); unmarshalErr != nil {
			return nil, fmt.Errorf("parse config file: %w", unmarshalErr)
		}
	}

	// Apply environment variable overrides for sensitive values
	applyEnvOverrides(cfg)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

// applyEnvOverrides applies environment variable overrides for sensitive configuration
func applyEnvOverrides(cfg *Config) {
	// Pulumi access token
	if token := os.Getenv("PULUMI_ACCESS_TOKEN"); token != "" {
		cfg.Pulumi.AccessToken = token
	}

	// Server host/port overrides
	if host := os.Getenv("MCP_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("MCP_PORT"); port != "" {
		var portNum int
		if _, err := fmt.Sscanf(port, "%d", &portNum); err == nil {
			cfg.Server.Port = portNum
		}
	}

	// Log level override
	if level := os.Getenv("MCP_LOG_LEVEL"); level != "" {
		cfg.Server.LogLevel = level
		cfg.Observability.Logging.Level = level
	}

	// PulumiCost core path override
	if corePath := os.Getenv("PULUMICOST_CORE_PATH"); corePath != "" {
		cfg.PulumiCost.CorePath = corePath
	}

	// Plugin directory override
	if pluginDir := os.Getenv("PULUMICOST_PLUGIN_DIR"); pluginDir != "" {
		cfg.PulumiCost.PluginDir = pluginDir
	}
}

// Validate checks that required configuration values are set and valid
// nolint:gocognit // sequential validation checks - complexity is inherent to comprehensive validation
func (c *Config) Validate() error {
	// Validate server config
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d (must be 1-65535)", c.Server.Port)
	}

	if c.Server.LogLevel != "" {
		validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
		if !validLevels[c.Server.LogLevel] {
			return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.Server.LogLevel)
		}
	}

	// Validate PulumiCost config
	if c.PulumiCost.CorePath == "" {
		return fmt.Errorf("pulumicost.core_path is required")
	}

	if c.PulumiCost.PluginDir == "" {
		return fmt.Errorf("pulumicost.plugin_dir is required")
	}

	if c.PulumiCost.BatchSize < 1 {
		return fmt.Errorf("pulumicost.batch_size must be at least 1")
	}

	// Validate plugins config
	if c.Plugins.MaxConcurrent < 1 {
		return fmt.Errorf("plugins.max_concurrent must be at least 1")
	}

	if c.Plugins.RetryAttempts < 0 {
		return fmt.Errorf("plugins.retry_attempts cannot be negative")
	}

	// Validate MCP config
	if c.MCP.MaxMessageSize < 1024 {
		return fmt.Errorf("mcp.max_message_size must be at least 1024 bytes")
	}

	// Validate metrics config
	if c.Observability.Metrics.Enabled && (c.Observability.Metrics.Port < 1 || c.Observability.Metrics.Port > 65535) {
		return fmt.Errorf("invalid metrics port: %d (must be 1-65535)", c.Observability.Metrics.Port)
	}

	// Validate TLS config
	if c.Security.TLS.Enabled {
		if c.Security.TLS.CertFile == "" {
			return fmt.Errorf("security.tls.cert_file is required when TLS is enabled")
		}
		if c.Security.TLS.KeyFile == "" {
			return fmt.Errorf("security.tls.key_file is required when TLS is enabled")
		}
	}

	// Validate auth config
	if c.Security.Auth.Enabled {
		validMethods := map[string]bool{"token": true, "oauth": true, "mtls": true}
		if !validMethods[c.Security.Auth.Method] {
			return fmt.Errorf("invalid auth method: %s (must be token, oauth, or mtls)", c.Security.Auth.Method)
		}
	}

	return nil
}

// Default returns a Config with sensible default values
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            8080,
			Host:            "0.0.0.0",
			LogLevel:        "info",
			LogFormat:       "json",
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		PulumiCost: PulumiCostConfig{
			CorePath:    "/usr/local/bin/pulumicost",
			PluginDir:   "~/.pulumicost/plugins",
			SpecVersion: "0.1.0",
			BatchSize:   100,
		},
		Plugins: PluginsConfig{
			Timeout:             30 * time.Second,
			MaxConcurrent:       10,
			HealthCheckInterval: 60 * time.Second,
			RetryAttempts:       3,
			RetryDelay:          5 * time.Second,
		},
		MCP: MCPConfig{
			EnableStreaming:   true,
			MaxMessageSize:    10 * 1024 * 1024, // 10MB
			ConnectionTimeout: 30 * time.Second,
		},
		Cache: CacheConfig{
			Enabled: true,
			TTL: CacheTTL{
				PluginMetadata: 5 * time.Minute,
				CostData:       30 * time.Second,
				PulumiState:    1 * time.Minute,
			},
		},
		Observability: ObservabilityConfig{
			Metrics: MetricsConfig{
				Enabled: true,
				Port:    9090,
				Path:    "/metrics",
			},
			Tracing: TracingConfig{
				Enabled:    false,
				Endpoint:   "http://jaeger:14268/api/traces",
				SampleRate: 0.1,
			},
			Logging: LoggingConfig{
				Format: "json",
				Level:  "info",
				Output: "stdout",
			},
		},
		Security: SecurityConfig{
			TLS: TLSConfig{
				Enabled: false,
			},
			Auth: AuthConfig{
				Enabled:     false,
				Method:      "token",
				TokenHeader: "X-API-Token",
			},
		},
		Pulumi: PulumiConfig{
			BackendURL: "https://api.pulumi.com",
		},
		Features: FeaturesConfig{
			Experimental:     false,
			Forecasting:      true,
			AnomalyDetection: true,
			Recommendations:  true,
		},
		RateLimiting: RateLimitingConfig{
			Enabled:           false,
			RequestsPerMinute: 60,
			Burst:             10,
		},
		CORS: CORSConfig{
			Enabled:        false,
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
			ExposedHeaders: []string{"X-Request-ID"},
			MaxAge:         "3600",
		},
		Development: DevelopmentConfig{
			Debug:     false,
			Profiling: false,
			HotReload: false,
		},
	}
}
