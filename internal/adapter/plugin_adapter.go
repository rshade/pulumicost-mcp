package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rshade/pulumicost-mcp/gen/plugin"
	"github.com/rshade/pulumicost-mcp/internal/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// PluginAdapter handles plugin discovery and gRPC communication
type PluginAdapter struct {
	pluginDir       string
	logger          *logging.Logger
	connections     map[string]*grpc.ClientConn
	connMutex       sync.RWMutex
	circuitBreakers map[string]*circuitBreaker
	cbMutex         sync.RWMutex
}

// circuitBreaker tracks plugin failures and prevents cascade failures
type circuitBreaker struct {
	failures      int
	lastFailure   time.Time
	state         string // "closed", "open", "half-open"
	threshold     int
	timeout       time.Duration
	resetInterval time.Duration
	mu            sync.RWMutex
}

// pluginMetadata represents the plugin.json structure
type pluginMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Providers   string `json:"providers"`
	GRPCAddress string `json:"grpc_address"`
	Capabilities struct {
		SupportsProjectedCost bool `json:"supports_projected_cost"`
		SupportsActualCost    bool `json:"supports_actual_cost"`
		SupportsOptimization  bool `json:"supports_optimization"`
	} `json:"capabilities"`
}

// NewPluginAdapter creates a new plugin adapter
func NewPluginAdapter(pluginDir string, logger *logging.Logger) *PluginAdapter {
	if logger == nil {
		logger = logging.Default()
	}
	return &PluginAdapter{
		pluginDir:       pluginDir,
		logger:          logger,
		connections:     make(map[string]*grpc.ClientConn),
		circuitBreakers: make(map[string]*circuitBreaker),
	}
}

// DiscoverPlugins scans the plugin directory and loads metadata (T056)
func (a *PluginAdapter) DiscoverPlugins(ctx context.Context) ([]*plugin.Plugin, error) {
	a.logger.Info("discovering plugins", "dir", a.pluginDir)

	// Check if plugin directory exists
	if _, err := os.Stat(a.pluginDir); os.IsNotExist(err) {
		a.logger.Warn("plugin directory does not exist", "dir", a.pluginDir)
		return []*plugin.Plugin{}, nil
	}

	entries, err := os.ReadDir(a.pluginDir)
	if err != nil {
		return nil, fmt.Errorf("read plugin directory: %w", err)
	}

	var plugins []*plugin.Plugin

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(a.pluginDir, entry.Name())
		metadataPath := filepath.Join(pluginPath, "plugin.json")

		// Check if metadata exists
		if _, statErr := os.Stat(metadataPath); os.IsNotExist(statErr) {
			a.logger.Debug("skipping directory without metadata", "dir", entry.Name())
			continue
		}

		// Load metadata
		data, readErr := os.ReadFile(metadataPath)
		if readErr != nil {
			a.logger.Warn("failed to read plugin metadata", "plugin", entry.Name(), "error", readErr)
			continue
		}

		var meta pluginMetadata
		if unmarshalErr := json.Unmarshal(data, &meta); unmarshalErr != nil {
			a.logger.Warn("failed to parse plugin metadata", "plugin", entry.Name(), "error", unmarshalErr)
			continue
		}

		// Convert to plugin type
		p := &plugin.Plugin{
			Name:    meta.Name,
			Version: meta.Version,
		}

		if meta.Description != "" {
			p.Description = &meta.Description
		}

		// Add to list
		plugins = append(plugins, p)

		a.logger.Info("discovered plugin", "name", meta.Name, "version", meta.Version)
	}

	return plugins, nil
}

// EstablishConnection establishes gRPC connection to a plugin (T058)
func (a *PluginAdapter) EstablishConnection(ctx context.Context, p *plugin.Plugin) error {
	a.connMutex.Lock()
	defer a.connMutex.Unlock()

	// Check if already connected
	if _, exists := a.connections[p.Name]; exists {
		return nil
	}

	// Check circuit breaker
	if a.isCircuitOpenLocked(p.Name) {
		return fmt.Errorf("circuit breaker open for plugin %s", p.Name)
	}

	// Load metadata to get gRPC address
	metadataPath := filepath.Join(a.pluginDir, p.Name, "plugin.json")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		a.recordFailure(p.Name)
		return fmt.Errorf("read plugin metadata: %w", err)
	}

	var meta pluginMetadata
	if unmarshalErr := json.Unmarshal(data, &meta); unmarshalErr != nil {
		a.recordFailure(p.Name)
		return fmt.Errorf("parse plugin metadata: %w", unmarshalErr)
	}

	// Establish gRPC connection
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// nolint:staticcheck // SA1019: grpc.DialContext is deprecated but supported throughout 1.x
	conn, err := grpc.DialContext(dialCtx, meta.GRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // nolint:staticcheck // SA1019: grpc.WithBlock deprecated but needed for blocking dial
	)
	if err != nil {
		a.recordFailure(p.Name)
		return fmt.Errorf("dial plugin %s at %s: %w", p.Name, meta.GRPCAddress, err)
	}

	a.connections[p.Name] = conn
	a.logger.Info("established connection to plugin", "name", p.Name, "address", meta.GRPCAddress)

	return nil
}

// GetPluginCapabilities queries plugin capabilities via gRPC (T059)
func (a *PluginAdapter) GetPluginCapabilities(ctx context.Context, p *plugin.Plugin) (*plugin.PluginCapabilities, error) {
	a.connMutex.RLock()
	conn, exists := a.connections[p.Name]
	a.connMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no connection to plugin %s", p.Name)
	}

	// Check circuit breaker
	if a.IsCircuitOpen(p.Name) {
		return nil, fmt.Errorf("circuit breaker open for plugin %s", p.Name)
	}

	// For now, return capabilities from metadata
	// In full implementation, this would query the plugin via gRPC
	metadataPath := filepath.Join(a.pluginDir, p.Name, "plugin.json")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("read plugin metadata: %w", err)
	}

	var meta pluginMetadata
	if unmarshalErr := json.Unmarshal(data, &meta); unmarshalErr != nil {
		return nil, fmt.Errorf("parse plugin metadata: %w", unmarshalErr)
	}

	capabilities := &plugin.PluginCapabilities{
		SupportsProjected: meta.Capabilities.SupportsProjectedCost,
		SupportsActual:    meta.Capabilities.SupportsActualCost,
	}

	// Use connection to verify it's alive
	_ = conn

	return capabilities, nil
}

// Health status constants
const (
	statusHealthy   = "healthy"
	statusUnhealthy = "unhealthy"
)

// HealthCheck performs health check on plugin (FR-017)
func (a *PluginAdapter) HealthCheck(ctx context.Context, p *plugin.Plugin) (status string, latency int64, err error) {
	// Check circuit breaker first
	if a.IsCircuitOpen(p.Name) {
		return statusUnhealthy, 0, fmt.Errorf("circuit breaker open")
	}

	a.connMutex.RLock()
	conn, exists := a.connections[p.Name]
	a.connMutex.RUnlock()

	if !exists {
		// Try to establish connection first
		if connErr := a.EstablishConnection(ctx, p); connErr != nil {
			return statusUnhealthy, 0, connErr
		}

		a.connMutex.RLock()
		conn, exists = a.connections[p.Name]
		a.connMutex.RUnlock()

		if !exists {
			return statusUnhealthy, 0, fmt.Errorf("failed to establish connection")
		}
	}

	// Perform health check via gRPC
	start := time.Now()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := healthClient.Check(checkCtx, &grpc_health_v1.HealthCheckRequest{
		Service: "",
	})

	latency = time.Since(start).Milliseconds()

	if err != nil {
		a.recordFailure(p.Name)
		return statusUnhealthy, latency, fmt.Errorf("health check failed: %w", err)
	}

	a.resetFailures(p.Name)

	if resp.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
		return statusHealthy, latency, nil
	}

	return statusUnhealthy, latency, fmt.Errorf("plugin not serving")
}

// IsCircuitOpen checks if circuit breaker is open for a plugin (T060)
func (a *PluginAdapter) IsCircuitOpen(pluginName string) bool {
	a.cbMutex.RLock()
	defer a.cbMutex.RUnlock()
	return a.isCircuitOpenLocked(pluginName)
}

func (a *PluginAdapter) isCircuitOpenLocked(pluginName string) bool {
	cb, exists := a.circuitBreakers[pluginName]
	if !exists {
		return false
	}

	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.state == "open" {
		// Check if timeout has passed to transition to half-open
		if time.Since(cb.lastFailure) > cb.timeout {
			return false // Transition to half-open
		}
		return true
	}

	return false
}

func (a *PluginAdapter) recordFailure(pluginName string) {
	a.cbMutex.Lock()
	defer a.cbMutex.Unlock()

	cb, exists := a.circuitBreakers[pluginName]
	if !exists {
		cb = &circuitBreaker{
			threshold:     5,                // Open after 5 failures
			timeout:       30 * time.Second, // Stay open for 30s
			resetInterval: 60 * time.Second, // Reset counter after 60s
			state:         "closed",
		}
		a.circuitBreakers[pluginName] = cb
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.failures >= cb.threshold {
		cb.state = "open"
		a.logger.Warn("circuit breaker opened", "plugin", pluginName, "failures", cb.failures)
	}
}

func (a *PluginAdapter) resetFailures(pluginName string) {
	a.cbMutex.Lock()
	defer a.cbMutex.Unlock()

	cb, exists := a.circuitBreakers[pluginName]
	if !exists {
		return
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Reset if plugin is healthy
	if time.Since(cb.lastFailure) > cb.resetInterval {
		cb.failures = 0
		cb.state = "closed"
	}
}

// Close closes all plugin connections
func (a *PluginAdapter) Close() error {
	a.connMutex.Lock()
	defer a.connMutex.Unlock()

	for name, conn := range a.connections {
		if err := conn.Close(); err != nil {
			a.logger.Warn("failed to close plugin connection", "plugin", name, "error", err)
		}
	}

	a.connections = make(map[string]*grpc.ClientConn)
	return nil
}
