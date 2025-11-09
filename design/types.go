package design

import (
	. "goa.design/goa/v3/dsl"
)

// Common types used across services

// TimeRange defines a time period for cost queries
var TimeRange = Type("TimeRange", func() {
	Description("Time range for cost analysis")
	Attribute("start", String, "Start time (ISO 8601)", func() {
		Format(FormatDateTime)
		Example("2025-01-01T00:00:00Z")
	})
	Attribute("end", String, "End time (ISO 8601)", func() {
		Format(FormatDateTime)
		Example("2025-01-31T23:59:59Z")
	})
	Required("start", "end")
})

// ResourceFilter defines filtering criteria for resources
var ResourceFilter = Type("ResourceFilter", func() {
	Description("Filtering criteria for resources")
	Attribute("provider", String, "Cloud provider (aws, azure, gcp, kubernetes)", func() {
		Enum("aws", "azure", "gcp", "kubernetes", "custom")
	})
	Attribute("resource_type", String, "Resource type (e.g., ec2/instance, s3/bucket)")
	Attribute("region", String, "Cloud region")
	Attribute("tags", MapOf(String, String), "Resource tags to filter by")
	Attribute("name_pattern", String, "Resource name pattern (regex)")
})

// ResourceCost represents cost information for a single resource
var ResourceCost = Type("ResourceCost", func() {
	Description("Cost information for a resource")
	Attribute("urn", String, "Pulumi resource URN", func() {
		Example("urn:pulumi:prod::myapp::aws:ec2/instance:Instance::web-server")
	})
	Attribute("name", String, "Resource name")
	Attribute("type", String, "Resource type")
	Attribute("provider", String, "Cloud provider")
	Attribute("monthly_cost", Float64, "Estimated monthly cost")
	Attribute("hourly_cost", Float64, "Hourly cost rate")
	Attribute("currency", String, "Currency code (ISO 4217)", func() {
		Default("USD")
		Example("USD")
	})
	Attribute("billing_mode", String, "Billing mode", func() {
		Example("per_hour")
	})
	Attribute("adapter", String, "Pricing adapter used")
	Attribute("notes", String, "Additional cost notes")
	Attribute("tags", MapOf(String, String), "Resource tags")
	Attribute("properties", MapOf(String, Any), "Resource-specific properties")
	Required("urn", "name", "type", "monthly_cost", "currency")
})

// CostBreakdown represents aggregated cost information
var CostBreakdown = Type("CostBreakdown", func() {
	Description("Aggregated cost breakdown")
	Attribute("by_provider", MapOf(String, Float64), "Cost grouped by provider")
	Attribute("by_service", MapOf(String, Float64), "Cost grouped by service")
	Attribute("by_region", MapOf(String, Float64), "Cost grouped by region")
	Attribute("by_tag", MapOf(String, MapOf(String, Float64)), "Cost grouped by tag")
})

// ActualCostDataPoint represents a historical cost data point
var ActualCostDataPoint = Type("ActualCostDataPoint", func() {
	Description("Historical cost data point")
	Attribute("timestamp", String, "Timestamp (ISO 8601)", func() {
		Format(FormatDateTime)
	})
	Attribute("cost", Float64, "Cost amount")
	Attribute("usage_amount", Float64, "Usage amount")
	Attribute("usage_unit", String, "Usage unit")
	Attribute("currency", String, "Currency code")
	Attribute("tags", MapOf(String, String), "Cost tags")
	Required("timestamp", "cost", "currency")
})

// OptimizationRecommendation represents a cost optimization suggestion
var OptimizationRecommendation = Type("OptimizationRecommendation", func() {
	Description("Cost optimization recommendation")
	Attribute("id", String, "Recommendation ID")
	Attribute("title", String, "Recommendation title")
	Attribute("description", String, "Detailed description")
	Attribute("category", String, "Category", func() {
		Enum("rightsizing", "reserved_instance", "spot_instance", "storage_optimization", "cleanup", "scheduling", "other")
	})
	Attribute("impact", String, "Impact level", func() {
		Enum("low", "medium", "high")
	})
	Attribute("estimated_monthly_savings", Float64, "Estimated monthly savings")
	Attribute("currency", String, "Currency code")
	Attribute("resources", ArrayOf(String), "Affected resource URNs")
	Attribute("risk_level", String, "Risk level", func() {
		Enum("low", "medium", "high")
		Default("low")
	})
	Attribute("implementation_effort", String, "Implementation effort", func() {
		Enum("low", "medium", "high")
	})
	Attribute("action_items", ArrayOf(String), "Steps to implement")
	Required("id", "title", "description", "category", "impact", "estimated_monthly_savings", "currency")
})

// PluginInfo represents information about a cost source plugin
var PluginInfo = Type("PluginInfo", func() {
	Description("Cost source plugin information")
	Attribute("name", String, "Plugin name")
	Attribute("version", String, "Plugin version")
	Attribute("description", String, "Plugin description")
	Attribute("provider", String, "Cloud provider")
	Attribute("capabilities", ArrayOf(String), "Supported capabilities")
	Attribute("status", String, "Plugin status", func() {
		Enum("active", "inactive", "error")
	})
	Attribute("health", String, "Health status", func() {
		Enum("healthy", "degraded", "unhealthy")
	})
	Required("name", "version", "status")
})

// ValidationResult represents plugin validation results
var ValidationResult = Type("ValidationResult", func() {
	Description("Plugin validation results")
	Attribute("passed", Boolean, "Whether validation passed")
	Attribute("level", String, "Conformance level tested", func() {
		Enum("basic", "standard", "advanced")
	})
	Attribute("total_tests", Int, "Total number of tests")
	Attribute("passed_tests", Int, "Number of passed tests")
	Attribute("failed_tests", Int, "Number of failed tests")
	Attribute("errors", ArrayOf(String), "Validation errors")
	Attribute("warnings", ArrayOf(String), "Validation warnings")
	Attribute("report", String, "Full validation report")
	Required("passed", "level", "total_tests", "passed_tests", "failed_tests")
})

// PricingSpec represents detailed pricing specification
var PricingSpec = Type("PricingSpec", func() {
	Description("Detailed pricing specification from pulumicost-spec")
	Attribute("provider", String, "Cloud provider")
	Attribute("resource_type", String, "Resource type")
	Attribute("sku", String, "SKU identifier")
	Attribute("region", String, "Region")
	Attribute("billing_mode", String, "Billing mode")
	Attribute("rate_per_unit", Float64, "Rate per unit")
	Attribute("currency", String, "Currency")
	Attribute("description", String, "Description")
	Attribute("source", String, "Source plugin")
	Required("provider", "resource_type", "billing_mode", "rate_per_unit", "currency")
})

// Note: Granularity, OutputFormat, and GroupBy are defined inline in service methods as String with Enum constraints

// ====================
// Core Cost Types
// ====================

// CostMetadata provides additional context about cost data
var CostMetadata = Type("CostMetadata", func() {
	Description("Additional context and metadata for cost results")
	Attribute("source", String, "Data source", func() {
		Enum("projected", "actual", "estimated")
		Example("projected")
	})
	Attribute("confidence", String, "Data quality confidence", func() {
		Enum("high", "medium", "low")
		Default("medium")
	})
	Attribute("notes", ArrayOf(String), "Additional context or warnings")
	Attribute("generated_at", String, "Timestamp of analysis", func() {
		Format(FormatDateTime)
	})
})

// CostResult represents the response from cost analysis
var CostResult = Type("CostResult", func() {
	Description("Cost analysis result with breakdown")
	Attribute("total_monthly", Float64, "Total estimated monthly cost", func() {
		Minimum(0)
	})
	Attribute("currency", String, "Currency code (ISO 4217)", func() {
		Default("USD")
		Example("USD")
	})
	Attribute("resources", ArrayOf(ResourceCost), "Per-resource cost breakdown")
	Attribute("by_provider", MapOf(String, Float64), "Costs grouped by cloud provider")
	Attribute("by_service", MapOf(String, Float64), "Costs grouped by service type")
	Attribute("by_region", MapOf(String, Float64), "Costs grouped by region")
	Attribute("timestamp", String, "ISO 8601 timestamp of analysis", func() {
		Format(FormatDateTime)
	})
	Attribute("metadata", CostMetadata, "Additional context")
	Required("total_monthly", "currency", "resources")
})

// ====================
// Plugin Types
// ====================

// PluginCapabilities describes what a plugin can do
var PluginCapabilities = Type("PluginCapabilities", func() {
	Description("Capabilities and features supported by a plugin")
	Attribute("supports_projected", Boolean, "Can calculate projected costs", func() {
		Default(false)
	})
	Attribute("supports_actual", Boolean, "Can retrieve actual historical costs", func() {
		Default(false)
	})
	Attribute("supports_providers", ArrayOf(String), "Supported cloud providers")
	Required("supports_projected", "supports_actual", "supports_providers")
})

// HealthStatus represents plugin health information
var HealthStatus = Type("HealthStatus", func() {
	Description("Health status of a plugin")
	Attribute("status", String, "Health status", func() {
		Enum("HEALTHY", "UNHEALTHY", "UNKNOWN")
		Default("UNKNOWN")
	})
	Attribute("last_check", String, "ISO 8601 timestamp of last health check", func() {
		Format(FormatDateTime)
	})
	Attribute("latency_ms", Int64, "Response time in milliseconds", func() {
		Minimum(0)
	})
	Attribute("error_message", String, "Last error if unhealthy")
	Required("status")
})

// Plugin represents a cost source plugin with metadata
var Plugin = Type("Plugin", func() {
	Description("Cost source plugin information")
	Attribute("name", String, "Plugin identifier")
	Attribute("version", String, "Semantic version (vX.Y.Z)")
	Attribute("description", String, "Plugin purpose")
	Attribute("capabilities", PluginCapabilities, "Supported features")
	Attribute("health_status", HealthStatus, "Current health state")
	Required("name", "version", "capabilities")
})

// ValidationTest represents a single conformance test result
var ValidationTest = Type("ValidationTest", func() {
	Description("Single conformance test result")
	Attribute("name", String, "Test identifier")
	Attribute("passed", Boolean, "Test outcome")
	Attribute("error_message", String, "Failure reason if not passed")
	Attribute("duration_ms", Int64, "Test execution time in milliseconds", func() {
		Minimum(0)
	})
	Required("name", "passed")
})

// PluginValidationReport represents results from conformance testing
var PluginValidationReport = Type("PluginValidationReport", func() {
	Description("Results from plugin conformance testing")
	Attribute("plugin_name", String, "Plugin being validated")
	Attribute("conformance_level", String, "Conformance level tested", func() {
		Enum("BASIC", "STANDARD", "FULL")
		Default("STANDARD")
	})
	Attribute("passed", Boolean, "Overall pass/fail")
	Attribute("test_results", ArrayOf(ValidationTest), "Individual test outcomes")
	Attribute("timestamp", String, "ISO 8601 timestamp of validation", func() {
		Format(FormatDateTime)
	})
	Required("plugin_name", "conformance_level", "passed", "test_results")
})

// ====================
// Analysis Types
// ====================

// Recommendation represents a cost optimization suggestion
var Recommendation = Type("Recommendation", func() {
	Description("AI-powered cost optimization recommendation")
	Attribute("id", String, "Unique recommendation ID")
	Attribute("type", String, "Recommendation type", func() {
		Enum("RIGHTSIZING", "RESERVED_INSTANCES", "SPOT_INSTANCES", "STORAGE_OPTIMIZATION")
	})
	Attribute("resource_urn", String, "Affected resource URN")
	Attribute("current_cost", Float64, "Current monthly cost", func() {
		Minimum(0)
	})
	Attribute("projected_savings", Float64, "Estimated monthly savings", func() {
		Minimum(0)
	})
	Attribute("confidence", String, "Confidence level", func() {
		Enum("LOW", "MEDIUM", "HIGH")
		Default("MEDIUM")
	})
	Attribute("description", String, "Human-readable explanation")
	Attribute("action_steps", ArrayOf(String), "Implementation guidance")
	Required("id", "type", "resource_urn", "current_cost", "projected_savings", "confidence", "description")
})

// Anomaly represents a detected cost irregularity
var Anomaly = Type("Anomaly", func() {
	Description("Detected cost anomaly")
	Attribute("id", String, "Unique anomaly ID")
	Attribute("timestamp", String, "ISO 8601 when detected", func() {
		Format(FormatDateTime)
	})
	Attribute("resource_urns", ArrayOf(String), "Affected resources")
	Attribute("severity", String, "Anomaly severity", func() {
		Enum("LOW", "MEDIUM", "HIGH", "CRITICAL")
	})
	Attribute("current_cost", Float64, "Current spending")
	Attribute("baseline_cost", Float64, "Expected spending")
	Attribute("deviation_percent", Float64, "Percentage difference (can be negative)")
	Attribute("potential_causes", ArrayOf(String), "Possible explanations")
	Required("id", "timestamp", "resource_urns", "severity", "current_cost", "baseline_cost", "deviation_percent")
})

// ForecastPoint represents a single point in forecast time series
var ForecastPoint = Type("ForecastPoint", func() {
	Description("Single point in cost forecast")
	Attribute("timestamp", String, "ISO 8601 date", func() {
		Format(FormatDateTime)
	})
	Attribute("predicted_cost", Float64, "Expected cost", func() {
		Minimum(0)
	})
	Attribute("lower_bound", Float64, "Confidence interval lower", func() {
		Minimum(0)
	})
	Attribute("upper_bound", Float64, "Confidence interval upper", func() {
		Minimum(0)
	})
	Required("timestamp", "predicted_cost", "lower_bound", "upper_bound")
})

// Forecast represents projected future costs
var Forecast = Type("Forecast", func() {
	Description("Cost forecast with confidence intervals")
	Attribute("stack_name", String, "Stack being forecasted")
	Attribute("forecast_period", TimeRange, "Future time range")
	Attribute("data_points", ArrayOf(ForecastPoint), "Time-series predictions")
	Attribute("confidence_level", Float64, "Confidence level (0.0-1.0)", func() {
		Minimum(0.0)
		Maximum(1.0)
		Default(0.95)
	})
	Attribute("methodology", String, "Forecasting approach used")
	Required("stack_name", "forecast_period", "data_points", "confidence_level", "methodology")
})

// Budget represents budget definition and tracking
var Budget = Type("Budget", func() {
	Description("Budget configuration and status")
	Attribute("budget_amount", Float64, "Budget amount", func() {
		Minimum(0)
	})
	Attribute("current_spending", Float64, "Actual spending so far", func() {
		Minimum(0)
	})
	Attribute("remaining", Float64, "Budget remaining")
	Attribute("burn_rate", Float64, "Daily spending rate")
	Attribute("projected_end_date", String, "ISO 8601 when budget exhausts", func() {
		Format(FormatDateTime)
	})
	Attribute("status", String, "Budget status", func() {
		Enum("OK", "WARNING", "EXCEEDED")
	})
	Attribute("alerts", ArrayOf(Any), "Threshold alerts") // Array of alert objects
	Required("budget_amount", "current_spending", "remaining", "status")
})
