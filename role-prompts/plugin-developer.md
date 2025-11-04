# Plugin Developer - PulumiCost MCP Server

## Role Context

You are a Plugin Developer building cost source plugins for the PulumiCost ecosystem. Your plugins integrate various cost data sources (Kubecost, Vantage, cloud provider APIs) with the PulumiCost core system via the pulumicost-spec gRPC protocol.

## Key Responsibilities

- **Plugin Implementation**: Build spec-compliant cost source plugins
- **Integration**: Integrate with cost data APIs
- **Testing**: Ensure conformance to specification
- **Documentation**: Document plugin capabilities and usage
- **Performance**: Optimize plugin response times
- **Maintenance**: Update plugins for spec changes

## Plugin Development Overview

### What is a Cost Source Plugin?

A plugin is a standalone gRPC server that implements the `CostSourceService` interface defined in pulumicost-spec. It translates between the standard PulumiCost interface and a specific cost data source.

### Plugin Architecture

```
┌─────────────────────────────────────────┐
│       PulumiCost Core/MCP Server        │
└────────────────┬────────────────────────┘
                 │ gRPC
     ┌───────────┼───────────┐
     │           │           │
┌────▼────┐ ┌───▼────┐ ┌───▼────┐
│ Kubecost│ │Vantage │ │  AWS   │
│ Plugin  │ │ Plugin │ │ Plugin │
└────┬────┘ └───┬────┘ └───┬────┘
     │          │          │
┌────▼────┐ ┌───▼────┐ ┌───▼────┐
│Kubecost │ │Vantage │ │AWS Cost│
│   API   │ │  API   │ │Explorer│
└─────────┘ └────────┘ └────────┘
```

## Getting Started

### Prerequisites

```bash
# Go 1.24+
go version

# Pulumicost-spec
go get github.com/rshade/pulumicost-spec/sdk/go/proto
go get github.com/rshade/pulumicost-spec/sdk/go/types

# Protocol Buffers compiler
brew install protobuf

# Optional: Buf for proto management
brew install buf
```

### Plugin Scaffolding

```bash
# Use the MCP server to generate scaffold
# (Once MCP server is running)
curl -X POST http://localhost:8080/rpc -d '{
  "jsonrpc": "2.0",
  "method": "generate_plugin",
  "params": {
    "plugin_name": "my-cost-source",
    "provider": "custom",
    "billing_models": ["per_hour", "per_gb_month"]
  },
  "id": 1
}'

# Or manually create structure
mkdir -p my-plugin/{cmd,internal/{service,client},pkg/types}
```

## Implementation Guide

### Step 1: Implement CostSourceService

```go
// internal/service/cost_service.go
package service

import (
    "context"
    "fmt"
    "time"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/timestamppb"

    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

type MyCostService struct {
    pbc.UnimplementedCostSourceServiceServer
    client *MyAPIClient  // Your API client
}

// Name returns the plugin name
func (s *MyCostService) Name(ctx context.Context,
    req *pbc.NameRequest) (*pbc.NameResponse, error) {

    return &pbc.NameResponse{
        Name: "my-cost-source",
    }, nil
}

// Supports checks if a resource type is supported
func (s *MyCostService) Supports(ctx context.Context,
    req *pbc.SupportsRequest) (*pbc.SupportsResponse, error) {

    resource := req.GetResource()

    // Check if we support this provider and resource type
    if resource.GetProvider() == "aws" &&
       resource.GetResourceType() == "ec2/instance" {
        return &pbc.SupportsResponse{
            Supported: true,
        }, nil
    }

    return &pbc.SupportsResponse{
        Supported: false,
        Reason:    fmt.Sprintf("Provider %s type %s not supported",
            resource.GetProvider(), resource.GetResourceType()),
    }, nil
}

// GetActualCost retrieves historical cost data
func (s *MyCostService) GetActualCost(ctx context.Context,
    req *pbc.GetActualCostRequest) (*pbc.GetActualCostResponse, error) {

    // Extract parameters
    resourceID := req.GetResourceId()
    start := req.GetStart().AsTime()
    end := req.GetEnd().AsTime()
    tags := req.GetTags()

    // Validate inputs
    if resourceID == "" {
        return nil, status.Error(codes.InvalidArgument,
            "resource_id is required")
    }

    if end.Before(start) {
        return nil, status.Error(codes.InvalidArgument,
            "end time must be after start time")
    }

    // Query your cost API
    costs, err := s.client.GetCosts(ctx, resourceID, start, end, tags)
    if err != nil {
        return nil, status.Errorf(codes.Internal,
            "failed to fetch costs: %v", err)
    }

    // Convert to spec format
    results := make([]*pbc.ActualCostResult, len(costs))
    for i, cost := range costs {
        results[i] = &pbc.ActualCostResult{
            Timestamp:   timestamppb.New(cost.Timestamp),
            Cost:        cost.Amount,
            Currency:    cost.Currency,
            UsageAmount: cost.Usage,
            UsageUnit:   cost.UsageUnit,
            Source:      "my-cost-source",
            Tags:        cost.Tags,
        }
    }

    return &pbc.GetActualCostResponse{
        Results: results,
    }, nil
}

// GetProjectedCost calculates estimated costs
func (s *MyCostService) GetProjectedCost(ctx context.Context,
    req *pbc.GetProjectedCostRequest) (*pbc.GetProjectedCostResponse, error) {

    resource := req.GetResource()
    config := req.GetConfiguration()

    // Extract resource properties
    instanceType := config["instance_type"]
    region := resource.GetRegion()

    // Get pricing data
    pricing, err := s.client.GetPricing(ctx, instanceType, region)
    if err != nil {
        return nil, status.Errorf(codes.Internal,
            "failed to fetch pricing: %v", err)
    }

    // Calculate monthly cost (assuming 730 hours/month)
    monthlyHours := 730.0
    monthlyCost := pricing.HourlyRate * monthlyHours

    return &pbc.GetProjectedCostResponse{
        UnitPrice:     pricing.HourlyRate,
        Currency:      "USD",
        CostPerMonth:  monthlyCost,
        BillingDetail: fmt.Sprintf("%s in %s", instanceType, region),
    }, nil
}

// GetPricingSpec returns detailed pricing specification
func (s *MyCostService) GetPricingSpec(ctx context.Context,
    req *pbc.GetPricingSpecRequest) (*pbc.GetPricingSpecResponse, error) {

    resource := req.GetResource()

    // Build pricing specification
    spec := &pbc.PricingSpec{
        Provider:     resource.GetProvider(),
        ResourceType: resource.GetResourceType(),
        Sku:          req.GetSku(),
        Region:       resource.GetRegion(),
        BillingMode:  "per_hour",
        RatePerUnit:  0.0104,  // Example rate
        Currency:     "USD",
        Description:  "EC2 t3.micro hourly rate",
        MetricHints: []*pbc.UsageMetricHint{
            {
                Metric:            "compute_hours",
                Unit:              "hour",
                AggregationMethod: "sum",
                SamplingInterval:  3600,  // 1 hour
            },
        },
        PluginMetadata: map[string]string{
            "api_version": "v1",
            "last_updated": time.Now().Format(time.RFC3339),
        },
        Source: "my-cost-source",
    }

    return &pbc.GetPricingSpecResponse{
        Spec: spec,
    }, nil
}
```

### Step 2: Create gRPC Server

```go
// cmd/my-plugin/main.go
package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/health"
    "google.golang.org/grpc/health/grpc_health_v1"
    "google.golang.org/grpc/reflection"

    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
    "my-plugin/internal/service"
)

var (
    port = flag.Int("port", 50051, "The server port")
)

func main() {
    flag.Parse()

    // Create listener
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    // Create gRPC server
    server := grpc.NewServer(
        grpc.UnaryInterceptor(loggingInterceptor),
    )

    // Register service
    costService := service.NewMyCostService()
    pbc.RegisterCostSourceServiceServer(server, costService)

    // Register health check
    healthServer := health.NewServer()
    grpc_health_v1.RegisterHealthServer(server, healthServer)
    healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

    // Register reflection (for debugging)
    reflection.Register(server)

    // Graceful shutdown
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
        <-sigCh

        log.Println("Shutting down gracefully...")
        healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        stopped := make(chan struct{})
        go func() {
            server.GracefulStop()
            close(stopped)
        }()

        select {
        case <-ctx.Done():
            server.Stop()
        case <-stopped:
        }

        log.Println("Server stopped")
        os.Exit(0)
    }()

    log.Printf("Server listening on :%d", *port)
    if err := server.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}

func loggingInterceptor(ctx context.Context, req interface{},
    info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

    start := time.Now()
    resp, err := handler(ctx, req)
    duration := time.Since(start)

    log.Printf("Method: %s, Duration: %v, Error: %v",
        info.FullMethod, duration, err)

    return resp, err
}
```

### Step 3: Add Configuration

```go
// internal/config/config.go
package config

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    API      APIConfig      `yaml:"api"`
    Cache    CacheConfig    `yaml:"cache"`
    Logging  LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
    Port            int    `yaml:"port"`
    ShutdownTimeout string `yaml:"shutdown_timeout"`
}

type APIConfig struct {
    Endpoint string `yaml:"endpoint"`
    APIKey   string `yaml:"api_key"`
    Timeout  string `yaml:"timeout"`
}

type CacheConfig struct {
    Enabled bool   `yaml:"enabled"`
    TTL     string `yaml:"ttl"`
}

type LoggingConfig struct {
    Level  string `yaml:"level"`
    Format string `yaml:"format"`
}

func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config: %w", err)
    }

    // Load secrets from environment
    if apiKey := os.Getenv("API_KEY"); apiKey != "" {
        cfg.API.APIKey = apiKey
    }

    return &cfg, nil
}
```

```yaml
# config.yaml
server:
  port: 50051
  shutdown_timeout: 10s

api:
  endpoint: https://api.mycostprovider.com
  timeout: 30s

cache:
  enabled: true
  ttl: 5m

logging:
  level: info
  format: json
```

### Step 4: Testing

```go
// internal/service/cost_service_test.go
package service

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "google.golang.org/protobuf/types/known/timestamppb"

    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
    plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

func TestName(t *testing.T) {
    service := NewMyCostService()

    resp, err := service.Name(context.Background(), &pbc.NameRequest{})

    require.NoError(t, err)
    assert.Equal(t, "my-cost-source", resp.GetName())
}

func TestSupports(t *testing.T) {
    service := NewMyCostService()

    tests := []struct {
        name      string
        resource  *pbc.ResourceDescriptor
        supported bool
    }{
        {
            name: "supported resource",
            resource: &pbc.ResourceDescriptor{
                Provider:     "aws",
                ResourceType: "ec2/instance",
            },
            supported: true,
        },
        {
            name: "unsupported resource",
            resource: &pbc.ResourceDescriptor{
                Provider:     "azure",
                ResourceType: "vm",
            },
            supported: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := service.Supports(context.Background(), &pbc.SupportsRequest{
                Resource: tt.resource,
            })

            require.NoError(t, err)
            assert.Equal(t, tt.supported, resp.GetSupported())
        })
    }
}

func TestGetActualCost(t *testing.T) {
    service := NewMyCostService()

    start := time.Now().Add(-24 * time.Hour)
    end := time.Now()

    resp, err := service.GetActualCost(context.Background(), &pbc.GetActualCostRequest{
        ResourceId: "i-1234567890",
        Start:      timestamppb.New(start),
        End:        timestamppb.New(end),
    })

    require.NoError(t, err)
    assert.NotEmpty(t, resp.GetResults())

    for _, result := range resp.GetResults() {
        assert.Greater(t, result.GetCost(), 0.0)
        assert.NotEmpty(t, result.GetCurrency())
        assert.NotNil(t, result.GetTimestamp())
    }
}

// Conformance testing
func TestConformance(t *testing.T) {
    service := NewMyCostService()

    // Run basic conformance tests
    result := plugintesting.RunBasicConformanceTests(t, service)
    plugintesting.PrintConformanceReport(result)

    if result.FailedTests > 0 {
        t.Errorf("Failed %d conformance tests", result.FailedTests)
    }
}
```

## Plugin Testing

### Manual Testing

```bash
# Start plugin
go run cmd/my-plugin/main.go --port 50051

# Test with grpcurl
grpcurl -plaintext localhost:50051 list

# Call Name method
grpcurl -plaintext localhost:50051 \
  pulumicost.v1.CostSourceService/Name

# Call GetActualCost
grpcurl -plaintext -d '{
  "resource_id": "i-1234567890",
  "start": {"seconds": 1704067200},
  "end": {"seconds": 1704153600}
}' localhost:50051 \
  pulumicost.v1.CostSourceService/GetActualCost
```

### Conformance Testing

```bash
# Run conformance tests
go test -v ./... -run TestConformance

# Test specific level
go test -v ./... -run TestBasicConformance
go test -v ./... -run TestStandardConformance
go test -v ./... -run TestAdvancedConformance
```

### Integration Testing

```bash
# Use MCP server's test_plugin tool
curl -X POST http://localhost:8080/rpc -d '{
  "jsonrpc": "2.0",
  "method": "test_plugin",
  "params": {
    "plugin_path": "/path/to/my-plugin",
    "test_level": "standard"
  },
  "id": 1
}'
```

## Plugin Distribution

### Directory Structure

```
~/.pulumicost/plugins/
└── my-cost-source/
    └── 1.0.0/
        ├── pulumicost-my-cost-source  # Binary
        ├── config.yaml               # Default config
        └── README.md                 # Documentation
```

### Installation Script

```bash
#!/bin/bash
# install.sh

PLUGIN_NAME="my-cost-source"
VERSION="1.0.0"
PLUGIN_DIR="$HOME/.pulumicost/plugins/$PLUGIN_NAME/$VERSION"

mkdir -p "$PLUGIN_DIR"

# Build plugin
go build -o "$PLUGIN_DIR/pulumicost-$PLUGIN_NAME" cmd/my-plugin/main.go

# Copy config
cp config.yaml "$PLUGIN_DIR/"

# Copy docs
cp README.md "$PLUGIN_DIR/"

echo "Plugin installed to $PLUGIN_DIR"
```

### Plugin Manifest

```yaml
# plugin.yaml
name: my-cost-source
version: 1.0.0
description: Cost source plugin for MyProvider
author: Your Name <your@email.com>
license: Apache-2.0
homepage: https://github.com/yourusername/pulumicost-plugin-my-source

capabilities:
  providers:
    - aws
    - azure
  billing_modes:
    - per_hour
    - per_gb_month
  features:
    - actual_costs
    - projected_costs
    - pricing_specs

requirements:
  pulumicost_spec_version: ">=0.1.0"
  go_version: ">=1.24"

configuration:
  api_endpoint:
    required: true
    description: "API endpoint URL"
  api_key:
    required: true
    secret: true
    description: "API authentication key"
```

## Best Practices

### Error Handling

```go
// Use appropriate gRPC status codes
import "google.golang.org/grpc/codes"
import "google.golang.org/grpc/status"

// Invalid input
return nil, status.Error(codes.InvalidArgument, "resource_id required")

// Not found
return nil, status.Error(codes.NotFound, "resource not found")

// API errors
return nil, status.Errorf(codes.Internal, "API error: %v", err)

// Timeout
return nil, status.Error(codes.DeadlineExceeded, "request timeout")

// Unsupported
return nil, status.Error(codes.Unimplemented, "feature not supported")
```

### Performance Optimization

```go
// 1. Implement caching
type CachedClient struct {
    client *APIClient
    cache  *ttlcache.Cache
}

func (c *CachedClient) GetPricing(ctx context.Context, sku string) (*Pricing, error) {
    // Check cache
    if cached, ok := c.cache.Get(sku); ok {
        return cached.(*Pricing), nil
    }

    // Fetch from API
    pricing, err := c.client.GetPricing(ctx, sku)
    if err != nil {
        return nil, err
    }

    // Cache result
    c.cache.Set(sku, pricing, 5*time.Minute)
    return pricing, nil
}

// 2. Batch API calls
func (c *APIClient) GetCostsBatch(ctx context.Context,
    resourceIDs []string) (map[string]*Cost, error) {
    // Single API call for multiple resources
}

// 3. Use connection pooling
client := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:       10,
        IdleConnTimeout:    90 * time.Second,
        DisableCompression: true,
    },
}
```

### Security

```go
// 1. Validate all inputs
func validateResourceID(id string) error {
    if id == "" {
        return fmt.Errorf("resource_id cannot be empty")
    }
    if len(id) > 256 {
        return fmt.Errorf("resource_id too long")
    }
    // Add more validation
    return nil
}

// 2. Use TLS for API calls
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
}

// 3. Sanitize log output
log.Printf("Fetching costs for resource: %s", sanitize(resourceID))

// 4. Handle secrets securely
apiKey := os.Getenv("API_KEY")  // Don't hardcode
```

## Common Patterns

### Rate Limiting

```go
import "golang.org/x/time/rate"

type RateLimitedClient struct {
    client  *APIClient
    limiter *rate.Limiter
}

func (c *RateLimitedClient) GetCosts(ctx context.Context, id string) (*Cost, error) {
    if err := c.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    return c.client.GetCosts(ctx, id)
}
```

### Retry Logic

```go
import "github.com/cenkalti/backoff/v4"

func (c *APIClient) GetCostsWithRetry(ctx context.Context, id string) (*Cost, error) {
    var cost *Cost

    operation := func() error {
        var err error
        cost, err = c.GetCosts(ctx, id)
        return err
    }

    bo := backoff.NewExponentialBackOff()
    bo.MaxElapsedTime = 30 * time.Second

    if err := backoff.Retry(operation, bo); err != nil {
        return nil, err
    }

    return cost, nil
}
```

### Circuit Breaker

```go
import "github.com/sony/gobreaker"

type CircuitBreakerClient struct {
    client  *APIClient
    breaker *gobreaker.CircuitBreaker
}

func (c *CircuitBreakerClient) GetCosts(ctx context.Context, id string) (*Cost, error) {
    result, err := c.breaker.Execute(func() (interface{}, error) {
        return c.client.GetCosts(ctx, id)
    })

    if err != nil {
        return nil, err
    }

    return result.(*Cost), nil
}
```

## Resources

- [pulumicost-spec Documentation](https://github.com/rshade/pulumicost-spec)
- [Example Plugins](../examples/plugins/)
- [Plugin Development Guide](../docs/guides/plugin-development.md)
- [Conformance Testing](https://github.com/rshade/pulumicost-spec/blob/main/sdk/go/testing/README.md)

---

**Remember**: A good plugin is reliable, performant, and well-tested. Focus on correctness first, then optimize. Your plugin represents a critical data source in the cost analysis pipeline.
