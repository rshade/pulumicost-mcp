package design

import (
	. "goa.design/goa/v3/dsl"
	mcp "goa.design/goa-ai/dsl"
)

// API defines the global properties of the PulumiCost MCP API
var _ = API("pulumicost-mcp", func() {
	Title("PulumiCost MCP Server")
	Description("AI-powered cloud cost analysis via Model Context Protocol")
	Version("1.0.0")

	Server("pulumicost-mcp", func() {
		Description("PulumiCost MCP server")

		// HTTP server for JSON-RPC
		Host("localhost", func() {
			URI("http://localhost:8080")
		})

		// MCP-specific services
		Services("cost", "plugin", "analysis")
	})

	// MCP server configuration
	mcp.MCPServer("pulumicost", "1.0.0", func() {
		Description("Cloud cost analysis and optimization for infrastructure as code")
		mcp.ServerCapability("tools")
		mcp.ServerCapability("resources")
	})
})

// Common error responses
var InternalError = Type("InternalError", func() {
	Description("Internal server error")
	Attribute("message", String, "Error message")
	Attribute("request_id", String, "Request ID for tracking")
	Required("message")
})

var NotFoundError = Type("NotFoundError", func() {
	Description("Resource not found")
	Attribute("message", String, "Error message")
	Attribute("resource", String, "Resource identifier")
	Required("message")
})

var ValidationError = Type("ValidationError", func() {
	Description("Validation error")
	Attribute("message", String, "Error message")
	Attribute("field", String, "Field that failed validation")
	Attribute("value", String, "Invalid value")
	Required("message")
})
