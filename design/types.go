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

// Granularity represents time granularity for cost queries
var Granularity = Type("Granularity", func() {
	Description("Time granularity for cost aggregation")
	Enum("hourly", "daily", "weekly", "monthly")
})

// OutputFormat represents output format options
var OutputFormat = Type("OutputFormat", func() {
	Description("Output format for results")
	Enum("table", "json", "ndjson", "csv")
})

// GroupBy represents grouping options for cost aggregation
var GroupBy = Type("GroupBy", func() {
	Description("Grouping dimension for cost aggregation")
	Enum("provider", "service", "region", "resource_type", "tag", "date")
})
