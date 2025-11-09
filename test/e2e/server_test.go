package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/rshade/pulumicost-mcp/internal/adapter"
	"github.com/rshade/pulumicost-mcp/internal/config"
	"github.com/rshade/pulumicost-mcp/internal/service"
	mcpcost "github.com/rshade/pulumicost-mcp/gen/mcp_cost"
	costsvr "github.com/rshade/pulumicost-mcp/gen/jsonrpc/mcp_cost/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goahttp "goa.design/goa/v3/http"
)

// TestMCPProtocolFlow tests the complete MCP protocol flow
func TestMCPProtocolFlow(t *testing.T) {
	// Start test server
	server, url := startTestServer(t)
	defer server.Close()

	// Test 1: Initialize
	t.Run("Initialize", func(t *testing.T) {
		req := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "initialize",
			"params": map[string]interface{}{
				"protocolVersion": "2025-06-18",
				"capabilities":    map[string]interface{}{},
				"clientInfo": map[string]interface{}{
					"name":    "test-client",
					"version": "1.0.0",
				},
			},
			"id": 1,
		}

		resp := makeRequest(t, url, req)
		require.NotNil(t, resp["result"])

		result := resp["result"].(map[string]interface{})
		assert.Equal(t, "2025-06-18", result["protocolVersion"])
		assert.NotNil(t, result["serverInfo"])
		assert.NotNil(t, result["capabilities"])
	})

	// Test 2: Ping
	t.Run("Ping", func(t *testing.T) {
		req := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "ping",
			"params":  map[string]interface{}{},
			"id":      2,
		}

		resp := makeRequest(t, url, req)
		require.NotNil(t, resp["result"])
	})

	// Test 3: Tools List
	t.Run("ToolsList", func(t *testing.T) {
		req := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "tools/list",
			"params":  map[string]interface{}{},
			"id":      3,
		}

		resp := makeRequest(t, url, req)
		require.NotNil(t, resp["result"])

		result := resp["result"].(map[string]interface{})
		tools := result["tools"].([]interface{})
		assert.Greater(t, len(tools), 0, "Should have at least one tool")

		// Verify tool structure
		tool := tools[0].(map[string]interface{})
		assert.NotEmpty(t, tool["name"])
		assert.NotNil(t, tool["inputSchema"])
	})

	// Test 4: Tools Call (analyze_projected_costs) - SSE response
	t.Run("ToolsCall_AnalyzeProjected", func(t *testing.T) {
		req := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "tools/call",
			"params": map[string]interface{}{
				"name": "analyze_projected_costs",
				"arguments": map[string]interface{}{
					"PulumiJSON": `{
						"resources": [
							{
								"urn": "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::web",
								"type": "aws:ec2/instance:Instance",
								"inputs": {"instanceType": "t3.micro"}
							}
						]
					}`,
				},
			},
			"id": 4,
		}

		resp := makeSSERequest(t, url, req)
		require.NotNil(t, resp["result"])

		result := resp["result"].(map[string]interface{})
		assert.NotNil(t, result["content"])
	})
}

// TestConcurrentRequests tests the server handles concurrent requests
func TestConcurrentRequests(t *testing.T) {
	server, url := startTestServer(t)
	defer server.Close()

	// Initialize first
	initReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2025-06-18",
			"capabilities":    map[string]interface{}{},
			"clientInfo":      map[string]interface{}{"name": "test", "version": "1.0.0"},
		},
		"id": 1,
	}
	makeRequest(t, url, initReq)

	// Make 10 concurrent ping requests
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			req := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "ping",
				"params":  map[string]interface{}{},
				"id":      id + 100,
			}
			resp := makeRequest(t, url, req)
			assert.NotNil(t, resp["result"])
			done <- true
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent requests")
		}
	}
}

// Helper functions

func startTestServer(t *testing.T) (*http.Server, string) {
	// Load config with defaults (just to validate it works)
	_, err := config.Load("")
	require.NoError(t, err)

	// Create PulumiCost adapter with mock
	mockAdapterPath := "../../internal/adapter/testdata/mock_pulumicost.sh"
	pulumiAdapter := adapter.NewPulumiCostAdapter(mockAdapterPath)

	// Create Cost service
	costService := service.NewCostService(pulumiAdapter, nil)

	// Create MCP adapter
	mcpAdapter := mcpcost.NewMCPAdapter(costService, nil)

	// Create endpoints
	endpoints := mcpcost.NewEndpoints(mcpAdapter)

	// Create server
	mux := goahttp.NewMuxer()
	server := costsvr.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil)
	costsvr.Mount(mux, server)

	// Start HTTP server on a fixed test port
	httpServer := &http.Server{
		Addr:    ":18080",
		Handler: mux,
	}

	// Start in background
	go func() {
		httpServer.ListenAndServe()
	}()

	// Give server time to start
	time.Sleep(200 * time.Millisecond)

	url := "http://localhost:18080/"

	// Register cleanup
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	})

	return httpServer, url
}

func makeRequest(t *testing.T, baseURL string, req map[string]interface{}) map[string]interface{} {
	// Marshal request
	body, err := json.Marshal(req)
	require.NoError(t, err)

	// Make HTTP request (append "rpc" path)
	resp, err := http.Post(baseURL+"rpc", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Debug: print response
	if testing.Verbose() {
		t.Logf("Response body: %s", string(respBody))
	}

	// Parse response
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	require.NoError(t, err, "Failed to parse response: %s", string(respBody))

	return result
}

func makeSSERequest(t *testing.T, baseURL string, req map[string]interface{}) map[string]interface{} {
	// Marshal request
	body, err := json.Marshal(req)
	require.NoError(t, err)

	// Make HTTP request (append "rpc" path)
	resp, err := http.Post(baseURL+"rpc", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Debug: print response
	if testing.Verbose() {
		t.Logf("Response body: %s", string(respBody))
	}

	// Parse SSE format: extract JSON from "data: " line
	lines := bytes.Split(respBody, []byte("\n"))
	var jsonData []byte
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("data: ")) {
			jsonData = bytes.TrimPrefix(line, []byte("data: "))
			break
		}
	}

	require.NotEmpty(t, jsonData, "No SSE data line found in response")

	// Parse JSON from data line
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err, "Failed to parse SSE JSON: %s", string(jsonData))

	return result
}
