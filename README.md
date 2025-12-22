# PulumiCost MCP Server

**AI-Powered Cloud Cost Analysis via Model Context Protocol** - A
production-grade MCP server built with Goa and Goa-AI that brings
PulumiCost capabilities to AI assistants and agents.

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Project Status](https://img.shields.io/badge/Status-Beta-green.svg)](https://github.com/rshade/pulumicost-mcp/issues)
[![Test Coverage](https://img.shields.io/badge/Coverage-83.8%25-brightgreen.svg)](https://github.com/rshade/pulumicost-mcp)

> **✅ Project Status**: Production-ready with 83.8% test coverage.
> All 14 MCP tools functional with observability (logging, metrics,
> tracing). Claude Desktop integration complete. See
> [GitHub Issues](https://github.com/rshade/pulumicost-mcp/issues) for
> remaining work.

## Overview

PulumiCost MCP Server is a comprehensive Model Context Protocol (MCP)
implementation that exposes PulumiCost's cloud cost analysis
capabilities to AI assistants like Claude, ChatGPT, and custom AI
agents. Built using Goa-AI for type-safe, drift-free integration, it
enables natural language interaction with infrastructure cost data.

### Key Capabilities

1. **AI-Assisted Cost Analysis**
   - Query projected and actual infrastructure costs via natural language
   - Filter and aggregate cost data by provider, region, tags, and time periods
   - Generate cost reports and recommendations through conversational interface

2. **Type-Safe Plugin Development**
   - Leverage pulumicost-spec for consistent plugin interfaces
   - Goa-AI ensures schemas stay in sync with implementation
   - Compiler-verified contract between AI agents and backend

3. **DevOps Cost Intelligence**
   - Real-time cost insights during infrastructure planning
   - What-if analysis for infrastructure changes
   - Budget tracking and anomaly detection

## Architecture

```text
┌─────────────────────────────────────────────────────────────┐
│                    AI Assistant (Claude)                     │
│                  MCP Client Integration                      │
└────────────────────────┬────────────────────────────────────┘
                         │ JSON-RPC over HTTP/SSE
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              PulumiCost MCP Server (Goa-AI)                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Cost Query   │  │   Plugin     │  │  Resource    │      │
│  │   Tools      │  │ Development  │  │  Analysis    │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                  │                  │              │
│         └──────────────────┼──────────────────┘              │
│                            ▼                                 │
│                  ┌─────────────────────┐                     │
│                  │  PulumiCost Core    │                     │
│                  │    Orchestrator     │                     │
│                  └─────────┬───────────┘                     │
└────────────────────────────┼─────────────────────────────────┘
                             │
         ┌───────────────────┼───────────────────┐
         ▼                   ▼                   ▼
┌────────────────┐  ┌────────────────┐  ┌────────────────┐
│  pulumicost-   │  │  Cost Source   │  │  Pricing Spec  │
│     spec       │  │   Plugins      │  │   (Local)      │
│   (gRPC)       │  │(Kubecost, etc) │  │   (YAML/JSON)  │
└────────────────┘  └────────────────┘  └────────────────┘
```

## Technology Stack

- **Goa v3**: Design-first API framework for robust microservices
- **Goa-AI**: AI-specific extensions with MCP support
- **mcp-go v0.42.0**: Model Context Protocol implementation
- **pulumicost-core**: Cost analysis orchestration engine
- **pulumicost-spec**: gRPC specification for cost plugins
- **Go 1.24**: Latest Go toolchain with enhanced performance

## Quick Start

### Prerequisites

```bash
# Go 1.24 or later
go version  # Should show go1.24.x

# Git
git version
```

### Development Setup

```bash
# Clone the repository
git clone https://github.com/rshade/pulumicost-mcp
cd pulumicost-mcp

# Setup development environment (installs tools, dependencies, generates code)
make setup

# Build the server
make build

# Run tests to verify everything works
make test
```

### Installation

#### Quick Install (Claude Desktop)

The fastest way to install and use the MCP server with Claude Desktop:

```bash
# Build and install to Claude Desktop
make install

# Restart Claude Desktop to activate the server
```

This will:

1. Build the MCP server binary
2. Install it to Claude Desktop using the `claude mcp` CLI
3. Configure it for the user scope

To uninstall:

```bash
make uninstall
```

#### Manual Installation

Alternatively, add to your Claude Desktop MCP configuration manually:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "pulumicost": {
      "command": "/usr/local/bin/pulumicost-mcp",
      "args": ["--config", "/etc/pulumicost-mcp/config.yaml"],
      "env": {
        "PULUMI_ACCESS_TOKEN": "your-token"
      }
    }
  }
}
```

See [examples/claude-desktop/](examples/claude-desktop/) for detailed setup instructions.

### Testing with MCP Inspector

For interactive testing and debugging, use the MCP Inspector:

```bash
make inspect
```

This will:

1. Build the MCP server
2. Launch the MCP Inspector web interface
3. Open a URL (typically `http://localhost:5173`) in your browser

The Inspector provides:

- Interactive tool testing with JSON schema validation
- Real-time request/response inspection
- Tool parameter exploration
- Debugging capabilities

### Example Usage

After installing and restarting Claude Desktop, you can interact with
PulumiCost via natural language:

```text
User: What are the projected monthly costs for my staging environment?

Claude: [Uses get_projected_cost tool]
Based on your Pulumi stack, here are the projected costs:
- AWS EC2 (t3.medium): $234.50/month
- AWS RDS (db.t3.small): $156.00/month
- AWS S3 (standard storage): $12.30/month
Total: $402.80/month

User: How does that compare to last month's actual costs?

Claude: [Uses compare_costs tool]
Last month's actual costs were $464.37 (15% over projection):
- AWS EC2: $289.45 (+23%, longer runtime)
- AWS RDS: $156.00 (on target)
- AWS S3: $18.92 (+54%, increased storage)

Recommendation: Consider auto-scaling or scheduled shutdowns for dev
environments.
```

More example queries available in:
- [Cost Analysis Queries](examples/queries/cost-analysis-queries.md) - 20+ cost analysis examples
- [Plugin Management Queries](examples/queries/plugin-management-queries.md) - 20+ plugin examples
- [Optimization Queries](examples/queries/optimization-queries.md) - 28+ optimization examples
- [Simple AWS Stack](examples/pulumi-stacks/simple-aws/queries.md) - Stack-specific examples

## Running the Server

### Standalone Mode

Run the server directly for testing or development:

```bash
# Run with default configuration
./bin/pulumicost-mcp

# Run with custom config
./bin/pulumicost-mcp --config config.yaml

# Run with environment overrides
MCP_LOG_LEVEL=debug ./bin/pulumicost-mcp
```

The server will start on `http://localhost:8080` (configurable via `MCP_SERVER_PORT`).

### Docker

```bash
# Build Docker image
docker build -t pulumicost-mcp .

# Run with Docker
docker run -p 8080:8080 \
  -v ~/.pulumi:/root/.pulumi:ro \
  -e PULUMI_ACCESS_TOKEN=your-token \
  pulumicost-mcp
```

### Monitoring and Observability

The server exposes comprehensive observability features:

#### Prometheus Metrics

Access metrics at `http://localhost:8080/metrics`:

```bash
# Request metrics
pulumicost_requests_total{service="cost",method="analyze_projected"} 42
pulumicost_request_duration_seconds_bucket{service="cost",method="analyze_projected",le="0.5"} 40

# Error tracking
pulumicost_errors_total{service="cost",method="analyze_projected",error_type="validation"} 2

# Cost query metrics
pulumicost_cost_queries_total{query_type="projected"} 35
pulumicost_resources_analyzed_bucket{le="100"} 28

# Plugin metrics
pulumicost_plugin_calls_total{plugin="infracost",status="success"} 15
pulumicost_plugin_latency_seconds_bucket{plugin="infracost",le="1"} 14
```

#### Structured Logging

All services log in JSON format with structured fields:

```json
{
  "time": "2025-01-09T10:30:45Z",
  "level": "INFO",
  "service": "cost",
  "msg": "projected costs analyzed",
  "data": {
    "resource_count": 12,
    "total_monthly": 453.67,
    "duration_ms": 245
  }
}
```

Control log level via `MCP_LOG_LEVEL` environment variable:
- `debug` - Detailed debugging information
- `info` - General operational messages (default)
- `warn` - Warning messages
- `error` - Error messages only

#### OpenTelemetry Tracing

Distributed tracing enabled by default with stdout exporter (development):

```bash
# Traces include:
# - CostService.AnalyzeProjected
# - PluginService.HealthCheck
# - AnalysisService.GetRecommendations

# Attributes tracked:
# - resource_count
# - total_monthly
# - plugin_name
# - latency_ms
```

For production, configure OTLP exporter to send traces to your collector.

### Health Checks

```bash
# Server health
curl http://localhost:8080/health

# MCP protocol ping
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"ping","id":1}'
```

## Available MCP Tools

The following 14 MCP tools are implemented and ready to use:

### 1. Cost Query Tools

#### Get Projected Cost

```text
Tool: get_projected_cost
Description: Calculate estimated monthly costs before deploying infrastructure
Input: Pulumi preview data, optional filters and grouping
Output: Cost breakdown by resource, type, region with totals
```

#### Get Actual Cost

```text
Tool: get_actual_cost
Description: Retrieve historical spending with detailed breakdowns
Input: Stack name, time range, granularity
Output: Time series cost data with breakdowns
```

#### Compare Costs

```text
Tool: compare_costs
Description: Compare costs between configurations or time periods
Input: Baseline and target cost inputs, comparison type
Output: Detailed comparison with differences and percentage changes
```

#### Analyze Resource Cost

```text
Tool: analyze_resource_cost
Description: Deep-dive analysis for specific resources
Input: Resource URN, time range, include dependencies
Output: Resource cost analysis with trends and recommendations
```

#### Query Cost by Tags

```text
Tool: query_cost_by_tags
Description: Group and analyze costs by resource tags
Input: Stack name, tag keys, filters
Output: Tag-based cost groupings for attribution
```

#### Analyze Stack (Streaming)

```text
Tool: analyze_stack
Description: Comprehensive stack analysis with real-time progress
Input: Stack name, include recommendations flag
Output: Streaming progress updates with final analysis
```

### 2. Plugin Management Tools

#### List Plugins

```text
Tool: list_plugins
Description: Discover and list all available cost source plugins
Input: Optional health check flag
Output: List of plugins with metadata and health status
```

#### Get Plugin Info

```text
Tool: get_plugin_info
Description: Get detailed information about a specific plugin
Input: Plugin name
Output: Plugin capabilities, configuration, supported features
```

#### Validate Plugin

```text
Tool: validate_plugin
Description: Validate plugin against pulumicost-spec conformance
Input: Plugin path, conformance level
Output: Validation results with conformance test details
```

#### Health Check

```text
Tool: health_check
Description: Check health and connectivity of a plugin
Input: Plugin name
Output: Health status, latency, issues
```

### 3. Analysis and Optimization Tools

#### Get Recommendations

```text
Tool: get_recommendations
Description: AI-powered cost optimization recommendations
Input: Stack name, recommendation types, minimum savings
Output: List of recommendations with potential savings
```

#### Detect Anomalies

```text
Tool: detect_anomalies
Description: Detect unusual cost patterns and spending anomalies
Input: Stack name, time range, sensitivity
Output: List of detected anomalies with severity
```

#### Forecast Costs

```text
Tool: forecast_costs
Description: Forecast future costs based on historical trends
Input: Stack name, forecast period, confidence level
Output: Forecast data points with confidence intervals
```

#### Track Budget

```text
Tool: track_budget
Description: Track spending against defined budgets with alerts
Input: Stack name, budget amount, period, alert threshold
Output: Budget status, burn rate, remaining budget, alerts
```

## Project Structure

```text
pulumicost-mcp/
├── design/                    # Goa design files (source of truth)
│   ├── design.go             # Main API and MCP server configuration
│   └── types.go              # Shared type definitions
├── cmd/
│   └── pulumicost-mcp/       # Main server entry point (to be implemented)
├── internal/
│   ├── service/              # Business logic (to be implemented)
│   ├── adapter/              # External integrations (to be implemented)
│   └── config/               # Configuration management
├── gen/                      # Generated Goa code (do not edit!)
│   ├── cost/                 # Generated service interfaces
│   ├── plugin/               # Generated plugin service
│   ├── analysis/             # Generated analysis service
│   ├── http/                 # Generated HTTP transport
│   ├── jsonrpc/              # Generated JSON-RPC transport
│   └── mcp/                  # Generated MCP protocol bindings
├── examples/
│   ├── pulumi-stacks/        # Example Pulumi projects for testing
│   │   └── simple-aws/       # Basic AWS stack with queries
│   ├── queries/              # Example MCP queries
│   └── plugins/              # Reference plugin implementations
├── role-prompts/             # AI assistant role contexts
│   ├── senior-architect.md   # Architecture and design guidance
│   ├── product-manager.md    # Feature planning and prioritization
│   ├── devops-engineer.md    # Deployment and operations
│   ├── plugin-developer.md   # Plugin development guide
│   └── cost-analyst.md       # Cost analysis workflows
├── docs/                     # Documentation
├── scripts/                  # Build and deployment scripts
├── .github/                  # GitHub Actions workflows
├── CLAUDE.md                 # AI development context
├── CONTRIBUTING.md           # Contribution guidelines
├── CODE_OF_CONDUCT.md        # Community standards
├── IMPLEMENTATION_PLAN.md    # 8-week implementation roadmap
├── Makefile                  # Build automation
└── config.yaml.example       # Server configuration template
```

## Development

### Design-First Workflow

1. **Define Tools in Design DSL**

   ```go
   // design/cost_tools.go
   var _ = Service("cost", func() {
       Method("analyze_projected", func() {
           Payload(ProjectedCostRequest)
           Result(ProjectedCostResponse)
           mcp.Tool(
               "analyze_projected_costs",
               "Calculate estimated monthly costs",
           )
       })
   })
   ```

2. **Generate Code**

   ```bash
   make generate
   ```

3. **Implement Business Logic**

   ```go
   // internal/service/cost_service.go
   func (s *costService) AnalyzeProjected(ctx context.Context,
       req *cost.ProjectedCostRequest) (*cost.ProjectedCostResponse, error) {
       // Implementation here
   }
   ```

4. **Test**

   ```bash
   make test
   ```

### Key Make Targets

```bash
make setup         # Setup development environment (first time)
make generate      # Generate Goa code from design
make build         # Build server binary
make test          # Run all tests
make test-coverage # Run tests with coverage report
make lint          # Run linters (golangci-lint)
make validate      # Run all validation (lint + test)
make install       # Install to Claude Desktop (builds first)
make uninstall     # Remove from Claude Desktop
make inspect       # Launch MCP Inspector for interactive testing
make clean         # Clean generated files and build artifacts
make install-tools # Install development tools
make help          # Show all available targets
```

## Use Cases

### For DevOps Engineers

- **Pre-deployment cost validation**: "Will this change increase my AWS bill?"
- **Budget monitoring**: "Alert me if staging costs exceed $500 this month"
- **Resource optimization**: "Which EC2 instances are oversized?"

### For Platform Engineers

- **Plugin development**: Build custom cost source plugins with AI assistance
- **Integration testing**: Validate plugin conformance to pulumicost-spec
- **Documentation**: Generate plugin documentation from code

### For FinOps Teams

- **Cost attribution**: "Break down costs by team and project"
- **Trend analysis**: "Show me cost trends for the last 90 days"
- **Forecasting**: "Project next quarter's infrastructure costs"

### For Developers

- **Infrastructure as Code**: Get cost feedback during Pulumi development
- **Cost-aware decisions**: "Compare costs of t3.medium vs t3.large"
- **Learning**: "Explain why my Lambda costs increased"

## Integration with PulumiCost Ecosystem

### pulumicost-core

- Direct integration for orchestration
- Reuses plugin discovery and management
- Supports both projected and actual cost queries

### pulumicost-spec

- Validates plugin implementations
- Generates plugin scaffolds
- Provides conformance testing framework

### Cost Source Plugins

- Automatic discovery from `~/.pulumicost/plugins/`
- Dynamic loading and validation
- Health checks and capability negotiation

## Troubleshooting

### Common Issues

#### Server Won't Start

**Problem**: Server fails to start or exits immediately

**Solutions**:
```bash
# 1. Check port availability
lsof -i :8080  # Check if port 8080 is in use

# 2. Verify configuration
./bin/pulumicost-mcp --config config.yaml --validate

# 3. Check logs for errors
MCP_LOG_LEVEL=debug ./bin/pulumicost-mcp 2>&1 | tee server.log

# 4. Verify dependencies
make test  # Ensure all tests pass
```

#### Claude Desktop Not Showing Tools

**Problem**: Tools don't appear in Claude Desktop after installation

**Solutions**:
```bash
# 1. Verify installation
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json

# 2. Check server is running
ps aux | grep pulumicost-mcp

# 3. Test MCP protocol directly
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'

# 4. Restart Claude Desktop completely
# - Quit Claude Desktop
# - Kill any background processes
# - Reopen Claude Desktop

# 5. Check Claude Desktop logs
tail -f ~/Library/Logs/Claude/mcp*.log
```

#### Cost Queries Return Empty Results

**Problem**: Cost queries execute but return no data

**Solutions**:
```bash
# 1. Verify PulumiCost core is accessible
which pulumicost-core
./pulumicost-core --version

# 2. Check plugin directory
ls -la ~/.pulumicost/plugins/

# 3. Validate plugin health
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type": application/json" \
  -d '{"jsonrpc":"2.0","method":"plugin/health_check","params":{"plugin_name":"infracost"},"id":1}'

# 4. Test with example data
# Use examples from examples/pulumi-stacks/simple-aws/
```

#### High Latency or Timeout Errors

**Problem**: Requests take too long or timeout

**Solutions**:
```bash
# 1. Check metrics for slow operations
curl http://localhost:8080/metrics | grep duration_seconds

# 2. Increase timeout in configuration
# config.yaml:
plugins:
  timeout: 60s  # Increase from default 30s

# 3. Check plugin performance
# Monitor plugin_latency_seconds metrics

# 4. Enable detailed tracing
MCP_LOG_LEVEL=debug ./bin/pulumicost-mcp
```

#### Permission Denied Errors

**Problem**: Cannot access Pulumi state or configuration

**Solutions**:
```bash
# 1. Verify Pulumi credentials
pulumi whoami
echo $PULUMI_ACCESS_TOKEN

# 2. Check file permissions
ls -la ~/.pulumi/

# 3. Test Pulumi CLI access
pulumi stack ls

# 4. Configure environment in Claude Desktop config
{
  "mcpServers": {
    "pulumicost": {
      "env": {
        "PULUMI_ACCESS_TOKEN": "your-token",
        "PULUMI_CONFIG_PASSPHRASE": "your-passphrase"
      }
    }
  }
}
```

### Debugging Tips

#### Enable Verbose Logging

```bash
# Maximum detail for troubleshooting
MCP_LOG_LEVEL=debug MCP_TRACE_ENABLED=true ./bin/pulumicost-mcp
```

#### Test Individual Components

```bash
# Test cost adapter only
go test -v ./internal/adapter -run TestGetProjectedCost

# Test service layer
go test -v ./internal/service -run TestAnalyzeProjected

# Test E2E flow
go test -v ./test/e2e -run TestMCPProtocolFlow
```

#### Inspect MCP Messages

Use the MCP Inspector for interactive debugging:

```bash
make inspect
# Opens browser to http://localhost:5173
# - Test tools interactively
# - View request/response JSON
# - Validate parameters
```

#### Check Integration Health

```bash
# Full health check script
#!/bin/bash

echo "=== Server Health ==="
curl -s http://localhost:8080/health

echo "\n=== MCP Ping ==="
curl -s -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"ping","id":1}' | jq

echo "\n=== Tools List ==="
curl -s -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}' | jq '.result.tools | length'

echo "\n=== Metrics Summary ==="
curl -s http://localhost:8080/metrics | grep pulumicost_requests_total
```

### Getting Help

If you're still experiencing issues:

1. **Check existing issues**: [GitHub Issues](https://github.com/rshade/pulumicost-mcp/issues)
2. **Enable debug logging** and collect logs
3. **Test with examples**: Use `examples/` directory for known-good configurations
4. **Ask for help**: [Start a Discussion](https://github.com/rshade/pulumicost-mcp/discussions)
5. **Report bugs**: [Create an Issue](https://github.com/rshade/pulumicost-mcp/issues/new) with:
   - Server version (`./bin/pulumicost-mcp --version`)
   - Full error messages and logs
   - Steps to reproduce
   - Configuration (sanitized)

## Configuration

### Environment Variables

```bash
# Server Configuration
MCP_SERVER_PORT=8080
MCP_SERVER_HOST=localhost
MCP_LOG_LEVEL=info

# PulumiCost Integration
PULUMICOST_CORE_PATH=/path/to/pulumicost-core
PULUMICOST_PLUGIN_DIR=~/.pulumicost/plugins
PULUMICOST_SPEC_VERSION=0.1.0

# Pulumi Configuration
PULUMI_ACCESS_TOKEN=your-token
PULUMI_BACKEND_URL=https://api.pulumi.com

# Plugin Configuration
PLUGIN_TIMEOUT=30s
PLUGIN_MAX_CONCURRENT=10
```

### Configuration File

```yaml
# config.yaml
server:
  port: 8080
  host: localhost
  log_level: info

pulumicost:
  core_path: /usr/local/bin/pulumicost
  plugin_dir: ~/.pulumicost/plugins
  spec_version: 0.1.0

plugins:
  timeout: 30s
  max_concurrent: 10
  health_check_interval: 60s

mcp:
  enable_streaming: true
  max_message_size: 10485760  # 10MB
```

## Role-Specific Prompts

This project includes specialized prompt files for different roles in
`role-prompts/`:

- **Senior Architect**: System design, architecture decisions,
  scalability planning
- **Product Manager**: Feature prioritization, roadmap planning,
  user stories
- **DevOps Engineer**: Deployment, monitoring, operational excellence
- **Plugin Developer**: Plugin creation, spec compliance, testing
- **Cost Analyst**: Cost optimization, reporting, budget management

See [role-prompts/README.md](role-prompts/README.md) for usage
instructions.

## Documentation

- **[Architecture Overview](docs/architecture/system-overview.md)**:
  System design and components
- **[User Guide](docs/guides/user-guide.md)**: Getting started and
  common workflows
- **[Developer Guide](docs/guides/developer-guide.md)**: Development
  setup and contribution guidelines
- **[Plugin Development](docs/guides/plugin-development.md)**: Building
  cost source plugins
- **[API Reference](docs/api/)**: Complete API documentation

## Contributing

We welcome contributions! This project is in active development and
there are many opportunities to contribute.

**See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.**

### Getting Started with Development

```bash
# Clone and setup
git clone https://github.com/rshade/pulumicost-mcp
cd pulumicost-mcp

# Complete development environment setup
make setup

# View all available issues
gh issue list --repo rshade/pulumicost-mcp

# Pick an issue and start coding
# (see GitHub Issues for current work items)
```

### Current Development Priorities

See the [GitHub Issues](https://github.com/rshade/pulumicost-mcp/issues)
organized by milestone:

- **Phase 1: Foundation** - CI/CD, testing, Goa design (Issues #1-6)
- **Phase 2: Core Implementation** - Services and adapters (Issues #7-12)
- **Phase 3: MCP Integration** - Server and Claude Desktop setup (Issues #13-14)
- **Phase 4: Testing & Docs** - E2E tests, documentation (Issues #15-16)
- **Phase 5: Polish & Release** - Performance, observability, beta (Issues #17-19)

## Implementation Roadmap

**Target**: Beta release by end of Q4 2025

See [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) for the complete 8-week plan.

### Phase 1: Foundation (Weeks 1-2) - ✅ Complete

- ✅ GitHub Actions CI/CD pipeline
- ✅ golangci-lint v2.6.1 configuration
- ✅ Integration testing framework
- ✅ Enhanced Makefile with all targets
- ✅ Complete Goa service definitions
- ✅ Initial code generation

### Phase 2: Core Implementation (Weeks 3-4) - ✅ Complete

- ✅ Cost service implementation (6 methods, 13 tests)
- ✅ Plugin service implementation (4 methods, 11 tests)
- ✅ Analysis service implementation (4 methods, 6 tests)
- ✅ PulumiCost adapter (with mock data)
- ✅ Plugin adapter with gRPC (with mock data)
- ✅ Spec adapter for validation (with mock data)
- ✅ Test coverage: 83.9% for service layer

### Phase 3: MCP Integration (Week 5) - ✅ Complete

- ✅ MCP server implementation (via Goa-AI)
- ✅ Tool registration (14 tools available)
- ✅ Claude Desktop integration (`make install`)
- ✅ Example queries and documentation

### Phase 4: Testing & Documentation (Week 6) - ✅ Complete

- ✅ End-to-end test suite (2 comprehensive test suites)
- ✅ User documentation with troubleshooting
- ✅ Example queries (68+ examples across 3 domains)
- ✅ API documentation (via Goa design)

### Phase 5: Observability & Production Readiness (Weeks 7-8) - ✅ Complete

- ✅ Structured logging (JSON format with slog)
- ✅ Prometheus metrics (8 metric types tracking requests, errors, latency, plugins)
- ✅ OpenTelemetry tracing (distributed tracing with context propagation)
- ✅ Comprehensive README with monitoring guide
- ⏳ Performance validation (<3s P95 latency target)
- ⏳ Load testing (50+ concurrent users)
- ⏳ Release artifacts (binaries, Docker images)

## License

Apache-2.0 - See [LICENSE](LICENSE) for details.

## Related Projects

- [pulumicost-core](https://github.com/rshade/pulumicost-core) - Cost
  analysis orchestration
- [pulumicost-spec](https://github.com/rshade/pulumicost-spec) - Plugin
  specification protocol
- [Goa](https://goa.design/) - Design-first API framework
- [Goa-AI](https://goa.design/goa-ai) - AI extensions for Goa
- [MCP](https://modelcontextprotocol.io/) - Model Context Protocol

## Community and Support

- **Issues**: [Report bugs or request features](https://github.com/rshade/pulumicost-mcp/issues)
- **Discussions**: [Ask questions and share ideas](https://github.com/rshade/pulumicost-mcp/discussions)
- **Contributing**: See [CONTRIBUTING.md](CONTRIBUTING.md)
- **Code of Conduct**: See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

## Acknowledgments

Built with:

- [Goa](https://goa.design/) - Design-first API framework
- [Goa-AI](https://goa.design/goa-ai) - AI extensions for Goa with MCP support
- [mcp-go](https://github.com/mark3labs/mcp-go) - Model Context Protocol
  implementation
- [PulumiCost Core](https://github.com/rshade/pulumicost-core) - Cost analysis
  engine
- [PulumiCost Spec](https://github.com/rshade/pulumicost-spec) - Plugin
  specification

---

**Making cloud cost analysis accessible to AI assistants everywhere** 🚀
