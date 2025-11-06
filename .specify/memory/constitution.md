<!--
=============================================================================
SYNC IMPACT REPORT
=============================================================================
Version Change: None → 1.0.0 (Initial Constitution)
Modified Principles: N/A (new constitution)
Added Sections:
  - Core Principles (7 principles)
  - Development Workflow
  - Quality Standards
  - Governance

Templates Requiring Updates:
  ✅ spec-template.md - Reviewed: Already compatible (no constitution-specific refs)
  ✅ plan-template.md - Reviewed: "Constitution Check" section aligns with principles
  ✅ tasks-template.md - Reviewed: Testing requirements align with Principle III
  ✅ commands/*.md - Reviewed: No agent-specific references found

Follow-up TODOs: None

=============================================================================
-->

# PulumiCost MCP Server Constitution

## Core Principles

### I. Design-First Development (NON-NEGOTIABLE)

All APIs, services, and contracts MUST be defined in Goa DSL before
implementation. The design files in `design/*.go` are the single source of
truth for all APIs, types, validation rules, and transport bindings.

**Rules**:

- API contracts live exclusively in `design/*.go`
- Generated code in `gen/` MUST NEVER be modified directly
- `make generate` MUST be run after any design change
- Both design files and generated code MUST be committed together
- CI enforces design/generated code synchronization

**Rationale**: Design-first development eliminates schema drift, ensures type
safety throughout the stack, provides compiler-verified contracts, and enables
automatic generation of transport, validation, and documentation code.

### II. Type Safety (NON-NEGOTIABLE)

Strong typing MUST be enforced throughout the entire codebase. No
stringly-typed interfaces, no `map[string]interface{}` except where absolutely
required by external libraries, and all types MUST be defined in Goa DSL.

**Rules**:

- Use explicit types defined in Goa DSL
- Avoid `interface{}` and `map[string]interface{}`
- Compiler MUST catch integration errors
- No type assertions without error checking
- All external data MUST be validated against typed schemas

**Rationale**: Type safety catches errors at compile time, provides excellent
IDE support, enables refactoring confidence, and serves as living documentation.

### III. Test-First Development

Tests MUST be written before implementation for all service methods. Tests
serve as executable specifications and guard against regressions.

**Rules**:

- Write tests that FAIL before implementing functionality
- Follow Red-Green-Refactor cycle
- Unit tests for service layer (>80% coverage target)
- Integration tests for adapters
- Contract tests for external integrations
- End-to-end tests for MCP protocol
- All tests MUST pass before merge

**Rationale**: Test-first development clarifies requirements, drives better API
design, provides regression protection, and serves as living documentation.

### IV. Separation of Concerns

Code MUST be organized in distinct layers with clear responsibilities: Design
Layer (contracts), Generated Layer (transport/validation), Service Layer
(business logic), and Adapter Layer (external integrations).

**Rules**:

- Design Layer: API contracts, types, validation rules (`design/*.go`)
- Generated Layer: Transport, encoding, validation (`gen/*`) - NEVER edited
- Service Layer: Business logic implementation (`internal/service/*`)
- Adapter Layer: External system integration (`internal/adapter/*`)
- Each layer depends only on layers below it
- No business logic in adapters or generated code

**Rationale**: Clear separation enables independent testing, simplifies
maintenance, allows parallel development, and makes the system easier to reason
about.

### V. Error Handling Excellence

Errors MUST be handled explicitly and provide actionable information. All
external calls MUST include timeout and cancellation support via
`context.Context`.

**Rules**:

- Always pass `context.Context` through call stack
- Return errors, never panic (except in init functions)
- Use Goa's error types for proper transport mapping
- Include actionable error messages with context
- Log errors at appropriate levels
- Wrap errors with `fmt.Errorf("%w", err)` for stack traces
- Handle timeouts and cancellation gracefully

**Rationale**: Explicit error handling makes failures visible, context
propagation enables cancellation and timeouts, actionable messages reduce
debugging time, and proper logging aids troubleshooting.

### VI. Integration Standards

External integrations MUST use the Adapter pattern to isolate dependencies and
enable testing. Adapters MUST implement interfaces defined in the service
layer.

**Rules**:

- All external dependencies accessed through adapters
- Adapters implement service-layer interfaces
- Adapters located in `internal/adapter/*`
- Use dependency injection for adapters
- Mock adapters for service-layer testing
- Integration tests for adapter implementations
- Timeouts and circuit breakers for all external calls

**Rationale**: The Adapter pattern isolates external dependencies, enables
testing with mocks, allows swapping implementations, and provides resilience
through circuit breakers.

### VII. Documentation as Code

Documentation MUST be maintained alongside code and generated from
authoritative sources where possible. API documentation is automatically
generated from Goa designs.

**Rules**:

- API documentation generated from Goa DSL
- Update CLAUDE.md with project-specific patterns
- Document architectural decisions in code comments
- Maintain examples in `examples/` directory
- Keep README synchronized with capabilities
- Document all exported functions and types
- Include usage examples in documentation

**Rationale**: Documentation as code keeps documentation synchronized with
implementation, reduces documentation drift, provides single source of truth,
and makes documentation part of the development workflow.

## Development Workflow

### Code Generation Workflow

1. **Modify Design**: Edit files in `design/*.go` to change APIs
2. **Generate Code**: Run `make generate` to produce transport/validation code
3. **Implement Services**: Update service implementations in `internal/service/*`
4. **Update Tests**: Modify tests to match new signatures
5. **Validate**: Run `make validate` (lint + test)
6. **Commit Together**: Commit both design and generated code

**Critical**: Never modify generated code directly. Always start with design changes.

### Testing Workflow

1. **Write Failing Test**: Create test that captures requirement
2. **Verify Failure**: Ensure test fails for the right reason
3. **Implement**: Write minimal code to pass test
4. **Refactor**: Improve code while keeping tests green
5. **Coverage**: Verify >80% coverage with `make test-coverage`
6. **Integration**: Test external integrations separately

### Pull Request Workflow

1. **Branch**: Create feature branch from `develop`
2. **Implement**: Follow design-first workflow
3. **Validate**: Run `make validate` locally
4. **Push**: Push branch and open pull request
5. **CI**: Ensure all CI checks pass (generate-check, test, lint, build)
6. **Review**: Address reviewer feedback
7. **Merge**: Squash and merge to `develop`

## Quality Standards

### Code Quality Requirements

- **All linters pass**: Zero errors from golangci-lint (configured in `.golangci.yml`)
- **Test coverage**: Maintain >80% coverage (target 85%+)
- **No compiler warnings**: All code must compile cleanly
- **No security issues**: gosec linter must pass
- **Formatted code**: Use gofumpt for consistent formatting
- **Clean git history**: Meaningful commit messages following conventional commits

### Performance Requirements

- **P95 latency**: <3s for cost queries
- **Memory usage**: <512MB server footprint
- **Binary size**: <50MB compiled binary
- **Plugin timeout**: 30s maximum for plugin calls
- **Cache efficiency**: >70% cache hit rate for repeated queries

### Reliability Requirements

- **Graceful shutdown**: Handle SIGTERM/SIGINT cleanly
- **Circuit breakers**: Prevent cascade failures from plugins
- **Retry logic**: Exponential backoff for transient failures
- **Health checks**: Expose health endpoints for monitoring
- **Error recovery**: Gracefully handle and log all errors

## Governance

### Amendment Process

Constitution amendments require:

1. **Proposal**: Document proposed change with rationale
2. **Discussion**: Allow time for team review and feedback
3. **Approval**: Consensus from maintainers required
4. **Migration**: Update code to comply with new principles
5. **Documentation**: Update dependent templates and guides

### Version Numbering

- **MAJOR**: Backward incompatible governance changes or principle removal
- **MINOR**: New principles added or material expansions to guidance
- **PATCH**: Clarifications, wording improvements, typo fixes

### Compliance Review

- All pull requests MUST verify constitution compliance
- CI pipeline enforces technical principles (design sync, tests, linting)
- Code reviews verify architectural principles (separation of concerns, type safety)
- Principle violations MUST be justified or corrected

### Complexity Justification

Any deviation from principles MUST be documented in the Implementation Plan
under "Complexity Tracking" with:

- Which principle is violated
- Why the complexity is necessary
- What simpler alternatives were considered and rejected

**Use CLAUDE.md for runtime development guidance and project-specific patterns.**

**Version**: 1.0.0 | **Ratified**: 2025-01-06 | **Last Amended**: 2025-01-06
