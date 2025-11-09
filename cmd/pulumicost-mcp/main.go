package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rshade/pulumicost-mcp/internal/adapter"
	"github.com/rshade/pulumicost-mcp/internal/config"
	"github.com/rshade/pulumicost-mcp/internal/service"
	mcpcost "github.com/rshade/pulumicost-mcp/gen/mcp_cost"
	costsvr "github.com/rshade/pulumicost-mcp/gen/jsonrpc/mcp_cost/server"
	goahttp "goa.design/goa/v3/http"
)

func main() {
	// Setup logger
	logger := log.New(os.Stderr, "[pulumicost-mcp] ", log.Ltime|log.Lshortfile)

	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Create PulumiCost adapter
	pulumiAdapter := adapter.NewPulumiCostAdapter(cfg.PulumiCost.CorePath)
	logger.Printf("PulumiCost adapter initialized (core path: %s)", cfg.PulumiCost.CorePath)

	// Create Cost service
	costService := service.NewCostService(pulumiAdapter, logger)
	logger.Printf("Cost service initialized")

	// Create MCP adapter wrapping the Cost service
	mcpAdapter := mcpcost.NewMCPAdapter(costService, nil)
	logger.Printf("MCP adapter initialized")

	// Create MCP endpoints
	mcpEndpoints := mcpcost.NewEndpoints(mcpAdapter)

	// Create HTTP muxer
	mux := goahttp.NewMuxer()

	// Mount JSON-RPC server for MCP Cost service
	costServer := costsvr.New(mcpEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil)
	costsvr.Mount(mux, costServer)
	logger.Printf("MCP Cost service mounted on JSON-RPC")

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
		logger.Printf("Starting server on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Printf("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Printf("Server exited")
}
