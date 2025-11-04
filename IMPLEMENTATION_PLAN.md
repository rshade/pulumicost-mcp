# PulumiCost MCP Server - Implementation Plan

## Executive Summary

This document outlines the complete implementation plan for building a production-grade Model Context Protocol (MCP) server that exposes PulumiCost's cloud cost analysis capabilities to AI assistants. The project leverages Goa and Goa-AI for design-first, type-safe development.

## Project Overview

### Goals
1. Enable AI assistants to perform cloud cost analysis via natural language
2. Provide type-safe plugin development workflow for pulumicost-spec
3. Offer DevOps teams intelligent cost insights during infrastructure changes
4. Build extensible, scalable architecture using proven Go patterns

### Non-Goals
- Reimplementing pulumicost-core functionality
- Building a web UI (MCP server only)
- Replacing existing cost management platforms
- Supporting non-gRPC plugin protocols

## Architecture Decision Records

### ADR-001: Use Goa + Goa-AI

**Context**: Need design-first API framework with MCP support

**Decision**: Use Goa v3 with Goa-AI extensions

**Rationale**:
- Design-first eliminates schema drift
- Code generation ensures consistency
- Goa-AI provides native MCP support
- Strong typing throughout stack
- Battle-tested in production

**Consequences**:
- Learning curve for Goa DSL
- Generated code cannot be modified
- Tied to Goa's release cycle
- Excellent tooling and community

### ADR-002: gRPC for Plugin Communication

**Context**: Need standardized protocol for cost source plugins

**Decision**: Use pulumicost-spec gRPC protocol

**Rationale**:
- Already defined and proven
- Strong typing and versioning
- Efficient binary protocol
- Good Go support
- Industry standard

**Consequences**:
- Plugins must implement gRPC
- Network overhead for local plugins
- Excellent ecosystem and tools

### ADR-003: JSON-RPC for MCP Transport

**Context**: MCP specification requires JSON-RPC

**Decision**: Use Goa's JSON-RPC generator with SSE for streaming

**Rationale**:
- MCP standard protocol
- Goa has built-in support
- SSE for streaming responses
- HTTP-based (firewall-friendly)

**Consequences**:
- Larger payloads than gRPC
- Text-based (less efficient)
- Standard compliance

### ADR-004: Stateless Server Design

**Context**: Need scalable, reliable architecture

**Decision**: Stateless server with external plugins

**Rationale**:
- Horizontal scaling
- Simple deployment
- Plugin independence
- Cloud-native

**Consequences**:
- No session state
- Plugin discovery per request
- Caching layer needed

## Implementation Phases

### Phase 1: Foundation (Weeks 1-2)

#### Week 1: Project Setup & Core Structure

**Deliverables:**
- [x] Project repository structure
- [x] Go module configuration
- [x] Role-based prompts
- [x] README and documentation
- [ ] Initial Goa design files
- [ ] Build system (Makefile)
- [ ] CI/CD pipeline setup

**Tasks:**
1. Create directory structure
   - design/ for Goa DSL
   - cmd/ for entry points
   - internal/ for implementation
   - pkg/ for shared libraries
   - docs/ for documentation

2. Setup Go modules
   - Import Goa v3
   - Import Goa-AI
   - Import pulumicost dependencies
   - Import MCP-go

3. Create build system
   - Makefile targets
   - Code generation scripts
   - Testing scripts
   - Linting configuration

4. Setup CI/CD
   - GitHub Actions workflows
   - Code generation verification
   - Test execution
   - Linting and validation

**Acceptance Criteria:**
- `make generate` produces valid Go code
- `make test` runs successfully
- `make lint` passes
- CI pipeline runs on push

#### Week 2: Core Goa Design

**Deliverables:**
- [ ] Service definitions in Goa DSL
- [ ] Type definitions
- [ ] MCP tool annotations
- [ ] Generated code compiles

**Tasks:**
1. Define core services
   ```go
   // design/design.go
   var API = APIDesign("pulumicost-mcp", func() {
       Title("PulumiCost MCP Server")
       Description("AI-powered cloud cost analysis")
       Version("1.0.0")
   })
   ```

2. Cost Query Service
   ```go
   // design/cost_service.go
   var _ = Service("cost", func() {
       Description("Cost analysis and querying")

       Method("analyze_projected", func() {
           Payload(ProjectedCostRequest)
           Result(ProjectedCostResponse)
           mcp.Tool("analyze_projected_costs", "...")
       })

       Method("get_actual", func() {
           Payload(ActualCostRequest)
           Result(ActualCostResponse)
           mcp.Tool("get_actual_costs", "...")
       })

       Method("compare_costs", func() {
           Payload(CostComparisonRequest)
           Result(CostComparisonResponse)
           mcp.Tool("compare_costs", "...")
       })
   })
   ```

3. Plugin Management Service
   ```go
   // design/plugin_service.go
   var _ = Service("plugin", func() {
       Description("Plugin development and management")

       Method("validate_plugin", func() {
           Payload(ValidatePluginRequest)
           Result(ValidatePluginResponse)
           mcp.Tool("validate_plugin_spec", "...")
       })

       Method("generate_plugin", func() {
           Payload(GeneratePluginRequest)
           Result(GeneratePluginResponse)
           mcp.Tool("generate_plugin", "...")
       })

       Method("test_plugin", func() {
           Payload(TestPluginRequest)
           StreamingResult(TestPluginProgress)
           mcp.Tool("test_plugin", "...")
       })
   })
   ```

4. Resource Analysis Service
   ```go
   // design/analysis_service.go
   var _ = Service("analysis", func() {
       Description("Resource-level cost analysis")

       Method("analyze_resource", func() {
           Payload(ResourceAnalysisRequest)
           Result(ResourceAnalysisResponse)
           mcp.Tool("analyze_resource_cost", "...")
       })

       Method("get_recommendations", func() {
           Payload(RecommendationRequest)
           StreamingResult(RecommendationProgress)
           mcp.Tool("get_optimization_recommendations", "...")
       })
   })
   ```

5. Type Definitions
   ```go
   // design/types.go
   var ProjectedCostRequest = Type("ProjectedCostRequest", func() {
       Attribute("stack_name", String, "Pulumi stack name")
       Attribute("pulumi_json", String, "Pulumi preview JSON")
       Attribute("filters", Filters, "Resource filters")
       Required("pulumi_json")
   })

   var ProjectedCostResponse = Type("ProjectedCostResponse", func() {
       Attribute("total_monthly", Float64, "Total monthly cost")
       Attribute("currency", String, "Currency code")
       Attribute("resources", ArrayOf(ResourceCost), "Per-resource breakdown")
       Attribute("by_provider", MapOf(String, Float64), "Cost by provider")
       Attribute("by_service", MapOf(String, Float64), "Cost by service")
       Required("total_monthly", "currency", "resources")
   })
   ```

**Acceptance Criteria:**
- All services defined in DSL
- Types are comprehensive
- `goa gen` produces valid code
- API documentation generated
- JSON schemas exported

### Phase 2: Core Implementation (Weeks 3-4)

#### Week 3: Service Layer Implementation

**Deliverables:**
- [ ] Cost service implementation
- [ ] Plugin service implementation
- [ ] Analysis service implementation
- [ ] Unit tests for services

**Tasks:**
1. Implement Cost Service
   ```go
   // internal/service/cost_service.go
   type costService struct {
       adapter *adapter.PulumiCostAdapter
       logger  *log.Logger
   }

   func (s *costService) AnalyzeProjected(ctx context.Context,
       req *cost.ProjectedCostRequest) (*cost.ProjectedCostResponse, error) {

       // Parse Pulumi JSON
       // Call pulumicost-core
       // Transform response
       // Return results
   }
   ```

2. Implement Plugin Service
   ```go
   // internal/service/plugin_service.go
   func (s *pluginService) ValidatePlugin(ctx context.Context,
       req *plugin.ValidatePluginRequest) (*plugin.ValidatePluginResponse, error) {

       // Load plugin
       // Run conformance tests
       // Generate report
       // Return validation result
   }
   ```

3. Implement Analysis Service
   ```go
   // internal/service/analysis_service.go
   func (s *analysisService) GetRecommendations(ctx context.Context,
       req *analysis.RecommendationRequest,
       stream analysis.GetRecommendationsServerStream) error {

       // Analyze resources
       // Stream recommendations
       // Track progress
   }
   ```

4. Write unit tests
   - Test each method in isolation
   - Mock external dependencies
   - Cover error cases
   - Achieve >80% coverage

**Acceptance Criteria:**
- All service methods implemented
- Unit tests passing
- >80% test coverage
- Error handling comprehensive

#### Week 4: Adapter Layer Implementation

**Deliverables:**
- [ ] PulumiCost adapter
- [ ] Plugin adapter
- [ ] Spec adapter
- [ ] Integration tests

**Tasks:**
1. PulumiCost Core Adapter
   ```go
   // internal/adapter/pulumicost_adapter.go
   type PulumiCostAdapter struct {
       corePath string
       executor *exec.Commander
   }

   func (a *PulumiCostAdapter) GetProjectedCost(ctx context.Context,
       pulumiJSON string) (*types.CostResult, error) {

       // Write JSON to temp file
       // Execute pulumicost binary
       // Parse output
       // Return structured result
   }
   ```

2. Plugin Manager Adapter
   ```go
   // internal/adapter/plugin_adapter.go
   type PluginAdapter struct {
       pluginDir string
       manager   *plugin.Manager
       clients   map[string]pulumicostpb.CostSourceServiceClient
   }

   func (a *PluginAdapter) DiscoverPlugins(ctx context.Context) ([]*Plugin, error) {
       // Scan plugin directory
       // Load plugin metadata
       // Establish gRPC connections
       // Health check
   }
   ```

3. Spec Validator Adapter
   ```go
   // internal/adapter/spec_adapter.go
   func (a *SpecAdapter) ValidatePlugin(ctx context.Context,
       pluginPath string, level ConformanceLevel) (*ValidationReport, error) {

       // Load plugin
       // Run conformance suite
       // Generate report
   }
   ```

**Acceptance Criteria:**
- All adapters implemented
- Integration tests passing
- Error handling robust
- Timeouts and retries configured

### Phase 3: MCP Integration (Week 5)

**Deliverables:**
- [ ] MCP server implementation
- [ ] Claude Desktop integration
- [ ] Example queries
- [ ] Integration documentation

**Tasks:**
1. Server Setup
   ```go
   // cmd/pulumicost-mcp/main.go
   func main() {
       // Load configuration
       // Initialize services
       // Create Goa server
       // Setup MCP transport
       // Start server
   }
   ```

2. Configuration Management
   ```yaml
   # config.yaml
   server:
     port: 8080
     host: localhost

   pulumicost:
     core_path: /usr/local/bin/pulumicost
     plugin_dir: ~/.pulumicost/plugins

   mcp:
     enable_streaming: true
   ```

3. Claude Desktop Integration
   ```json
   // claude_desktop_config.json
   {
     "mcpServers": {
       "pulumicost": {
         "command": "/path/to/pulumicost-mcp",
         "args": ["--config", "config.yaml"]
       }
     }
   }
   ```

4. Example Queries
   - Cost analysis examples
   - Plugin validation examples
   - Resource analysis examples
   - Optimization examples

**Acceptance Criteria:**
- Server starts successfully
- Claude Desktop connects
- All tools discoverable
- Example queries work

### Phase 4: Testing & Documentation (Week 6)

**Deliverables:**
- [ ] Comprehensive test suite
- [ ] User documentation
- [ ] Developer documentation
- [ ] Example plugins

**Tasks:**
1. Testing
   - End-to-end tests
   - Performance tests
   - Load tests
   - Conformance tests

2. Documentation
   - User guide
   - Developer guide
   - API reference
   - Plugin development guide

3. Examples
   - Sample Pulumi stacks
   - Example cost queries
   - Plugin templates
   - Integration examples

**Acceptance Criteria:**
- Test coverage >85%
- All documentation complete
- Examples working
- No critical bugs

### Phase 5: Polish & Release (Week 7-8)

#### Week 7: Performance & Reliability

**Tasks:**
1. Performance optimization
   - Profile and optimize hot paths
   - Implement caching
   - Optimize plugin communication
   - Reduce allocations

2. Reliability improvements
   - Circuit breakers
   - Retry logic
   - Graceful degradation
   - Health checks

3. Observability
   - Metrics (Prometheus)
   - Logging (structured JSON)
   - Tracing (OpenTelemetry)
   - Dashboards (Grafana)

**Acceptance Criteria:**
- P95 latency <3s for queries
- <1% error rate
- Graceful handling of failures
- Full observability

#### Week 8: Launch Preparation

**Tasks:**
1. Release artifacts
   - Binaries for major platforms
   - Docker images
   - Homebrew formula
   - Debian packages

2. Launch materials
   - Blog post
   - Demo video
   - GitHub README
   - Launch announcement

3. Community setup
   - GitHub templates
   - Contributing guide
   - Code of conduct
   - Discussion forums

**Acceptance Criteria:**
- Release artifacts ready
- Documentation polished
- Community infrastructure setup
- Launch plan finalized

## Technology Stack Details

### Core Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| goa.design/goa/v3 | v3.20.2 | API framework |
| goa.design/goa-ai | v0.1.0 | MCP extensions |
| github.com/mark3labs/mcp-go | v0.42.0 | MCP protocol |
| github.com/pulumi/pulumi/sdk/v3 | v3.204.0 | Pulumi integration |
| google.golang.org/grpc | v1.72.1 | gRPC for plugins |
| google.golang.org/protobuf | v1.36.6 | Protocol buffers |

### Development Tools

- **Go 1.24**: Latest Go toolchain
- **buf**: Protocol buffer management
- **golangci-lint**: Comprehensive linting
- **gotestsum**: Better test output
- **mockery**: Mock generation
- **gofumpt**: Code formatting

## Quality Gates

### Code Quality
- [ ] All tests passing
- [ ] >85% test coverage
- [ ] Zero linter errors
- [ ] No security vulnerabilities
- [ ] All TODOs addressed

### Documentation Quality
- [ ] README complete
- [ ] All APIs documented
- [ ] Examples working
- [ ] Role prompts complete
- [ ] Architecture documented

### Performance
- [ ] P95 latency <3s
- [ ] Memory usage <512MB
- [ ] Binary size <50MB
- [ ] Plugin calls <30s timeout
- [ ] Cache hit rate >70%

### Reliability
- [ ] Graceful shutdown
- [ ] Circuit breakers implemented
- [ ] Retry logic in place
- [ ] Health checks working
- [ ] Error recovery tested

## Risk Management

### Technical Risks

1. **Goa-AI Maturity**
   - Risk: New library, may have bugs
   - Mitigation: Extensive testing, contribute fixes
   - Contingency: Fork if needed

2. **Plugin Performance**
   - Risk: Slow plugins block requests
   - Mitigation: Timeouts, circuit breakers, caching
   - Contingency: Async plugin queries

3. **Schema Drift**
   - Risk: pulumicost-core changes break integration
   - Mitigation: Version pinning, integration tests
   - Contingency: Adapter versioning

### Schedule Risks

1. **Learning Curve**
   - Risk: Team unfamiliar with Goa
   - Mitigation: Training, examples, pair programming
   - Contingency: Extra time budgeted

2. **Scope Creep**
   - Risk: Feature requests delay launch
   - Mitigation: Strict scope control, phased approach
   - Contingency: MVP first, iterate

3. **Dependency Updates**
   - Risk: Breaking changes in dependencies
   - Mitigation: Careful update review, tests
   - Contingency: Pin versions, delay updates

## Success Metrics

### Adoption Metrics (Month 1)
- 50 installations
- 10 active users
- 100 cost queries/day
- 5 GitHub stars

### Engagement Metrics (Month 3)
- 200 installations
- 50 active users
- 1000 cost queries/day
- 3 community plugins
- 20 GitHub stars

### Quality Metrics (Ongoing)
- >99% uptime
- <3s P95 response time
- <1% error rate
- >4.5 user satisfaction
- <24h bug response time

## Next Steps

1. **Immediate (This Week)**
   - [x] Create project structure
   - [x] Write role prompts
   - [x] Setup Go modules
   - [ ] Create initial Goa designs
   - [ ] Setup Makefile

2. **Short Term (Next 2 Weeks)**
   - [ ] Complete Goa designs
   - [ ] Implement service layer
   - [ ] Implement adapter layer
   - [ ] Write unit tests

3. **Medium Term (Next 4 Weeks)**
   - [ ] MCP integration
   - [ ] Claude Desktop testing
   - [ ] Documentation
   - [ ] Performance optimization

4. **Long Term (Next 8 Weeks)**
   - [ ] Beta release
   - [ ] Community feedback
   - [ ] Plugin ecosystem
   - [ ] 1.0 release

## Appendix

### Useful Commands

```bash
# Generate code from design
make generate

# Run tests
make test

# Run linter
make lint

# Build binary
make build

# Run server
make run

# Clean generated files
make clean

# Full validation
make validate
```

### Key Files

- `design/design.go`: Main API design
- `design/types.go`: Type definitions
- `cmd/pulumicost-mcp/main.go`: Server entry point
- `internal/service/`: Business logic
- `internal/adapter/`: External integrations
- `gen/`: Generated code (gitignored)
- `Makefile`: Build automation
- `config.yaml`: Configuration

### Resources

- [Goa Documentation](https://goa.design)
- [Goa-AI Repository](https://github.com/goa-ai/goa-ai)
- [MCP Specification](https://modelcontextprotocol.io)
- [pulumicost-core](https://github.com/rshade/pulumicost-core)
- [pulumicost-spec](https://github.com/rshade/pulumicost-spec)

---

**Document Version**: 1.0.0
**Last Updated**: 2025-01-04
**Status**: Draft → Ready for Implementation
