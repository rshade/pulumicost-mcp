# CONTEXT.md

## Core Architectural Identity
**Lightweight MCP Server & Infrastructure Cost Presentation Layer**  
This project is a specialized bridge between AI Assistants (e.g., Claude, ChatGPT) and the PulumiCost cost analysis engine. It leverages the Model Context Protocol (MCP) and the Goa-AI framework to transform complex infrastructure data into actionable, tool-based insights for LLMs.

## Technical Boundaries
1. **No Native Calculation**: This project NEVER performs its own cost calculations, math, or currency conversions. All financial data is treated as a "pass-through" from downstream services.
2. **No Direct Cloud Mutability**: The server is strictly a read-only analysis layer. It does not have the authority or capability to create, modify, or delete cloud resources.
3. **Stateless Presentation**: There is no internal database. All historical data, trends, and anomalies are retrieved dynamically from `pulumicost-core` or external plugins.
4. **Logic Delegation**: Complex FinOps logic (e.g., "what-if" analysis or normalization) happens in `pulumicost-core`; this project is responsible only for orchestrating the request and formatting the response for the MCP client.

## Data Source of Truth
- **pulumicost-core**: The primary orchestration engine and source of truth for cost analysis logic.
- **Cost Source Plugins**: External gRPC services (found in `~/.pulumicost/plugins/`) responsible for provider-specific pricing data.
- **Pulumi CLI**: The source of infrastructure state and preview JSON data used as input for analysis.

## Interaction Model
- **Northbound (AI Clients)**: Standardized JSON-RPC via the Model Context Protocol (MCP) over Stdio or SSE.
- **Southbound (Orchestrator)**: CLI execution of the `pulumicost-core` binary using JSON-based I/O.
- **East-West (Plugins)**: High-performance gRPC for plugin discovery, health monitoring, and capability negotiation.
- **Observability**: Prometheus for metrics, slog for structured JSON logging, and OpenTelemetry for distributed tracing.
