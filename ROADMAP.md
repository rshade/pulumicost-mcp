# ROADMAP

This roadmap outlines the development progress and future goals for the PulumiCost MCP Server. All features and goals strictly adhere to the technical guardrails defined in [CONTEXT.md](./CONTEXT.md).

## Phase 1: Foundation (Weeks 1-2) - [Done]
Establish the project infrastructure and design-first API definitions.

- [x] **CI/CD Pipeline** ([#1](https://github.com/rshade/pulumicost-mcp/issues/1)) - GitHub Actions for automated linting and testing.
- [x] **Linting Configuration** ([#2](https://github.com/rshade/pulumicost-mcp/issues/2)) - Comprehensive `golangci-lint` rules for code quality.
- [x] **Testing Framework** ([#3](https://github.com/rshade/pulumicost-mcp/issues/3)) - Integration testing infrastructure and mock adapters.
- [x] **Build Automation** ([#4](https://github.com/rshade/pulumicost-mcp/issues/4)) - Enhanced Makefile with cross-platform targets and Claude Desktop helpers.
- [x] **Goa Service Definitions** ([#5](https://github.com/rshade/pulumicost-mcp/issues/5)) - Design-first DSL for Cost, Plugin, and Analysis services.
- [x] **Initial Code Generation** ([#6](https://github.com/rshade/pulumicost-mcp/issues/6)) - Scaffolding transport, encoding, and protocol layers via Goa-AI.

## Phase 2: Core Implementation (Weeks 3-4) - [Done]
Implement the business logic and external adapters.

- [x] **Cost Service** ([#7](https://github.com/rshade/pulumicost-mcp/issues/7)) - Business logic for projected vs. actual cost analysis.
- [x] **Plugin Service** ([#8](https://github.com/rshade/pulumicost-mcp/issues/8)) - Management and health monitoring for cost source plugins.
- [x] **Analysis Service** ([#9](https://github.com/rshade/pulumicost-mcp/issues/9)) - Optimization recommendations and forecasting logic.
- [x] **PulumiCost Adapter** ([#10](https://github.com/rshade/pulumicost-mcp/issues/10)) - CLI-based orchestration with `pulumicost-core`.
- [x] **Plugin Adapter** ([#11](https://github.com/rshade/pulumicost-mcp/issues/11)) - gRPC communication with circuit-breaking for plugin resilience.
- [x] **Spec Adapter** ([#12](https://github.com/rshade/pulumicost-mcp/issues/12)) - Conformance validation for `pulumicost-spec` compatibility.

## Phase 3: MCP Integration (Week 5) - [Done]
Expose capabilities to AI assistants via the Model Context Protocol.

- [x] **MCP Server Implementation** ([#13](https://github.com/rshade/pulumicost-mcp/issues/13)) - MCP-compliant server with dynamic tool registration.
- [x] **Claude Desktop Config** ([#14](https://github.com/rshade/pulumicost-mcp/issues/14)) - Automated installer for seamless local integration.

## Phase 4: Testing & Documentation (Week 6) - [Done]
Ensure reliability and developer readiness.

- [x] **E2E Test Suite** ([#15](https://github.com/rshade/pulumicost-mcp/issues/15)) - Full protocol-level tests covering all 14 MCP tools.
- [x] **User/Dev Documentation** ([#16](https://github.com/rshade/pulumicost-mcp/issues/16)) - Comprehensive guides, troubleshooting, and example query library.

## Phase 5: Production Readiness (Weeks 7-8) - [In Progress]
Polishing, observability, and public release.

- [x] **Observability** ([#18](https://github.com/rshade/pulumicost-mcp/issues/18)) - Integrated slog (JSON), Prometheus metrics, and OpenTelemetry tracing.
- [ ] **Performance & Benchmarking** ([#17](https://github.com/rshade/pulumicost-mcp/issues/17)) - Validating <3s P95 latency and load testing 50+ concurrent users.
- [ ] **Beta Release Artifacts** ([#19](https://github.com/rshade/pulumicost-mcp/issues/19)) - Cross-platform binaries and Docker image distribution.

---

## Next Horizon (Proposed)
Future improvements focused on presentation, UX, and workflow integration.

- [ ] **Cost Trend Visualization** - ASCII-based sparklines in CLI/Tool output for historical cost trends.
- [ ] **Watch Mode** - Real-time re-analysis of Pulumi previews when local state changes.
- [ ] **CI/CD Formatter** - Specialized output modes for GitHub Action comments (PR-friendly summaries).
- [ ] **Cost-to-Context Enrichment** - Proactive injection of high-level cost metadata into MCP conversation context.
- [ ] **Cross-Plugin Health Dashboard** - Unified view of 3rd-party cost plugin status and circuit-breaker metrics.

## Research & Development (Detailed Horizon)
These items focus on architectural exploration to solve complex FinOps-to-AI interaction challenges while maintaining strict logic delegation.

### 1. Mixed-Currency Aggregation Strategy
*   **Objective**: Solve the "Global Spend" problem where multi-regional infrastructure returns non-additive costs in different currencies.
*   **Technical Approach**: Analyze `CostResult` schemas to implement MCP "Sampling" for user currency preference.
*   **Anti-Guess Boundary**: The MCP Server is strictly forbidden from performing currency math or lookups; it must only group, label, and display values according to the currency codes returned by the orchestrator.
*   **Success Criteria**: A technical specification for grouping mixed-currency stacks by currency in the final AI response.

### 2. High-Latency Streaming Progress (UX)
*   **Objective**: Prevent client-side timeouts and user frustration during long-running (>30s) analysis of massive Pulumi stacks.
*   **Technical Approach**: Integrate MCP `notifications/progress` with `pulumicost-core` ndjson streaming output.
*   **Anti-Guess Boundary**: The MCP Server MUST NOT estimate or "hallucinate" progress percentages based on resource count; it must only pass through the specific progress increments emitted by the core orchestrator.
*   **Success Criteria**: Functional prototype showing a real-time progress bar in the MCP Inspector that mirrors orchestrator status.

### 3. Discovery-Driven Tool Hinting
*   **Objective**: Prevent AI "hallucinations" regarding available plugin capabilities (e.g., suggesting Infracost logic when only Kubecost is present).
*   **Technical Approach**: Dynamically update tool descriptions in the `tools/list` MCP endpoint based on discovered `plugin.json` metadata.
*   **Anti-Guess Boundary**: The MCP Server MUST NOT hardcode or assume the capabilities of a plugin; it must rely entirely on the metadata provided in the plugin's own capability declaration.
*   **Success Criteria**: AI tool descriptions automatically update their specific advice/context when plugins are added/removed.

### 4. Proactive Budget Sentinel (Notifications)
*   **Objective**: Transition from reactive queries to proactive AI alerts when cost deltas exceed thresholds.
*   **Technical Approach**: Research MCP SSE (Server-Sent Events) for pushing budget alerts triggered by core orchestrator events.
*   **Anti-Guess Boundary**: The MCP Server is forbidden from storing budget thresholds or performing budget-to-cost comparisons; it must act solely as a delivery vehicle for alerts generated by the core engine.
*   **Success Criteria**: Protocol log showing an unprompted budget alert being pushed from the server to the client.
