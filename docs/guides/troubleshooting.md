# Troubleshooting Guide

Complete troubleshooting guide for PulumiCost MCP Server.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Server Startup Issues](#server-startup-issues)
- [Claude Desktop Integration](#claude-desktop-integration)
- [Tool Execution Errors](#tool-execution-errors)
- [Plugin Issues](#plugin-issues)
- [Performance Problems](#performance-problems)
- [Network and Connectivity](#network-and-connectivity)
- [Data and Results Issues](#data-and-results-issues)
- [Development and Testing](#development-and-testing)
- [Logging and Debugging](#logging-and-debugging)

---

## Installation Issues

### Problem: `make install` fails

**Symptoms**:

```bash
$ make install
make: claude: command not found
```

**Cause**: Claude Desktop CLI is not installed or not in PATH.

**Solutions**:

1. **Install Claude Desktop CLI**:

   ```bash
   # macOS (via Homebrew)
   brew install anthropics/claude/claude

   # Or download from: https://www.anthropic.com/claude
   ```

2. **Verify installation**:

   ```bash
   claude --version
   ```

3. **Check PATH**:

   ```bash
   which claude
   # Should output: /usr/local/bin/claude or similar
   ```

4. **Manual installation** (if CLI unavailable):

   Edit Claude Desktop config manually:

   **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

   **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

   ```json
   {
     "mcpServers": {
       "pulumicost": {
         "command": "/path/to/pulumicost-mcp/bin/pulumicost-mcp",
         "args": []
       }
     }
   }
   ```

---

### Problem: Build fails with Go version error

**Symptoms**:

```bash
$ make build
go: go.mod requires go >= 1.24, but go version is 1.21
```

**Cause**: Go 1.24 or later is required.

**Solutions**:

1. **Upgrade Go**:

   ```bash
   # macOS (via Homebrew)
   brew upgrade go

   # Linux (via go install)
   go install golang.org/dl/go1.24@latest
   go1.24 download

   # Verify
   go version  # Should show go1.24.x
   ```

2. **Update PATH** (if using multiple Go versions):

   ```bash
   export PATH=/usr/local/go/bin:$PATH
   ```

---

### Problem: Missing dependencies during setup

**Symptoms**:

```bash
$ make setup
Error: golangci-lint not found
```

**Cause**: Development tools not installed.

**Solutions**:

1. **Run install-tools target**:

   ```bash
   make install-tools
   ```

2. **Manual installation**:

   ```bash
   # golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.6.0

   # goa and goa-ai
   go install goa.design/goa/v3/cmd/goa@v3.22.6
   go install goa.design/goa-ai/cmd/goa-ai@v0.1.0
   ```

3. **Verify installations**:

   ```bash
   golangci-lint --version
   goa version
   goa-ai version
   ```

---

## Server Startup Issues

### Problem: Server won't start

**Symptoms**:

```bash
$ ./bin/pulumicost-mcp
panic: runtime error: invalid memory address
```

**Diagnostic Steps**:

1. **Check server logs**:

   ```bash
   ./bin/pulumicost-mcp 2>&1 | tee server.log
   ```

2. **Verify binary is up-to-date**:

   ```bash
   make clean
   make build
   ```

3. **Check generated code**:

   ```bash
   make generate
   # If changes detected, rebuild
   make build
   ```

4. **Run with verbose logging**:

   ```bash
   MCP_LOG_LEVEL=debug ./bin/pulumicost-mcp
   ```

---

### Problem: Port already in use

**Symptoms**:

```bash
Error: listen tcp :8080: bind: address already in use
```

**Cause**: Another process is using port 8080.

**Solutions**:

1. **Find conflicting process**:

   ```bash
   # macOS/Linux
   lsof -i :8080
   # Or
   netstat -an | grep 8080
   ```

2. **Kill conflicting process**:

   ```bash
   kill -9 <PID>
   ```

3. **Use different port**:

   ```bash
   MCP_SERVER_PORT=8081 ./bin/pulumicost-mcp
   ```

4. **Update Claude Desktop config**:

   ```json
   {
     "mcpServers": {
       "pulumicost": {
         "command": "/path/to/pulumicost-mcp",
         "env": {
           "MCP_SERVER_PORT": "8081"
         }
       }
     }
   }
   ```

---

### Problem: Permission denied

**Symptoms**:

```bash
$ ./bin/pulumicost-mcp
-bash: ./bin/pulumicost-mcp: Permission denied
```

**Cause**: Binary is not executable.

**Solutions**:

```bash
# Make executable
chmod +x ./bin/pulumicost-mcp

# Verify
ls -l ./bin/pulumicost-mcp
# Should show: -rwxr-xr-x
```

---

## Claude Desktop Integration

### Problem: MCP server not showing up in Claude Desktop

**Symptoms**:

- No "pulumicost" tools available in Claude
- No MCP connection indicators

**Diagnostic Steps**:

1. **Verify installation**:

   ```bash
   claude mcp list
   # Should show: pulumicost
   ```

2. **Check configuration file**:

   **macOS**:

   ```bash
   cat ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

   **Windows**:

   ```cmd
   type %APPDATA%\Claude\claude_desktop_config.json
   ```

   Should contain:

   ```json
   {
     "mcpServers": {
       "pulumicost": {
         "command": "/absolute/path/to/pulumicost-mcp",
         "args": []
       }
     }
   }
   ```

3. **Verify binary path is absolute**:

   ```bash
   # Get absolute path
   realpath ./bin/pulumicost-mcp
   # Or
   readlink -f ./bin/pulumicost-mcp
   ```

4. **Check binary is executable**:

   ```bash
   /absolute/path/to/pulumicost-mcp --version
   ```

5. **Restart Claude Desktop** (required after config changes):

   - macOS: Cmd+Q to quit, then relaunch
   - Windows: Exit from system tray, then relaunch

---

### Problem: "Server failed to start" error

**Symptoms**:

Claude Desktop shows: "MCP server 'pulumicost' failed to start"

**Diagnostic Steps**:

1. **Check Claude Desktop logs**:

   **macOS**:

   ```bash
   tail -f ~/Library/Logs/Claude/mcp-*.log
   ```

   **Windows**:

   ```cmd
   type %APPDATA%\Claude\Logs\mcp-*.log
   ```

2. **Test binary manually**:

   ```bash
   /path/to/pulumicost-mcp
   # Should start without errors
   ```

3. **Check for missing dependencies**:

   ```bash
   ldd /path/to/pulumicost-mcp  # Linux
   otool -L /path/to/pulumicost-mcp  # macOS
   ```

4. **Verify configuration syntax**:

   ```bash
   # Use JSON validator
   cat claude_desktop_config.json | jq .
   # Should output formatted JSON without errors
   ```

---

### Problem: Tools work inconsistently

**Symptoms**:

- Some tools work, others don't
- Intermittent "tool not found" errors

**Solutions**:

1. **Verify server is running**:

   ```bash
   # Check process
   ps aux | grep pulumicost-mcp
   ```

2. **Check server logs** for errors:

   ```bash
   MCP_LOG_LEVEL=debug ./bin/pulumicost-mcp 2>&1 | tee debug.log
   ```

3. **Test tools via curl**:

   ```bash
   curl -X POST http://localhost:8080/rpc \
     -H "Content-Type: application/json" \
     -d '{
       "jsonrpc": "2.0",
       "method": "tools/list",
       "id": 1
     }'
   ```

   Should return list of 14 tools.

4. **Restart both server and Claude Desktop**:

   ```bash
   # Kill server
   pkill pulumicost-mcp

   # Restart Claude Desktop
   # Then it will auto-start the server
   ```

---

## Tool Execution Errors

### Problem: "Validation error" when using tools

**Symptoms**:

```json
{
  "error": "validation_error",
  "message": "stack_name cannot be empty"
}
```

**Cause**: Missing required parameters or invalid input.

**Solutions**:

1. **Check required parameters**:

   See [MCP Tools Reference](mcp-tools.md) for each tool's required parameters.

2. **Verify parameter format**:

   ```bash
   # Correct format for time_range
   {
     "time_range": {
       "start": "2024-01-01T00:00:00Z",  # ISO 8601
       "end": "2024-01-31T23:59:59Z"
     }
   }
   ```

3. **Test with minimal payload**:

   ```bash
   curl -X POST http://localhost:8080/rpc \
     -H "Content-Type: application/json" \
     -d '{
       "jsonrpc": "2.0",
       "method": "cost.analyze_projected",
       "params": {
         "pulumi_json": "{\"resources\": []}"
       },
       "id": 1
     }'
   ```

---

### Problem: "Internal error" from tools

**Symptoms**:

```json
{
  "error": "internal_error",
  "message": "failed to execute pulumicost-core"
}
```

**Diagnostic Steps**:

1. **Check server logs** for detailed error:

   ```bash
   tail -f /tmp/pulumicost-mcp.log
   ```

2. **Verify pulumicost-core** (when implemented):

   ```bash
   # Check if binary exists
   which pulumicost

   # Test directly
   pulumicost --version
   ```

3. **Check adapter configuration**:

   Current implementation uses mock data, so this error indicates a code issue.
   Report at: https://github.com/rshade/pulumicost-mcp/issues

---

### Problem: Timeouts on large stacks

**Symptoms**:

```json
{
  "error": "internal_error",
  "message": "context deadline exceeded"
}
```

**Cause**: Operation took longer than timeout (default: 30s).

**Solutions**:

1. **Increase timeout** via configuration:

   ```yaml
   # config.yaml
   server:
     timeout: 60s  # Increase to 60 seconds
   ```

2. **Use streaming tools** for large operations:

   ```bash
   # Instead of: analyze_projected_cost
   # Use: analyze_stack (with streaming)
   ```

3. **Filter results** to reduce processing:

   ```json
   {
     "filters": {
       "provider": "aws",
       "region": "us-east-1"
     }
   }
   ```

---

## Plugin Issues

### Problem: Plugins not discovered

**Symptoms**:

`list_plugins` returns empty array or missing plugins.

**Diagnostic Steps**:

1. **Check plugin directory**:

   ```bash
   ls -la ~/.pulumicost/plugins/
   # Should show plugin binaries
   ```

2. **Verify plugin directory configuration**:

   ```bash
   echo $PULUMICOST_PLUGIN_DIR
   # Or check config.yaml
   ```

3. **Set correct permissions**:

   ```bash
   chmod +x ~/.pulumicost/plugins/*
   ```

4. **Test plugin directly**:

   ```bash
   ~/.pulumicost/plugins/aws-cost-source --version
   ```

---

### Problem: Plugin health check fails

**Symptoms**:

```json
{
  "status": "unhealthy",
  "error_message": "connection refused"
}
```

**Diagnostic Steps**:

1. **Check if plugin is running**:

   ```bash
   ps aux | grep aws-cost-source
   ```

2. **Verify gRPC address**:

   ```bash
   netstat -an | grep 50051
   # Plugin should be listening
   ```

3. **Test gRPC connection**:

   ```bash
   grpcurl -plaintext localhost:50051 list
   # Should list gRPC services
   ```

4. **Check plugin logs**:

   ```bash
   tail -f ~/.pulumicost/plugins/aws-cost-source.log
   ```

5. **Restart plugin**:

   ```bash
   # Kill plugin
   pkill aws-cost-source

   # Server will auto-restart it on next health check
   ```

---

### Problem: Plugin validation fails

**Symptoms**:

```json
{
  "passed": false,
  "test_results": [
    {
      "name": "gRPC Interface",
      "passed": false,
      "error": "method not found"
    }
  ]
}
```

**Cause**: Plugin doesn't implement required pulumicost-spec interface.

**Solutions**:

1. **Check spec version compatibility**:

   ```bash
   # Plugin should implement same spec version as server
   grep "spec_version" ~/.pulumicost/plugins/plugin-info.yaml
   ```

2. **Run conformance tests**:

   ```bash
   pulumicost-spec validate ~/.pulumicost/plugins/my-plugin
   ```

3. **Review plugin implementation**:

   See [Plugin Development Guide](plugin-development.md)

4. **Update plugin** to latest spec version.

---

## Performance Problems

### Problem: Slow response times

**Symptoms**:

- Tools take >5 seconds to respond
- Claude appears to hang

**Diagnostic Steps**:

1. **Check server resource usage**:

   ```bash
   # CPU and memory
   top -p $(pgrep pulumicost-mcp)
   ```

2. **Enable performance logging**:

   ```bash
   MCP_LOG_LEVEL=debug ./bin/pulumicost-mcp 2>&1 | grep -i "duration"
   ```

3. **Profile the server**:

   ```bash
   # CPU profiling
   go test -cpuprofile=cpu.prof -bench=.

   # Analyze
   go tool pprof cpu.prof
   ```

**Solutions**:

1. **Enable caching** (when implemented):

   ```yaml
   cache:
     enabled: true
     ttl: 300s  # 5 minutes
   ```

2. **Reduce concurrent plugin calls**:

   ```yaml
   plugins:
     max_concurrent: 5  # Reduce from 10
   ```

3. **Filter data** to reduce processing:

   Use `filters` parameter in tool calls.

4. **Upgrade hardware**:

   - More CPU cores for concurrent processing
   - More RAM for caching

---

### Problem: High memory usage

**Symptoms**:

```bash
$ ps aux | grep pulumicost-mcp
rshade  12345  50.2  15.3   # 15.3% memory usage
```

**Diagnostic Steps**:

1. **Profile memory**:

   ```bash
   go test -memprofile=mem.prof -bench=.
   go tool pprof mem.prof
   ```

2. **Check for memory leaks**:

   ```bash
   # Run for 1 hour and monitor
   watch -n 60 'ps aux | grep pulumicost-mcp'
   ```

**Solutions**:

1. **Disable caching** (if enabled):

   ```yaml
   cache:
     enabled: false
   ```

2. **Reduce cache size**:

   ```yaml
   cache:
     max_size: 100MB  # Reduce from 1GB
   ```

3. **Restart server periodically**:

   ```bash
   # Cron job to restart daily
   0 0 * * * pkill pulumicost-mcp
   ```

---

## Network and Connectivity

### Problem: Can't connect to MCP server

**Symptoms**:

```
Connection refused: http://localhost:8080
```

**Diagnostic Steps**:

1. **Check server is running**:

   ```bash
   ps aux | grep pulumicost-mcp
   ```

2. **Check listening port**:

   ```bash
   netstat -an | grep 8080
   # Or
   lsof -i :8080
   ```

3. **Test connection**:

   ```bash
   curl http://localhost:8080/health
   # Or
   telnet localhost 8080
   ```

**Solutions**:

1. **Start server** if not running:

   ```bash
   ./bin/pulumicost-mcp &
   ```

2. **Check firewall**:

   ```bash
   # macOS
   sudo /usr/libexec/ApplicationFirewall/socketfilterfw --listapps

   # Linux
   sudo ufw status
   ```

3. **Verify host/port configuration**:

   ```yaml
   server:
     host: localhost  # Or 0.0.0.0 for all interfaces
     port: 8080
   ```

---

### Problem: gRPC plugin connection fails

**Symptoms**:

```
rpc error: code = Unavailable desc = connection refused
```

**Diagnostic Steps**:

1. **Check plugin is running**:

   ```bash
   ps aux | grep plugin-name
   ```

2. **Verify gRPC port**:

   ```bash
   netstat -an | grep 50051
   ```

3. **Test with grpcurl**:

   ```bash
   grpcurl -plaintext localhost:50051 list
   ```

**Solutions**:

1. **Restart plugin**:

   ```bash
   pkill plugin-name
   # Server will auto-start
   ```

2. **Check network connectivity**:

   ```bash
   ping localhost
   nc -zv localhost 50051
   ```

3. **Review plugin configuration**:

   ```bash
   cat ~/.pulumicost/plugins/config.yaml
   ```

---

## Data and Results Issues

### Problem: Cost results seem inaccurate

**Symptoms**:

- Costs don't match cloud provider bills
- Unexpected $0.00 costs
- Missing resources

**Diagnostic Steps**:

1. **Verify mock data is in use**:

   Current implementation uses mock data. Real data will be available when
   pulumicost-core integration is complete.

2. **Check filters** aren't excluding resources:

   ```json
   {
     "filters": {
       "provider": "aws"  // Will exclude Azure, GCP resources
     }
   }
   ```

3. **Verify time range** for actual costs:

   ```json
   {
     "time_range": {
       "start": "2024-01-01T00:00:00Z",
       "end": "2024-01-31T23:59:59Z"
     }
   }
   ```

**Solutions**:

1. **Use appropriate tool**:

   - `analyze_projected_cost` - For Pulumi preview data
   - `get_actual_cost` - For historical cloud bills

2. **Check data source**:

   When using plugins, verify plugin is pulling latest data.

3. **Report data issues**:

   https://github.com/rshade/pulumicost-mcp/issues

---

### Problem: Missing resources in analysis

**Symptoms**:

- Stack has 50 resources, but only 30 returned

**Diagnostic Steps**:

1. **Check for filters**:

   Remove `filters` parameter to get all resources.

2. **Verify Pulumi JSON** is complete:

   ```bash
   pulumi preview --json | jq '.resources | length'
   ```

3. **Check for pagination** (future feature):

   Currently all results returned in single response.

**Solutions**:

1. **Remove filters**:

   ```json
   {
     // Don't use filters
     "pulumi_json": "..."
   }
   ```

2. **Verify resource types are supported**:

   See [MCP Tools Reference](mcp-tools.md) for supported providers.

---

## Development and Testing

### Problem: Tests failing

**Symptoms**:

```bash
$ make test
FAIL: internal/service/cost_service_test.go:45
```

**Diagnostic Steps**:

1. **Run specific test**:

   ```bash
   go test -v -run TestAnalyzeProjected ./internal/service
   ```

2. **Check for code generation**:

   ```bash
   make generate
   # If files changed, tests may need updates
   ```

3. **Review test output**:

   ```bash
   go test -v ./... 2>&1 | tee test-output.log
   ```

**Solutions**:

1. **Update generated code**:

   ```bash
   make clean
   make generate
   make test
   ```

2. **Fix test assertions**:

   Check if generated types changed and update test expectations.

3. **Clear test cache**:

   ```bash
   go clean -testcache
   go test ./...
   ```

---

### Problem: Linter errors

**Symptoms**:

```bash
$ make lint
internal/service/cost_service.go:45:2: unused variable
```

**Solutions**:

1. **Fix reported issues**:

   Review and address each linter error.

2. **Auto-fix** where possible:

   ```bash
   golangci-lint run --fix
   ```

3. **Skip specific linters** (not recommended):

   ```yaml
   # .golangci.yml
   linters:
     disable:
       - unused  # Don't do this unless necessary
   ```

---

### Problem: Code generation fails

**Symptoms**:

```bash
$ make generate
goa gen: design error: unknown type
```

**Diagnostic Steps**:

1. **Check design syntax**:

   ```bash
   go build ./design
   ```

2. **Review error message**:

   ```bash
   make generate 2>&1 | tee gen-error.log
   ```

3. **Validate Goa DSL**:

   ```bash
   goa gen --help
   goa gen github.com/rshade/pulumicost-mcp/design -v
   ```

**Solutions**:

1. **Fix design errors**:

   Review design files for syntax errors.

2. **Update Goa version**:

   ```bash
   go get -u goa.design/goa/v3
   go mod tidy
   ```

3. **Clean and regenerate**:

   ```bash
   make clean
   rm -rf gen/
   make generate
   ```

---

## Logging and Debugging

### Enable Debug Logging

**Environment Variable**:

```bash
MCP_LOG_LEVEL=debug ./bin/pulumicost-mcp
```

**Log Levels**:

- `debug` - Verbose debugging information
- `info` - General informational messages (default)
- `warn` - Warning messages
- `error` - Error messages only

**Sample Debug Output**:

```
2024-01-08T10:30:00Z DEBUG [mcp] Received request: method=cost.analyze_projected id=1
2024-01-08T10:30:00Z DEBUG [service] Processing AnalyzeProjected request
2024-01-08T10:30:00Z DEBUG [adapter] Executing pulumicost-core with args: [--json, ...]
2024-01-08T10:30:02Z DEBUG [adapter] pulumicost-core returned 234 resources
2024-01-08T10:30:02Z DEBUG [service] Applied filters: provider=aws, region=us-east-1
2024-01-08T10:30:02Z DEBUG [mcp] Sending response: id=1 duration=2.1s
```

---

### Log File Locations

**Server Logs**:

- Default: stdout/stderr
- Custom: Set `MCP_LOG_FILE` environment variable

**Claude Desktop Logs**:

- **macOS**: `~/Library/Logs/Claude/mcp-*.log`
- **Windows**: `%APPDATA%\Claude\Logs\mcp-*.log`

**Plugin Logs**:

- `~/.pulumicost/plugins/<plugin-name>.log`

---

### Capturing Logs

**Save to File**:

```bash
./bin/pulumicost-mcp 2>&1 | tee pulumicost-mcp.log
```

**Follow Logs in Real-time**:

```bash
tail -f pulumicost-mcp.log
```

**Filter Logs**:

```bash
# Only errors
grep ERROR pulumicost-mcp.log

# Specific service
grep "cost_service" pulumicost-mcp.log

# Time range
grep "2024-01-08T10:" pulumicost-mcp.log
```

---

### Debugging with Delve

**Install Delve**:

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

**Debug Server**:

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o bin/pulumicost-mcp-debug ./cmd/pulumicost-mcp

# Run with debugger
dlv exec ./bin/pulumicost-mcp-debug

# Set breakpoints
(dlv) break internal/service/cost_service.go:45
(dlv) continue
```

**Debug Tests**:

```bash
dlv test ./internal/service -- -test.run TestAnalyzeProjected
```

---

## Getting Help

### Before Reporting Issues

1. **Check documentation**:
   - [README](../../README.md)
   - [MCP Tools Reference](mcp-tools.md)
   - [Architecture Overview](../architecture/system-overview.md)

2. **Search existing issues**:
   - https://github.com/rshade/pulumicost-mcp/issues

3. **Gather information**:
   - Server logs
   - Claude Desktop logs
   - Configuration files
   - Steps to reproduce

### Reporting Bugs

Create a new issue: https://github.com/rshade/pulumicost-mcp/issues/new

**Include**:

1. **Environment**:
   - OS and version
   - Go version: `go version`
   - Server version: `./bin/pulumicost-mcp --version`
   - Claude Desktop version

2. **Configuration**:
   - `claude_desktop_config.json` (redact sensitive data)
   - `config.yaml` (if using)
   - Environment variables

3. **Logs**:
   - Server logs (with debug enabled)
   - Claude Desktop logs
   - Error messages

4. **Steps to Reproduce**:
   - Exact commands run
   - Tool calls made
   - Expected vs actual behavior

5. **Screenshots** (if applicable)

### Community Support

- **Discussions**: https://github.com/rshade/pulumicost-mcp/discussions
- **Discord**: (Coming soon)
- **Email**: support@pulumicost.io

---

## Common Error Reference

### Error Code Quick Reference

| Error | Meaning | Common Cause | Solution |
|-------|---------|--------------|----------|
| `validation_error` | Invalid input | Missing required field | Check [MCP Tools Reference](mcp-tools.md) |
| `not_found` | Resource not found | Invalid stack/plugin name | Verify resource exists |
| `bad_request` | Malformed request | Invalid JSON/format | Validate request format |
| `internal_error` | Server error | Code bug or dependency failure | Check logs, report issue |
| `timeout` | Operation timeout | Large dataset or slow plugin | Increase timeout or use filters |
| `connection_refused` | Can't connect | Server not running | Start server |
| `permission_denied` | Access denied | Binary not executable | `chmod +x` the binary |

---

## Performance Tuning

### Optimization Checklist

- [ ] Enable caching (when available)
- [ ] Use filters to reduce data volume
- [ ] Increase timeout for large operations
- [ ] Use streaming tools for large datasets
- [ ] Monitor resource usage
- [ ] Upgrade hardware if needed

### Recommended Configuration

```yaml
# config.yaml (production)
server:
  port: 8080
  timeout: 60s
  log_level: info

cache:
  enabled: true
  ttl: 300s
  max_size: 500MB

plugins:
  timeout: 30s
  max_concurrent: 10
  health_check_interval: 60s

mcp:
  enable_streaming: true
  max_message_size: 10485760  # 10MB
```

---

## See Also

- [MCP Tools Reference](mcp-tools.md) - Tool documentation
- [Architecture Overview](../architecture/system-overview.md) - System design
- [Developer Guide](developer-guide.md) - Development setup
- [Contributing Guide](../../CONTRIBUTING.md) - How to contribute
