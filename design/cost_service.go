package design

import (
	mcp "goa.design/goa-ai/dsl"
	. "goa.design/goa/v3/dsl"
)

var _ = Service("cost", func() {
	Description("Cost Query Service for analyzing cloud infrastructure costs")

	// Enable MCP for this service
	mcp.MCP("cost-mcp", "1.0.0")

	// JSON-RPC transport with SSE support
	JSONRPC(func() {
		POST("/rpc")
	})

	// Analyze Projected Costs
	Method("analyze_projected", func() {
		Description("Calculate projected costs from Pulumi preview JSON")
		Payload(func() {
			Attribute("pulumi_json", String, "Pulumi preview JSON output", func() {
				MinLength(1)
			})
			Attribute("filters", ResourceFilter, "Resource filtering criteria")
			Attribute("group_by", ArrayOf(String), "Group results by dimension", func() {
				Elem(func() {
					Enum("provider", "service", "region", "tag")
				})
			})
			Required("pulumi_json")
		})
		Result(CostResult)
		Error("invalid_input", ValidationError, "Invalid Pulumi JSON or parameters")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/cost/analyze_projected")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("analyze_projected_costs", "Calculate projected infrastructure costs before deployment")
	JSONRPC(func() {})
	})

	// Get Actual Costs
	Method("get_actual", func() {
		Description("Retrieve actual historical costs from cloud providers")
		Payload(func() {
			Attribute("stack_name", String, "Pulumi stack name", func() {
				MinLength(1)
			})
			Attribute("time_range", TimeRange, "Time period for cost data")
			Attribute("granularity", String, "Time aggregation granularity", func() {
				Enum("hourly", "daily", "monthly")
			})
			Attribute("filters", ResourceFilter, "Resource filtering criteria")
			Required("stack_name", "time_range")
		})
		Result(CostResult)
		Error("invalid_input", ValidationError, "Invalid stack name or time range")
		Error("not_found", NotFoundError, "Stack not found or no cost data available")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/cost/get_actual")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("not_found", StatusNotFound)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("get_actual_costs", "Retrieve historical cloud spending for deployed infrastructure")
	JSONRPC(func() {})
	})

	// Compare Costs
	Method("compare_costs", func() {
		Description("Compare costs between two configurations")
		Payload(func() {
			Attribute("baseline", func() {
				Description("Baseline configuration")
				Attribute("stack_name", String, "Stack name")
				Attribute("pulumi_json", String, "Pulumi JSON")
				Attribute("filters", ResourceFilter, "Filters")
			})
			Attribute("target", func() {
				Description("Target configuration to compare against baseline")
				Attribute("stack_name", String, "Stack name")
				Attribute("pulumi_json", String, "Pulumi JSON")
				Attribute("filters", ResourceFilter, "Filters")
			})
			Attribute("comparison_type", String, "Type of comparison", func() {
				Enum("absolute", "percentage", "both")
			})
			Required("baseline", "target")
		})
		Result(func() {
			Description("Cost comparison result")
			Attribute("baseline_cost", Float64, "Baseline total cost")
			Attribute("target_cost", Float64, "Target total cost")
			Attribute("difference", Float64, "Absolute cost difference")
			Attribute("difference_percent", Float64, "Percentage difference")
			Attribute("resource_changes", ArrayOf(Any), "Per-resource changes")
			Required("baseline_cost", "target_cost", "difference", "difference_percent")
		})
		Error("invalid_input", ValidationError, "Invalid comparison parameters")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/cost/compare")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("compare_costs", "Compare infrastructure costs between configurations")
	JSONRPC(func() {})
	})

	// Analyze Resource Cost
	Method("analyze_resource", func() {
		Description("Get detailed cost analysis for a specific resource")
		Payload(func() {
			Attribute("resource_urn", String, "Pulumi resource URN", func() {
				MinLength(1)
				Pattern("^urn:pulumi:")
			})
			Attribute("time_range", TimeRange, "Time period for analysis")
			Attribute("include_dependencies", Boolean, "Include dependent resources", func() {
				Default(true)
			})
			Required("resource_urn")
		})
		Result(func() {
			Description("Resource cost analysis")
			Attribute("resource", ResourceCost, "Primary resource cost")
			Attribute("dependencies", ArrayOf(ResourceCost), "Dependent resources")
			Attribute("trends", func() {
				Description("Cost trends")
				Attribute("daily_average", Float64, "Daily average cost")
				Attribute("weekly_average", Float64, "Weekly average cost")
				Attribute("monthly_average", Float64, "Monthly average cost")
			})
			Required("resource")
		})
		Error("invalid_input", ValidationError, "Invalid resource URN")
		Error("not_found", NotFoundError, "Resource not found")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/cost/analyze_resource")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("not_found", StatusNotFound)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("analyze_resource_cost", "Get detailed cost analysis for a specific resource")
	JSONRPC(func() {})
	})

	// Query Costs by Tags
	Method("query_by_tags", func() {
		Description("Group and aggregate costs by resource tags")
		Payload(func() {
			Attribute("stack_name", String, "Pulumi stack name", func() {
				MinLength(1)
			})
			Attribute("tag_keys", ArrayOf(String), "Tag keys to group by", func() {
				MinLength(1)
			})
			Attribute("filters", func() {
				Description("Tag filtering criteria")
				Attribute("key", String, "Tag key")
				Attribute("values", ArrayOf(String), "Acceptable tag values")
			})
			Required("stack_name", "tag_keys")
		})
		Result(func() {
			Description("Tag-based cost query result")
			Attribute("by_tag", MapOf(String, MapOf(String, Float64)), "Costs grouped by tag key/value")
			Required("by_tag")
		})
		Error("invalid_input", ValidationError, "Invalid stack or tag parameters")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/cost/query_by_tags")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("query_cost_by_tags", "Group costs by resource tags for cost attribution")
	JSONRPC(func() {})
	})

	// Analyze Stack (with streaming)
	Method("analyze_stack", func() {
		Description("Comprehensive stack cost analysis with streaming progress updates")
		Payload(func() {
			Attribute("stack_name", String, "Pulumi stack name", func() {
				MinLength(1)
			})
			Attribute("include_recommendations", Boolean, "Include optimization recommendations", func() {
				Default(false)
			})
			Required("stack_name")
		})
		StreamingResult(func() {
			Description("Streaming progress update")
			Attribute("progress", Float64, "Completion percentage (0-100)", func() {
				Minimum(0)
				Maximum(100)
			})
			Attribute("message", String, "Status message")
			Attribute("result", CostResult, "Final result (only in last message)")
		})
		Error("invalid_input", ValidationError, "Invalid stack name")
		Error("not_found", NotFoundError, "Stack not found")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/cost/analyze_stack")
			ServerSentEvents()
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("not_found", StatusNotFound)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("analyze_stack_comprehensive", "Comprehensive stack cost analysis with progress updates")
	JSONRPC(func() {
		ServerSentEvents()
	})
	})
})
