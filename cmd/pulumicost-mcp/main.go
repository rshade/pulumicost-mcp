package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/rshade/pulumicost-mcp/internal/adapter"
	"github.com/rshade/pulumicost-mcp/internal/config"
	"github.com/rshade/pulumicost-mcp/internal/logging"
	"github.com/rshade/pulumicost-mcp/internal/metrics"
	"github.com/rshade/pulumicost-mcp/internal/service"
	"github.com/rshade/pulumicost-mcp/internal/tracing"
	mcpcost "github.com/rshade/pulumicost-mcp/gen/mcp_cost"
	mcpplugin "github.com/rshade/pulumicost-mcp/gen/mcp_plugin"
	mcpanalysis "github.com/rshade/pulumicost-mcp/gen/mcp_analysis"
	costsvr "github.com/rshade/pulumicost-mcp/gen/jsonrpc/mcp_cost/server"
	pluginsvr "github.com/rshade/pulumicost-mcp/gen/jsonrpc/mcp_plugin/server"
	analysissvr "github.com/rshade/pulumicost-mcp/gen/jsonrpc/mcp_analysis/server"
	goahttp "goa.design/goa/v3/http"
)

func main() {
	// Setup structured logger
	logger := logging.New(logging.Config{
		Level:  "info",
		Format: "json",
		Output: os.Stderr,
	})
	stdLogger := log.New(os.Stderr, "[pulumicost-mcp] ", log.Ltime|log.Lshortfile)

	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		stdLogger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize tracing
	shutdownTracing, err := tracing.Init(tracing.Config{
		ServiceName:    "pulumicost-mcp",
		ServiceVersion: "0.1.0",
		Environment:    "development",
		Enabled:        true,
	})
	if err != nil {
		stdLogger.Fatalf("Failed to initialize tracing: %v", err)
	}
	defer shutdownTracing(context.Background())

	logger.Info("observability initialized")

	// Create PulumiCost adapter
	pulumiAdapter := adapter.NewPulumiCostAdapter(cfg.PulumiCost.CorePath)
	logger.Info("pulumicost adapter initialized", "core_path", cfg.PulumiCost.CorePath)

	// Create services
	costService := service.NewCostService(pulumiAdapter, logger)
	// Plugin directory - default to ~/.pulumicost/plugins or config value
	pluginDir := os.Getenv("PLUGIN_DIR")
	if pluginDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fall back to /tmp if home dir unavailable
			logger.Warn("failed to get home directory, using /tmp", "error", err)
			pluginDir = "/tmp/pulumicost/plugins"
		} else {
			pluginDir = filepath.Join(homeDir, ".pulumicost", "plugins")
		}
	}

	pluginService := service.NewPluginService(pluginDir, logger)
	analysisService := service.NewAnalysisService(nil, logger)
	logger.Info("services initialized")

	// Create MCP adapters
	mcpCostAdapter := mcpcost.NewMCPAdapter(costService, nil)
	mcpPluginAdapter := mcpplugin.NewMCPAdapter(pluginService, nil)
	mcpAnalysisAdapter := mcpanalysis.NewMCPAdapter(analysisService, nil)
	logger.Info("mcp adapters initialized")

	// Create MCP endpoints
	mcpCostEndpoints := mcpcost.NewEndpoints(mcpCostAdapter)
	mcpPluginEndpoints := mcpplugin.NewEndpoints(mcpPluginAdapter)
	mcpAnalysisEndpoints := mcpanalysis.NewEndpoints(mcpAnalysisAdapter)

	// Create HTTP muxer
	mux := goahttp.NewMuxer()

	// Add metrics endpoint
	mux.Handle("GET", "/metrics", metrics.Handler().ServeHTTP)

	// Mount JSON-RPC servers for MCP services
	costServer := costsvr.New(mcpCostEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil)
	costsvr.Mount(mux, costServer)

	pluginServer := pluginsvr.New(mcpPluginEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil)
	pluginsvr.Mount(mux, pluginServer)

	analysisServer := analysissvr.New(mcpAnalysisEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil)
	analysissvr.Mount(mux, analysisServer)

	logger.Info("mcp services mounted")

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("starting server", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stdLogger.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server exited")
}
