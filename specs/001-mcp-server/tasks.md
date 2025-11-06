# Implementation Tasks: PulumiCost MCP Server

**Feature**: 001-mcp-server
**Branch**: `001-mcp-server`
**Generated**: 2025-01-06
**Total Phases**: 6 (Setup + Foundation + 3 User Stories + Polish)

## Task Summary

**Total Tasks**: 88
**By Phase**:

- Phase 1 (Setup): 10 tasks
- Phase 2 (Foundation): 12 tasks
- Phase 3 (US1 - Cost Analysis): 24 tasks
- Phase 4 (US2 - Plugin Management): 18 tasks
- Phase 5 (US3 - Optimization & Analytics): 16 tasks
- Phase 6 (Polish): 8 tasks

**Parallel Opportunities**: 42 parallelizable tasks marked with [P]

## Implementation Strategy

**MVP Scope**: Phase 1 + Phase 2 + Phase 3 (User Story 1)

This delivers the core value proposition - AI-powered cost analysis for
infrastructure planning. Users can ask Claude about projected and actual costs,
compare configurations, and analyze resources.

**Incremental Delivery**:

1. **Week 1-2**: MVP (Phases 1-3) - Cost Query Service functional
2. **Week 3**: Phase 4 - Plugin ecosystem support
3. **Week 4**: Phase 5 - AI-powered optimization features
4. **Week 5**: Phase 6 - Performance tuning and observability

## User Story Dependencies

```text
Phase 1: Setup
  └─► Phase 2: Foundation (BLOCKS ALL)
        ├─► Phase 3: User Story 1 (P1) - Cost Analysis [INDEPENDENT]
        ├─► Phase 4: User Story 2 (P2) - Plugin Management [INDEPENDENT]
        └─► Phase 5: User Story 3 (P3) - Optimization [INDEPENDENT]
              └─► Phase 6: Polish & Cross-Cutting
```

**Key Insight**: User Stories 1, 2, and 3 can be developed in parallel after
Foundation phase completes. They are designed to be independently testable.

---

## Phase 1: Setup & Project Initialization

**Goal**: Bootstrap project infrastructure and development environment

**Tasks**: 10

- [ ] T001 Create Go module with dependencies in go.mod (Go 1.24, Goa v3,
  Goa-AI, gRPC, Pulumi SDK, mcp-go)
- [ ] T002 [P] Setup Makefile with targets: setup, generate, build, test, lint,
  validate, run, clean
- [ ] T003 [P] Create .golangci.yml configuration from CLAUDE.md specification
  (40+ linters, skip gen/, allow dot imports in design/)
- [ ] T004 [P] Create config.yaml.example with server, pulumicost, and MCP
  configuration sections
- [ ] T005 [P] Setup .github/workflows/ci.yml (generate-check, test, lint,
  build jobs)
- [ ] T006 [P] Create .gitignore for Go (gen/, build/, *.test, coverage.out)
- [ ] T007 [P] Create directory structure: design/, cmd/pulumicost-mcp/,
  internal/service/, internal/adapter/, internal/config/, examples/
- [ ] T008 [P] Create README.md with project overview, quick start, and build
  instructions
- [ ] T009 [P] Create CONTRIBUTING.md with development workflow and design-first
  principles
- [ ] T010 Verify `make setup` installs all development tools (goa, golangci-lint,
  gotestsum, mockery)

**Validation**: `make setup && make generate && make validate` should complete
successfully

---

## Phase 2: Foundation - Goa Design & Code Generation

**Goal**: Define all APIs in Goa DSL and generate transport/validation code.
This phase BLOCKS all user story implementation.

**Tasks**: 12

### Goa Design Files (Design Layer)

- [ ] T011 Create design/design.go with API metadata, MCP server configuration,
  and three service definitions
- [ ] T012 [P] Define all core types in design/types.go: CostQuery, CostResult,
  ResourceCost, ResourceFilter, TagFilter, TimeRange, CostMetadata (14 types
  total from data-model.md)
- [ ] T013 [P] Define plugin types in design/types.go: Plugin, PluginCapabilities,
  HealthStatus, PluginValidationReport, ValidationTest (5 types)
- [ ] T014 [P] Define analysis types in design/types.go: Recommendation, Anomaly,
  Forecast, ForecastPoint, Budget (5 types)
- [ ] T015 Create design/cost_service.go with Cost Query Service definition (6
  MCP tools from cost-service.yaml)
- [ ] T016 Create design/plugin_service.go with Plugin Management Service
  definition (4 MCP tools from plugin-service.yaml)
- [ ] T017 Create design/analysis_service.go with Analysis Service definition (4
  MCP tools from analysis-service.yaml)
- [ ] T018 Run `goa gen` to generate code in gen/cost/, gen/plugin/,
  gen/analysis/, gen/http/, gen/jsonrpc/, gen/mcp/
- [ ] T019 Verify generated code compiles without errors (`go build ./...`)
- [ ] T020 [P] Create internal/config/config.go with configuration loading
  (server, pulumicost, plugins sections)
- [ ] T021 [P] Create internal/config/config_test.go with configuration
  validation tests
- [ ] T022 Commit design/ and gen/ together per Constitution Principle I (never
  commit design without gen or vice versa)

**Validation**: `make generate` produces valid code, `go build ./...` succeeds,
gen/ directory contains all expected packages

---

## Phase 3: User Story 1 (P1) - AI-Powered Cost Analysis

**Goal**: Enable DevOps engineers to query infrastructure costs via AI assistants

**Independent Test**: Deploy Pulumi stack, export preview JSON, ask Claude "What
are projected costs?", receive accurate breakdown. This phase delivers the MVP.

**Tasks**: 24

### Test-First: Unit Tests

- [ ] T023 [P] [US1] Create internal/service/cost_service_test.go with
  TestAnalyzeProjected (RED test for FR-001)
- [ ] T024 [P] [US1] Add TestGetActual test (RED test for FR-002)
- [ ] T025 [P] [US1] Add TestCompareCosts test (RED test for FR-003)
- [ ] T026 [P] [US1] Add TestAnalyzeResource test (RED test for FR-004)
- [ ] T027 [P] [US1] Add TestQueryByTags test (RED test for FR-005)
- [ ] T028 [P] [US1] Add TestAnalyzeStack test with streaming (RED test for
  FR-006)

### Adapter Layer: PulumiCost Integration

- [ ] T029 [P] [US1] Create internal/adapter/pulumicost_adapter_test.go with
  TestGetProjectedCost (RED integration test)
- [ ] T030 [US1] Implement internal/adapter/pulumicost_adapter.go with
  GetProjectedCost (executes pulumicost-core binary, parses JSON output)
- [ ] T031 [US1] Run adapter tests - verify GREEN
- [ ] T032 [P] [US1] Add GetActualCost method to pulumicost_adapter.go (calls
  cloud provider APIs via pulumicost-core)
- [ ] T033 [P] [US1] Add mock implementation in
  internal/adapter/pulumicost_adapter_mock.go for service testing

### Service Layer: Business Logic

- [ ] T034 [US1] Implement internal/service/cost_service.go constructor with
  dependency injection (adapter, logger)
- [ ] T035 [US1] Implement AnalyzeProjected method (calls adapter, transforms to
  CostResult, handles errors with context)
- [ ] T036 [US1] Run TestAnalyzeProjected - verify GREEN, refactor if needed
- [ ] T037 [P] [US1] Implement GetActual method with time range and granularity
  handling
- [ ] T038 [P] [US1] Run TestGetActual - verify GREEN
- [ ] T039 [P] [US1] Implement CompareCosts method (baseline vs target
  comparison, percentage calc)
- [ ] T040 [P] [US1] Run TestCompareCosts - verify GREEN
- [ ] T041 [P] [US1] Implement AnalyzeResource method with URN parsing and
  dependency tracking
- [ ] T042 [P] [US1] Run TestAnalyzeResource - verify GREEN
- [ ] T043 [P] [US1] Implement QueryByTags method with tag grouping and
  aggregation
- [ ] T044 [P] [US1] Run TestQueryByTags - verify GREEN
- [ ] T045 [US1] Implement AnalyzeStack method with SSE streaming support for
  large stacks
- [ ] T046 [US1] Run TestAnalyzeStack - verify GREEN

### Server Integration

- [ ] T047 [US1] Create cmd/pulumicost-mcp/main.go with MCP server
  initialization, wire up Cost Service with adapters
- [ ] T048 [US1] Add graceful shutdown handling (SIGTERM/SIGINT) per FR-020
- [ ] T049 [US1] Run server locally, test with curl against JSON-RPC endpoints
- [ ] T050 [US1] Create examples/queries/cost-analysis-queries.md with 10+
  example questions for Claude

**Phase 3 Validation**: Start server, configure Claude Desktop with MCP config,
ask "What are projected costs for this Pulumi stack?", receive accurate
cost breakdown. Test all 5 acceptance scenarios from spec.md.

**Coverage Target**: >80% for cost_service.go and pulumicost_adapter.go

---

## Phase 4: User Story 2 (P2) - Plugin Management

**Goal**: Enable platform engineers to discover, validate, and manage plugins

**Independent Test**: Install plugin, ask Claude "List available plugins" and
"Validate kubecost plugin", receive plugin info and validation results.

**Tasks**: 18

### Test-First: Plugin Service Unit Tests

- [ ] T051 [P] [US2] Create internal/service/plugin_service_test.go with
  TestListPlugins (RED test for FR-007)
- [ ] T052 [P] [US2] Add TestGetPluginInfo test
- [ ] T053 [P] [US2] Add TestValidatePlugin test (RED test for FR-008)
- [ ] T054 [P] [US2] Add TestHealthCheck test

### Adapter Layer: gRPC Plugin Manager

- [ ] T055 [P] [US2] Create internal/adapter/plugin_adapter_test.go with
  TestDiscoverPlugins (RED integration test)
- [ ] T056 [US2] Implement internal/adapter/plugin_adapter.go with
  DiscoverPlugins (scans plugin directory, loads metadata)
- [ ] T057 [US2] Run TestDiscoverPlugins - verify GREEN
- [ ] T058 [P] [US2] Add EstablishConnection method with gRPC dial and
  health check
- [ ] T059 [P] [US2] Add GetPluginCapabilities method (queries plugin via gRPC
  per FR-009)
- [ ] T060 [P] [US2] Implement circuit breaker logic for plugin calls per FR-017
  (prevent cascade failures)

### Adapter Layer: Spec Validator

- [ ] T061 [P] [US2] Create internal/adapter/spec_adapter_test.go with
  TestValidatePlugin (RED test)
- [ ] T062 [US2] Implement internal/adapter/spec_adapter.go with ValidatePlugin
  (runs pulumicost-spec conformance tests)
- [ ] T063 [US2] Run TestValidatePlugin - verify GREEN

### Service Layer: Plugin Service

- [ ] T064 [US2] Implement internal/service/plugin_service.go constructor
- [ ] T065 [US2] Implement ListPlugins method (discovers + health checks all
  plugins)
- [ ] T066 [US2] Run TestListPlugins - verify GREEN
- [ ] T067 [P] [US2] Implement GetPluginInfo method
- [ ] T068 [P] [US2] Implement ValidatePlugin method
- [ ] T069 [P] [US2] Implement HealthCheck method with latency measurement
- [ ] T070 [US2] Wire up Plugin Service in cmd/pulumicost-mcp/main.go
- [ ] T071 [US2] Create examples/queries/plugin-management-queries.md with
  example plugin queries

**Phase 4 Validation**: Install test plugin, ask Claude "What plugins are
available?", receive list with health status. Ask "Validate the test plugin",
receive conformance report.

**Coverage Target**: >80% for plugin_service.go, plugin_adapter.go,
spec_adapter.go

---

## Phase 5: User Story 3 (P3) - Cost Optimization & Analytics

**Goal**: Enable FinOps analysts to get AI-powered recommendations and
budget tracking

**Independent Test**: Ask Claude "What are cost optimization recommendations?"
or "Show me anomalies for last 30 days", receive actionable insights.

**Tasks**: 16

### Test-First: Analysis Service Unit Tests

- [ ] T072 [P] [US3] Create internal/service/analysis_service_test.go with
  TestGetRecommendations (RED test for FR-010)
- [ ] T073 [P] [US3] Add TestDetectAnomalies test (RED test for FR-011)
- [ ] T074 [P] [US3] Add TestForecastCosts test (RED test for FR-012)
- [ ] T075 [P] [US3] Add TestTrackBudget test (RED test for FR-013)

### Service Layer: Analysis Service

- [ ] T076 [US3] Implement internal/service/analysis_service.go constructor
- [ ] T077 [US3] Implement GetRecommendations method (analyzes cost patterns,
  identifies rightsizing/reserved instance/spot opportunities)
- [ ] T078 [US3] Run TestGetRecommendations - verify GREEN
- [ ] T079 [P] [US3] Implement DetectAnomalies method (statistical analysis,
  baseline comparison, severity classification)
- [ ] T080 [P] [US3] Run TestDetectAnomalies - verify GREEN
- [ ] T081 [P] [US3] Implement ForecastCosts method (time-series prediction with
  confidence intervals)
- [ ] T082 [P] [US3] Run TestForecastCosts - verify GREEN
- [ ] T083 [P] [US3] Implement TrackBudget method (burn rate calc, threshold
  alerts per FR-013)
- [ ] T084 [P] [US3] Run TestTrackBudget - verify GREEN
- [ ] T085 [US3] Wire up Analysis Service in cmd/pulumicost-mcp/main.go
- [ ] T086 [US3] Create examples/queries/optimization-queries.md with 10+
  example analytics questions
- [ ] T087 [US3] Add edge case handling: unsupported resource types, missing
  data, currency normalization

**Phase 5 Validation**: Query stack for recommendations, receive actionable
suggestions with estimated savings. Track budget, receive burn rate and alerts.

**Coverage Target**: >80% for analysis_service.go

---

## Phase 6: Polish & Cross-Cutting Concerns

**Goal**: Production readiness - observability, performance, documentation

**Tasks**: 8

- [ ] T088 [P] Add structured logging throughout services (JSON format, log
  levels per FR-018)
- [ ] T089 [P] Add Prometheus metrics collection (request counts, latencies,
  errors per FR-019)
- [ ] T090 [P] Add OpenTelemetry tracing for request flows
- [ ] T091 [P] Performance testing: verify <3s P95 latency for 100-resource
  stacks (SC-001)
- [ ] T092 [P] Load testing: verify 50 concurrent requests without degradation
  (SC-002)
- [ ] T093 [P] Update README.md with complete usage examples, Claude Desktop
  setup, troubleshooting
- [ ] T094 [P] Create examples/pulumi-stacks/simple-aws/ example project with
  queries.md
- [ ] T095 Verify end-to-end: Start server, configure Claude, run all 14 MCP
  tools successfully, check all success criteria (SC-001 through SC-010)

**Phase 6 Validation**: All 10 success criteria from spec.md validated. Server
runs in production with full observability.

---

## Parallel Execution Examples

### Phase 1: Setup (Parallelizable: T002-T009)

Can run in parallel after T001:

```bash
# Terminal 1
T002: make targets

# Terminal 2
T003: golangci config

# Terminal 3
T004: config.yaml.example

# Terminal 4
T005: CI workflow

# Terminal 5
T006: .gitignore
```

### Phase 2: Foundation (Parallelizable: T012-T014, T020-T021)

After T011, parallel type definitions:

```bash
# Terminal 1
T012: Core types (CostQuery, CostResult, ResourceCost...)

# Terminal 2
T013: Plugin types (Plugin, PluginCapabilities...)

# Terminal 3
T014: Analysis types (Recommendation, Anomaly...)
```

### Phase 3: User Story 1 (Parallelizable: T023-T028, T029, T032-T033, T037+)

Tests can be written in parallel:

```bash
# Terminal 1
T023-T028: All 6 service tests

# Terminal 2
T029: Adapter test

# Terminal 3
T032: GetActualCost adapter method

# Terminal 4
T033: Mock adapter
```

Service methods after T034-T036:

```bash
# Terminal 1
T037-T038: GetActual

# Terminal 2
T039-T040: CompareCosts

# Terminal 3
T041-T042: AnalyzeResource

# Terminal 4
T043-T044: QueryByTags
```

### Phase 4: User Story 2 (Parallelizable: T051-T054, T055+)

All tests in parallel, then adapters:

```bash
# Tests
T051-T054: All plugin service tests in parallel

# Adapters (after T056-T057)
Terminal 1: T058 - EstablishConnection
Terminal 2: T059 - GetPluginCapabilities
Terminal 3: T060 - Circuit breaker
Terminal 4: T061-T063 - Spec validator
```

### Phase 5: User Story 3 (Parallelizable: T072-T075, T079+)

All tests and service methods can run in parallel:

```bash
# Tests
T072-T075: All analysis tests in parallel

# Service methods (after T076-T078)
Terminal 1: T079-T080 - DetectAnomalies
Terminal 2: T081-T082 - ForecastCosts
Terminal 3: T083-T084 - TrackBudget
```

### Phase 6: Polish (Parallelizable: T088-T094)

All polish tasks can run in parallel:

```bash
# Terminal 1
T088: Structured logging

# Terminal 2
T089: Prometheus metrics

# Terminal 3
T090: OpenTelemetry tracing

# Terminal 4
T091: Performance testing

# Terminal 5
T092: Load testing

# Terminal 6
T093: README updates

# Terminal 7
T094: Example projects
```

---

## Task Execution Checklist

When executing tasks:

1. **Read the spec** - Understand acceptance criteria before starting
2. **Write failing test** - RED phase (TDD)
3. **Implement minimum code** - GREEN phase
4. **Refactor** - Clean up while tests stay green
5. **Run `make validate`** - Ensure lint + tests pass
6. **Update coverage** - Maintain >80% target
7. **Commit atomically** - One task = one commit (design + gen together if
  applicable)
8. **Mark task complete** - Check off in this file

## Coverage Targets by Phase

- Phase 2 (Foundation): Config >80%
- Phase 3 (US1): cost_service.go >80%, pulumicost_adapter.go >80%
- Phase 4 (US2): plugin_service.go >80%, plugin_adapter.go >80%,
  spec_adapter.go >80%
- Phase 5 (US3): analysis_service.go >80%
- **Overall Project**: >80% (target 85%+)

## Constitution Compliance

All tasks follow the seven principles:

- **I. Design-First**: T011-T017 define APIs before implementation
- **II. Type Safety**: All types in Goa DSL (T012-T014)
- **III. Test-First**: RED tests before GREEN implementation (T023+)
- **IV. Separation of Concerns**: Four-layer architecture enforced
- **V. Error Handling**: context.Context everywhere, circuit breakers
- **VI. Integration Standards**: Adapter pattern (T029+, T055+, T061+)
- **VII. Documentation as Code**: Examples and README updates

---

**Last Updated**: 2025-01-06
**Status**: Ready for implementation
