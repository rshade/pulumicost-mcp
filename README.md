# PulumiCost MCP Server

**AI-Powered Cloud Cost Analysis via Model Context Protocol** - A production-grade MCP server built with Goa and Goa-AI that brings PulumiCost capabilities to AI assistants and agents.

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

## Overview

PulumiCost MCP Server is a comprehensive Model Context Protocol (MCP) implementation that exposes PulumiCost's cloud cost analysis capabilities to AI assistants like Claude, ChatGPT, and custom AI agents. Built using Goa-AI for type-safe, drift-free integration, it enables natural language interaction with infrastructure cost data.

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

```
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

# Install Goa tooling
go install goa.design/goa/v3/cmd/goa@latest
go install goa.design/goa-ai/cmd/goa-ai@latest
```

### Installation

```bash
# Clone the repository
git clone https://github.com/rshade/pulumicost-mcp
cd pulumicost-mcp

# Install dependencies
go mod download

# Generate Goa code from design
goa gen github.com/rshade/pulumicost-mcp/design

# Build the server
make build

# Run the MCP server
./bin/pulumicost-mcp
```

### Integration with Claude Desktop

Add to your Claude Desktop MCP configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "pulumicost": {
      "command": "/path/to/pulumicost-mcp",
      "args": ["--port", "8080"],
      "env": {
        "PULUMI_ACCESS_TOKEN": "your-token"
      }
    }
  }
}
```

### Example Usage

Once configured, you can interact with PulumiCost via natural language in Claude:

```
User: What are the projected monthly costs for my staging environment?

Claude: [Uses analyze_projected_costs tool]
Based on your Pulumi stack, here are the projected costs:
- AWS EC2: $234.50/month
- AWS RDS: $156.00/month
- AWS S3: $12.30/month
Total: $402.80/month

User: How does that compare to last month's actual costs?

Claude: [Uses get_actual_costs tool]
Last month's actual costs were:
- AWS EC2: $289.45 (23% higher than projected)
- AWS RDS: $156.00 (on target)
- AWS S3: $18.92 (54% higher than projected)
Total: $464.37 (15% over projection)

Recommendation: Your EC2 instances ran more hours than projected.
Consider implementing auto-scaling or scheduled shutdowns for dev environments.
```

## Core Features

### 1. Cost Query Tools

#### Projected Costs
```
Tool: analyze_projected_costs
Description: Calculate estimated monthly costs before deploying infrastructure
Parameters:
  - pulumi_stack: Stack name or path to Pulumi JSON
  - filters: Resource filters (type, provider, tags)
  - output_format: table, json, summary
```

#### Actual Costs
```
Tool: get_actual_costs
Description: Retrieve historical spending with detailed breakdowns
Parameters:
  - pulumi_stack: Stack name
  - date_range: Start and end dates
  - group_by: provider, service, resource, tag
  - granularity: hourly, daily, monthly
```

#### Cost Comparison
```
Tool: compare_costs
Description: Compare projected vs actual, or costs across time periods
Parameters:
  - stack_a: First stack/period
  - stack_b: Second stack/period
  - comparison_type: projected_vs_actual, period_over_period
```

### 2. Plugin Development Tools

#### Validate Plugin Spec
```
Tool: validate_plugin_spec
Description: Validate pulumicost-spec compliance for plugin development
Parameters:
  - plugin_path: Path to plugin source
  - spec_version: Target spec version (default: latest)
  - conformance_level: basic, standard, advanced
```

#### Generate Plugin Scaffold
```
Tool: generate_plugin
Description: Create new cost source plugin from template
Parameters:
  - plugin_name: Name of the plugin
  - provider: Cloud provider (aws, azure, gcp, kubernetes, custom)
  - billing_models: Supported billing models
```

#### Test Plugin
```
Tool: test_plugin
Description: Run conformance tests against plugin implementation
Parameters:
  - plugin_path: Path to plugin
  - test_level: basic, standard, advanced
```

### 3. Resource Analysis Tools

#### Analyze Resource Costs
```
Tool: analyze_resource
Description: Detailed cost breakdown for specific resources
Parameters:
  - resource_urn: Pulumi resource URN
  - include_dependencies: Include related resource costs
  - time_period: Analysis time range
```

#### Optimize Recommendations
```
Tool: get_optimization_recommendations
Description: AI-powered cost optimization suggestions
Parameters:
  - stack_name: Pulumi stack
  - optimization_goals: reduce_cost, improve_performance, both
  - aggressiveness: conservative, moderate, aggressive
```

## Project Structure

```
pulumicost-mcp/
├── design/                    # Goa design files (DSL)
│   ├── design.go             # Main API design
│   ├── cost_tools.go         # Cost query tool definitions
│   ├── plugin_tools.go       # Plugin development tool definitions
│   └── types.go              # Shared type definitions
├── cmd/
│   └── pulumicost-mcp/       # Main server entry point
│       └── main.go
├── internal/
│   ├── service/              # Business logic implementation
│   │   ├── cost_service.go
│   │   ├── plugin_service.go
│   │   └── analysis_service.go
│   ├── adapter/              # Integration adapters
│   │   ├── pulumicost_adapter.go  # pulumicost-core integration
│   │   ├── spec_adapter.go        # pulumicost-spec integration
│   │   └── plugin_adapter.go      # Plugin management
│   └── plugin/               # Plugin discovery and loading
│       ├── manager.go
│       └── validator.go
├── pkg/
│   ├── types/                # Shared types
│   └── client/               # MCP client utilities
├── gen/                      # Generated Goa code (gitignored)
│   ├── cost/                 # Generated service code
│   ├── http/                 # Generated HTTP transport
│   └── mcp/                  # Generated MCP protocol
├── docs/
│   ├── architecture/         # Architecture documentation
│   │   ├── system-overview.md
│   │   ├── mcp-integration.md
│   │   └── plugin-system.md
│   ├── guides/              # User guides
│   │   ├── user-guide.md
│   │   ├── developer-guide.md
│   │   └── plugin-development.md
│   └── api/                 # API documentation
├── role-prompts/            # AI assistant role prompts
│   ├── senior-architect.md
│   ├── product-manager.md
│   ├── devops-engineer.md
│   ├── plugin-developer.md
│   └── cost-analyst.md
├── examples/
│   ├── queries/             # Example cost queries
│   └── plugins/             # Example plugin implementations
├── scripts/                 # Build and deployment scripts
│   ├── generate.sh         # Run Goa generation
│   ├── build.sh            # Build binaries
│   └── install-tools.sh    # Install development tools
├── test/
│   ├── integration/        # Integration tests
│   └── unit/              # Unit tests
├── .claude/               # Claude Code configuration
├── CLAUDE.md             # AI development context
├── Makefile
├── go.mod
├── go.sum
└── README.md
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
make generate      # Generate Goa code from design
make build         # Build server binary
make test          # Run all tests
make lint          # Run linters
make install       # Install server to $GOPATH/bin
make clean         # Clean generated files
make docs          # Generate documentation
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

This project includes specialized prompt files for different roles in `role-prompts/`:

- **Senior Architect**: System design, architecture decisions, scalability planning
- **Product Manager**: Feature prioritization, roadmap planning, user stories
- **DevOps Engineer**: Deployment, monitoring, operational excellence
- **Plugin Developer**: Plugin creation, spec compliance, testing
- **Cost Analyst**: Cost optimization, reporting, budget management

See [role-prompts/README.md](role-prompts/README.md) for usage instructions.

## Documentation

- **[Architecture Overview](docs/architecture/system-overview.md)**: System design and components
- **[User Guide](docs/guides/user-guide.md)**: Getting started and common workflows
- **[Developer Guide](docs/guides/developer-guide.md)**: Development setup and contribution guidelines
- **[Plugin Development](docs/guides/plugin-development.md)**: Building cost source plugins
- **[API Reference](docs/api/)**: Complete API documentation

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup

```bash
# Clone and setup
git clone https://github.com/rshade/pulumicost-mcp
cd pulumicost-mcp
make setup

# Install development tools
make install-tools

# Run tests
make test

# Run linters
make lint
```

## Roadmap

### Phase 1: Core MCP Server (Current)
- ✅ Goa-AI integration
- ✅ Basic cost query tools
- ✅ Plugin validation tools
- ✅ MCP protocol implementation
- 🔄 Claude Desktop integration

### Phase 2: Enhanced Analysis
- ⏳ Advanced cost optimization recommendations
- ⏳ Anomaly detection
- ⏳ Budget tracking and alerts
- ⏳ Cost forecasting with ML

### Phase 3: Developer Experience
- ⏳ Interactive plugin scaffolding
- ⏳ Real-time cost feedback in IDE
- ⏳ CI/CD cost gates
- ⏳ Visual cost dashboards

### Phase 4: Enterprise Features
- ⏳ Multi-tenant support
- ⏳ RBAC and audit logging
- ⏳ Custom plugin marketplace
- ⏳ Advanced reporting

## License

Apache-2.0 - See [LICENSE](LICENSE) for details.

## Related Projects

- [pulumicost-core](https://github.com/rshade/pulumicost-core) - Cost analysis orchestration
- [pulumicost-spec](https://github.com/rshade/pulumicost-spec) - Plugin specification protocol
- [Goa](https://goa.design/) - Design-first API framework
- [Goa-AI](https://goa.design/goa-ai) - AI extensions for Goa
- [MCP](https://modelcontextprotocol.io/) - Model Context Protocol

## Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/rshade/pulumicost-mcp/issues)
- **Discussions**: [GitHub Discussions](https://github.com/rshade/pulumicost-mcp/discussions)
- **Community**: [Discord](https://discord.gg/pulumicost)

---

**Built with ❤️ using Goa and Goa-AI** - Making cloud cost analysis accessible to AI assistants everywhere.
