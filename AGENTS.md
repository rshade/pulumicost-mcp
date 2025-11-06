# AGENTS.md - Development Guidelines for PulumiCost MCP Server

## Build, Lint, and Test Commands

### Core Commands
- **Build**: `make build` (generates code and builds binary)
- **Test all**: `make test` (runs all tests with race detection)
- **Test coverage**: `make test-coverage` (generates coverage.html)
- **Lint**: `make lint` (golangci-lint with 40+ rules)
- **Validate**: `make validate` (lint + test)
- **Generate**: `make generate` (Goa code generation)

### Running Specific Tests
- **Single test**: `go test -v -run TestName ./internal/service`
- **Unit tests only**: `make test-unit`
- **Integration tests**: `make test-integration`
- **Benchmarks**: `make bench`

## Code Style Guidelines

### Formatting and Imports
- **Formatter**: `gofumpt -w .` (stricter than gofmt)
- **Imports**: `goimports` with local prefix `github.com/rshade/pulumicost-mcp`
- **Line length**: No hard limit, but prefer readable wrapping
- **Dot imports**: Allowed only in `design/` for Goa DSL

### Naming Conventions
- **Variables**: camelCase, descriptive names
- **Functions**: PascalCase for exported, camelCase for internal
- **Types**: PascalCase, descriptive and specific
- **Receivers**: Single letter (s, a, c) for structs
- **Constants**: ALL_CAPS with underscores

### Error Handling
- **Return errors**: Never panic, always return errors
- **Error wrapping**: Use `fmt.Errorf("context: %w", err)`
- **Context propagation**: Always pass `context.Context` for cancellation
- **Error types**: Use Goa's error types for HTTP/RPC mapping

### Type Safety
- **Explicit types**: Never use `map[string]interface{}`
- **Compiler verification**: Let compiler catch integration errors
- **Goa DSL**: All APIs defined in design before implementation

### Code Structure
- **Design-first**: API contracts in `design/` before implementation
- **Layered architecture**: Design → Generated → Service → Adapter
- **Separation of concerns**: Each layer has single responsibility
- **Generated code**: Never edit files in `gen/` directory

### Testing
- **Test coverage**: Target 80%+ coverage
- **Test naming**: `TestFunctionName` or `TestType_MethodName`
- **Table-driven tests**: Use for multiple scenarios
- **Mock external deps**: Isolate unit tests from external systems

### Documentation
- **Comments**: All exported items must have doc comments
- **Examples**: Include usage examples in comments
- **README**: Keep up to date with implementation changes

### Commit Messages
```
type(scope): subject

body

footer
```
Types: feat, fix, docs, style, refactor, test, chore

### MCP Protocol Compliance
- **JSON-RPC**: All tools properly annotated
- **Streaming**: SSE support for long-running operations
- **Tool registration**: Automatic discovery by MCP clients
- **Error codes**: Standard MCP error responses</content>
<parameter name="filePath">AGENTS.md