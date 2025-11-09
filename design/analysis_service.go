package design

import (
	_ "goa.design/goa-ai" // Import to register the MCP plugin
	mcp "goa.design/goa-ai/dsl"
	. "goa.design/goa/v3/dsl"
)

var _ = Service("analysis", func() {
	Description("Analysis Service for AI-powered cost optimization and budget tracking")

	// Enable MCP for this service
	mcp.MCPServer("analysis-mcp", "1.0.0")

	// JSON-RPC transport with SSE support
	JSONRPC(func() {
		POST("/rpc")
	})

	// Get Recommendations
	Method("get_recommendations", func() {
		Description("Get AI-powered cost optimization recommendations")
		Payload(func() {
			Attribute("stack_name", String, "Pulumi stack name", func() {
				MinLength(1)
			})
			Attribute("recommendation_types", ArrayOf(String), "Filter by recommendation types", func() {
				Elem(func() {
					Enum("RIGHTSIZING", "RESERVED_INSTANCES", "SPOT_INSTANCES", "STORAGE_OPTIMIZATION")
				})
			})
			Attribute("minimum_savings", Float64, "Minimum monthly savings threshold", func() {
				Minimum(0)
			})
			Required("stack_name")
		})
		Result(func() {
			Description("Cost optimization recommendations")
			Attribute("recommendations", ArrayOf(Recommendation), "List of recommendations")
			Required("recommendations")
		})
		Error("invalid_input", ValidationError, "Invalid stack name or parameters")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/analysis/get_recommendations")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("get_optimization_recommendations", "Get AI-powered cost optimization recommendations")
	JSONRPC(func() {})
	})

	// Detect Anomalies
	Method("detect_anomalies", func() {
		Description("Detect unusual spending patterns and cost anomalies")
		Payload(func() {
			Attribute("stack_name", String, "Pulumi stack name", func() {
				MinLength(1)
			})
			Attribute("time_range", TimeRange, "Time period to analyze")
			Attribute("sensitivity", String, "Detection sensitivity", func() {
				Enum("LOW", "MEDIUM", "HIGH")
				Default("MEDIUM")
			})
			Required("stack_name", "time_range")
		})
		Result(func() {
			Description("Detected cost anomalies")
			Attribute("anomalies", ArrayOf(Anomaly), "List of detected anomalies")
			Required("anomalies")
		})
		Error("invalid_input", ValidationError, "Invalid stack name or time range")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/analysis/detect_anomalies")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("detect_cost_anomalies", "Identify unusual spending patterns")
	JSONRPC(func() {})
	})

	// Forecast Costs
	Method("forecast", func() {
		Description("Generate cost forecast based on historical trends")
		Payload(func() {
			Attribute("stack_name", String, "Pulumi stack name", func() {
				MinLength(1)
			})
			Attribute("forecast_period", TimeRange, "Future time range to forecast")
			Attribute("confidence_level", Float64, "Confidence level (0.0-1.0)", func() {
				Minimum(0.0)
				Maximum(1.0)
				Default(0.95)
			})
			Required("stack_name", "forecast_period")
		})
		Result(Forecast)
		Error("invalid_input", ValidationError, "Invalid stack name or forecast period")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/analysis/forecast")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("forecast_costs", "Project future costs based on trends")
	JSONRPC(func() {})
	})

	// Track Budget
	Method("track_budget", func() {
		Description("Monitor spending against budget with alerts")
		Payload(func() {
			Attribute("stack_name", String, "Pulumi stack name", func() {
				MinLength(1)
			})
			Attribute("budget_amount", Float64, "Budget amount", func() {
				Minimum(0)
			})
			Attribute("period", String, "Budget period", func() {
				Enum("DAILY", "WEEKLY", "MONTHLY", "YEARLY")
			})
			Attribute("alert_thresholds", ArrayOf(Float64), "Alert threshold percentages", func() {
				Elem(func() {
					Minimum(0)
					Maximum(100)
				})
			})
			Required("stack_name", "budget_amount", "period")
		})
		Result(Budget)
		Error("invalid_input", ValidationError, "Invalid budget parameters")
		Error("internal_error", InternalError, "Internal server error")

		HTTP(func() {
			POST("/analysis/track_budget")
			Response(StatusOK)
			Response("invalid_input", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})

	mcp.Tool("track_budget", "Monitor spending against budget with alerts")
	JSONRPC(func() {})
	})
})
