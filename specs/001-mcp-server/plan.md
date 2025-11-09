# Implementation Plan: PulumiCost MCP Server

**Branch**: `001-mcp-server` | **Date**: 2025-01-06 | **Spec**:
[spec.md](./spec.md)

**Input**: Feature specification from
`/specs/001-mcp-server/spec.md`

**Note**: This plan is filled in by the `/speckit.plan` command.

## Summary

Build a production-grade MCP server exposing PulumiCost cloud cost analysis to
AI assistants via natural language. The server implements three core services
(Cost Query, Plugin Management, Analysis) using design-first development with
Goa+Goa-AI for type-safe MCP integration, gRPC for plugin communication, and
stateless architecture for horizontal scaling. Primary users are DevOps
engineers, platform engineers, and FinOps analysts who need conversational cost
insights during infrastructure planning and optimization.

## Technical Context

**Language/Version**: Go 1.24

**Primary Dependencies**:

- goa.design/goa/v3 v3.20.2 (API framework)
- goa.design/goa-ai v0.1.0 (MCP extensions)
- github.com/mark3labs/mcp-go v0.42.0 (MCP protocol)
- google.golang.org/grpc v1.72.1 (plugin communication)
- github.com/pulumi/pulumi/sdk/v3 v3.204.0 (Pulumi integration)

**Storage**: N/A (stateless server, no persistent storage)

**Testing**: Go testing framework with gotestsum, mockery for mocks, >80%
coverage target

**Target Platform**: Linux server (primary), macOS and Windows (development)

**Project Type**: Single Go binary server

**Performance Goals**: P95 latency <3s for cost queries, 50+ concurrent
requests, 99% uptime

**Constraints**: <512MB memory footprint, <50MB binary size, 30s plugin timeout,
stateless operation

**Scale/Scope**: 500 concurrent users with horizontal scaling, 19 GitHub issues
across 5 implementation phases

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Design-First Development (NON-NEGOTIABLE)

**Status**: ✅ PASS

**Evidence**: Goa DSL will define all APIs in `design/*.go` before
implementation. Generated code in `gen/` never modified. CI enforces
synchronization.

### Principle II: Type Safety (NON-NEGOTIABLE)

**Status**: ✅ PASS

**Evidence**: All types defined in Goa DSL. No `map[string]interface{}` except
required by external libraries. Compiler catches integration errors.

### Principle III: Test-First Development

**Status**: ✅ PASS

**Evidence**: TDD workflow required for all service methods. >80% coverage
target. Unit, integration, contract, and E2E tests planned across phases 1-4.

### Principle IV: Separation of Concerns

**Status**: ✅ PASS

**Evidence**: Four-layer architecture enforced - Design (`design/*.go`),
Generated (`gen/*`), Service (`internal/service/*`), Adapter
(`internal/adapter/*`).

### Principle V: Error Handling Excellence

**Status**: ✅ PASS

**Evidence**: All service methods accept `context.Context`. Goa error types for
transport mapping. Circuit breakers for plugins (FR-017). Graceful shutdown
(FR-020).

### Principle VI: Integration Standards

**Status**: ✅ PASS

**Evidence**: Adapter pattern for pulumicost-core, plugins, and spec
integrations. Dependency injection. Mock adapters for testing. Circuit breakers
and timeouts.

### Principle VII: Documentation as Code

**Status**: ✅ PASS

**Evidence**: API docs generated from Goa DSL. CLAUDE.md maintained. Examples
in `examples/` directory. README synchronized with capabilities.

**GATE RESULT**: ✅ ALL GATES PASS - Proceed to Phase 0 research

## Project Structure

### Documentation (this feature)

```text
specs/001-mcp-server/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (skipped - tech stack known)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   ├── cost-service.yaml       # Cost Query Service OpenAPI
│   ├── plugin-service.yaml     # Plugin Management Service OpenAPI
│   └── analysis-service.yaml   # Analysis Service OpenAPI
└── tasks.md             # Phase 2 output (/speckit.tasks - NOT by /speckit.plan)
```

### Source Code (repository root)

```text
pulumicost-mcp/
├── design/                    # Goa DSL (Design Layer)
│   ├── design.go             # Main API and MCP server configuration
│   ├── cost_service.go       # Cost Query Service definition
│   ├── plugin_service.go     # Plugin Management Service definition
│   ├── analysis_service.go   # Analysis Service definition
│   └── types.go              # Shared type definitions
│
├── gen/                       # Generated code (Generated Layer)
│   ├── cost/                 # Generated cost service
│   ├── plugin/               # Generated plugin service
│   ├── analysis/             # Generated analysis service
│   ├── http/                 # Generated HTTP transport
│   ├── jsonrpc/              # Generated JSON-RPC transport
│   └── mcp/                  # Generated MCP protocol bindings
│
├── internal/                  # Implementation (not exported)
│   ├── service/              # Service Layer (business logic)
│   │   ├── cost_service.go
│   │   ├── cost_service_test.go
│   │   ├── plugin_service.go
│   │   ├── plugin_service_test.go
│   │   ├── analysis_service.go
│   │   └── analysis_service_test.go
│   │
│   ├── adapter/              # Adapter Layer (external integrations)
│   │   ├── pulumicost_adapter.go      # pulumicost-core integration
│   │   ├── pulumicost_adapter_test.go
│   │   ├── plugin_adapter.go          # gRPC plugin management
│   │   ├── plugin_adapter_test.go
│   │   ├── spec_adapter.go            # pulumicost-spec validation
│   │   └── spec_adapter_test.go
│   │
│   └── config/               # Configuration management
│       ├── config.go
│       └── config_test.go
│
├── cmd/
│   └── pulumicost-mcp/       # Main server entry point
│       └── main.go
│
├── examples/
│   ├── pulumi-stacks/        # Example Pulumi projects
│   ├── queries/              # Example MCP queries
│   └── claude-desktop/       # Claude Desktop config examples
│
├── docs/                      # Documentation
├── scripts/                   # Build and deployment scripts
├── .github/workflows/         # CI/CD pipelines
├── Makefile                   # Build automation
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
└── config.yaml.example        # Configuration template
```

**Structure Decision**: Single Go binary server project. The design-first
approach with Goa generates transport/validation code in `gen/`, keeping
implementation in `internal/` following four-layer architecture. No
frontend/backend split needed - this is a pure MCP server exposing JSON-RPC
endpoints.

## Complexity Tracking

**No violations** - All constitution principles satisfied.

The design-first approach with Goa, four-layer architecture, and adapter pattern
align perfectly with all seven constitutional principles. No complexity
justification required.
