package design

import (
	mcp "goa.design/goa-ai/dsl"
	. "goa.design/goa/v3/dsl"
)

var _ = Service("plugin", func() {
	Description("Plugin Management Service for discovering and validating cost source plugins")

	// Enable MCP for this service
	mcp.MCP("plugin-mcp", "1.0.0")

	// JSON-RPC transport with SSE support
	JSONRPC(func() {
		POST("/rpc")
	})

	// List Plugins
	Method("list", func() {
		Description("List all available cost source plugins")
		Payload(func() {
			Attribute("include_health", Boolean, "Include health check results", func() {
				Default(true)
			})
		})
		Result(func() {
			Description("List of available plugins")
			Attribute("plugins", ArrayOf(Plugin), "Available cost source plugins")
			Required("plugins")
		})
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/plugin/list")
			Response(StatusOK)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("list_cost_plugins", "Discover available cost source plugins")
	JSONRPC(func() {})
	})

	// Get Plugin Info
	Method("get_info", func() {
		Description("Get detailed information about a specific plugin")
		Payload(func() {
			Attribute("plugin_name", String, "Plugin identifier", func() {
				MinLength(1)
			})
			Required("plugin_name")
		})
		Result(func() {
			Description("Detailed plugin information")
			Attribute("name", String, "Plugin identifier")
			Attribute("version", String, "Semantic version")
			Attribute("description", String, "Plugin description")
			Attribute("capabilities", PluginCapabilities, "Supported features")
			Attribute("health_status", HealthStatus, "Current health state")
			Attribute("grpc_address", String, "gRPC endpoint")
			Attribute("configuration", MapOf(String, Any), "Plugin configuration")
			Required("name", "version", "capabilities")
		})
		Error("invalid_input", ValidationError, "Invalid plugin name")
		Error("not_found", NotFoundError, "Plugin not found")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/plugin/get_info")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("not_found", StatusNotFound)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("get_plugin_info", "Get capabilities and configuration for a plugin")
	JSONRPC(func() {})
	})

	// Validate Plugin
	Method("validate", func() {
		Description("Validate plugin against pulumicost-spec conformance tests")
		Payload(func() {
			Attribute("plugin_path", String, "Path to plugin binary or directory", func() {
				MinLength(1)
			})
			Attribute("conformance_level", String, "Conformance level to test", func() {
				Enum("BASIC", "STANDARD", "FULL")
				Default("STANDARD")
			})
			Required("plugin_path")
		})
		Result(PluginValidationReport)
		Error("invalid_input", ValidationError, "Invalid plugin path or conformance level")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/plugin/validate")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("validate_plugin_spec", "Run conformance tests on a plugin")
	JSONRPC(func() {})
	})

	// Health Check
	Method("health_check", func() {
		Description("Check plugin health and connectivity")
		Payload(func() {
			Attribute("plugin_name", String, "Plugin identifier", func() {
				MinLength(1)
			})
			Required("plugin_name")
		})
		Result(HealthStatus)
		Error("invalid_input", ValidationError, "Invalid plugin name")
		Error("not_found", NotFoundError, "Plugin not found")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/plugin/health_check")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("not_found", StatusNotFound)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("check_plugin_health", "Verify plugin connectivity and response time")
	JSONRPC(func() {})
	})
})
