# Quick Start: PulumiCost MCP Server

**Feature**: 001-mcp-server
**Status**: Implementation Planning Phase
**Prerequisites**: Go 1.24, make, Goa CLI tools

## Overview

This quick start guides developers through setting up a local development
environment and implementing the PulumiCost MCP Server following design-first
principles.

## Phase 1: Environment Setup (10 minutes)

### 1. Install Dependencies

```bash
# Install Go 1.24
# https://go.dev/dl/

# Verify installation
go version  # Should show go1.24.x

# Install Goa CLI
go install goa.design/goa/v3/cmd/goa@latest

# Install development tools
make install-tools
```

### 2. Verify Project Structure

```bash
# Confirm you're in the project root
pwd  # Should end with /pulumicost-mcp

# Verify directory structure
ls -la design/  # Should see .go files
ls -la internal/  # Should see service/ and adapter/ dirs
```

## Phase 2: Design-First Development (30 minutes)

### 1. Review Goa Design Files

Start with the design layer - the single source of truth:

```bash
# Main API configuration
cat design/design.go

# Cost Query Service definition
cat design/cost_service.go

# Type definitions
cat design/types.go
```

**Key Concept**: All APIs, types, and validation rules are defined in Goa DSL
before any implementation.

### 2. Generate Code from Design

```bash
# Generate transport, validation, and protocol code
make generate

# Verify generated code
ls -la gen/cost/  # Generated cost service
ls -la gen/mcp/   # Generated MCP protocol bindings
```

**Important**: Never modify files in `gen/` - they're regenerated from design.

### 3. Understand the Four-Layer Architecture

```text
1. Design Layer (design/*.go)
   └─► Defines: APIs, types, validation, MCP tools

2. Generated Layer (gen/*)
   └─► Generates: Transport, encoding, validation, docs

3. Service Layer (internal/service/*)
   └─► Implements: Business logic, orchestration

4. Adapter Layer (internal/adapter/*)
   └─► Integrates: pulumicost-core, plugins, external systems
```

## Phase 3: Test-First Implementation (1-2 hours)

### User Story 1: AI-Powered Cost Analysis (P1 - MVP)

#### Step 1: Write Failing Tests

```bash
# Create test file for Cost Service
cat > internal/service/cost_service_test.go << 'EOF'
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAnalyzeProjected(t *testing.T) {
    // Given: Pulumi preview JSON for AWS EC2 instance
    pulumiJSON := `{...}`  // Sample Pulumi JSON

    // When: Analyzing projected costs
    result, err := service.AnalyzeProjected(context.Background(), &cost.ProjectedCostRequest{
        PulumiJSON: pulumiJSON,
    })

    // Then: Should return cost breakdown
    require.NoError(t, err)
    assert.Greater(t, result.TotalMonthly, 0.0)
    assert.Equal(t, "USD", result.Currency)
    assert.NotEmpty(t, result.Resources)
}
EOF
```

#### Step 2: Run Tests (They Should Fail)

```bash
# Run tests - expect failures (RED)
make test

# Output should show: FAIL - function not implemented
```

#### Step 3: Implement Service Logic

```bash
# Edit service implementation
vim internal/service/cost_service.go

# Implement minimal logic to pass tests
# Use adapter to call pulumicost-core
```

#### Step 4: Run Tests Again (Should Pass)

```bash
# Run tests - expect success (GREEN)
make test

# All tests should pass
```

#### Step 5: Refactor

```bash
# Improve code quality while keeping tests green
# Add error handling, logging, validation
# Run tests after each change to ensure nothing breaks
```

## Phase 4: Integration with Adapters (1 hour)

### Implement PulumiCost Adapter

```bash
# Create adapter for pulumicost-core integration
cat > internal/adapter/pulumicost_adapter.go << 'EOF'
package adapter

import (
    "context"
    "os/exec"
)

type PulumiCostAdapter interface {
    GetProjectedCost(ctx context.Context, pulumiJSON string) (*CostResult, error)
}

type pulumiCostAdapter struct {
    corePath string
}

func NewPulumiCostAdapter(corePath string) PulumiCostAdapter {
    return &pulumiCostAdapter{corePath: corePath}
}

func (a *pulumiCostAdapter) GetProjectedCost(
    ctx context.Context,
    pulumiJSON string) (*CostResult, error) {
    // Write JSON to temp file
    // Execute pulumicost binary
    // Parse output
    // Return structured result
    return nil, nil  // TODO: Implement
}
EOF
```

### Wire Up Dependency Injection

```bash
# Update main.go to inject adapters
vim cmd/pulumicost-mcp/main.go

# Example:
# adapter := adapter.NewPulumiCostAdapter("/usr/local/bin/pulumicost")
# service := service.NewCostService(adapter, logger)
```

## Phase 5: Run and Test Locally (30 minutes)

### 1. Build the Server

```bash
# Compile binary
make build

# Binary location
./build/pulumicost-mcp
```

### 2. Configure Server

```bash
# Copy example config
cp config.yaml.example config.yaml

# Edit configuration
vim config.yaml

# Set paths:
# - pulumicost core binary path
# - plugin directory
# - server port
```

### 3. Run Server

```bash
# Start server
make run

# Or run binary directly
./build/pulumicost-mcp --config config.yaml

# Server should start on port 8080
# Logs should show: "MCP server listening on :8080"
```

### 4. Test with cURL

```bash
# Test JSON-RPC endpoint
curl -X POST http://localhost:8080/jsonrpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "cost.analyze_projected",
    "params": {
      "pulumi_json": "{...}"
    },
    "id": 1
  }'

# Should return cost analysis result
```

## Phase 6: Claude Desktop Integration (15 minutes)

### 1. Configure Claude Desktop

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
(path may be long)

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "pulumicost": {
      "command": "/path/to/pulumicost-mcp",
      "args": ["--config", "/path/to/config.yaml"],
      "env": {
        "PULUMI_ACCESS_TOKEN": "your-token"
      }
    }
  }
}
```

### 2. Restart Claude Desktop

```bash
# macOS
killall Claude && open -a Claude

# Windows
# Close and reopen Claude Desktop
```

### 3. Test in Claude

Ask Claude:

```text
What MCP tools are available?
```

Should show:

- analyze_projected_costs
- get_actual_costs
- compare_costs
- analyze_resource_cost
- query_cost_by_tags

### 4. Run Cost Analysis

Ask Claude:

```text
Here's my Pulumi preview JSON: {...}
What are the projected monthly costs?
```

Claude should use the MCP tool to analyze costs and provide a breakdown.

## Development Workflow Reference

### Daily Development Cycle

```bash
# 1. Start with design changes
vim design/cost_service.go

# 2. Regenerate code
make generate

# 3. Write failing tests
vim internal/service/cost_service_test.go

# 4. Implement functionality
vim internal/service/cost_service.go

# 5. Run tests
make test

# 6. Validate everything
make validate  # Runs lint + test

# 7. Commit design + implementation together
git add design/ internal/ gen/
git commit -m "feat(cost): add projected cost analysis"
```

### Common Commands

```bash
# Setup environment (first time)
make setup

# Generate from design
make generate

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linters
make lint

# Run full validation
make validate

# Build binary
make build

# Run server
make run

# Clean generated files
make clean
```

## Troubleshooting

### "make generate" fails

**Problem**: Goa code generation errors

**Solution**:

```bash
# Check Go version
go version  # Must be 1.24+

# Reinstall Goa
go install goa.design/goa/v3/cmd/goa@latest

# Verify design syntax
goa gen -h
```

### Tests fail after design change

**Problem**: Type mismatches after regenerating

**Solution**:

```bash
# Regenerate code
make generate

# Update service implementations to match new types
vim internal/service/*.go

# Update tests
vim internal/service/*_test.go
```

### Server won't start

**Problem**: Configuration or dependency issues

**Solution**:

```bash
# Verify config file exists
ls -la config.yaml

# Check logs for specific error
./build/pulumicost-mcp --config config.yaml 2>&1 | tee server.log

# Common issues:
# - Port 8080 already in use
# - pulumicost binary not found
# - Invalid YAML syntax
```

## Next Steps

After completing this quickstart:

1. **Implement User Story 2** (Plugin Management): Follow same TDD workflow for
   plugin discovery and validation
2. **Implement User Story 3** (Cost Optimization): Add recommendations, anomaly
   detection, forecasting
3. **Add Integration Tests**: Test with real pulumicost-core and plugins
4. **Performance Testing**: Ensure <3s P95 latency target
5. **Documentation**: Update README with examples and usage

## Resources

- [Goa Documentation](https://goa.design)
- [Goa-AI Guide](https://goa.design/goa-ai)
- [MCP Specification](https://modelcontextprotocol.io)
- [Project IMPLEMENTATION_PLAN.md](../../IMPLEMENTATION_PLAN.md)
- [Project CLAUDE.md](../../CLAUDE.md)
- [Constitution](../../.specify/memory/constitution.md)

## Getting Help

- **GitHub Issues**: Report bugs or request features
- **Discussions**: Ask questions about implementation
- **CLAUDE.md**: Project-specific patterns and conventions
