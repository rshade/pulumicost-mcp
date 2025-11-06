# Data Model: PulumiCost MCP Server

**Feature**: 001-mcp-server
**Created**: 2025-01-06
**Status**: Design Phase

## Overview

This document defines the data entities and their relationships for the
PulumiCost MCP Server. All types will be defined in Goa DSL (`design/types.go`)
and generated code will provide type-safe marshaling/unmarshaling.

## Core Entities

### CostQuery

**Purpose**: Represents a request for cost analysis with filters and grouping
criteria

**Fields**:

- `stack_name` (string, optional): Pulumi stack identifier
- `pulumi_json` (string, optional): Pulumi preview JSON for projected cost
  analysis
- `resource_filters` (ResourceFilter[], optional): Filter resources by provider,
  type, region
- `tag_filters` (TagFilter[], optional): Filter by resource tags
- `time_range` (TimeRange, optional): Period for actual cost queries
- `granularity` (string, optional): Time bucket size (hourly, daily, monthly)
- `group_by` (string[], optional): Group results by provider, service, region,
  or tag

**Validation Rules**:

- Either `pulumi_json` OR `stack_name` must be provided
- `time_range` required for actual cost queries
- `granularity` must be one of: hourly, daily, monthly
- `group_by` values limited to: provider, service, region, tag

**State Transitions**: N/A (immutable request)

### CostResult

**Purpose**: Contains aggregated cost data with breakdowns

**Fields**:

- `total_monthly` (float64): Total estimated monthly cost
- `currency` (string): Currency code (default: USD)
- `resources` (ResourceCost[]): Per-resource cost breakdown
- `by_provider` (map[string]float64): Costs grouped by cloud provider
- `by_service` (map[string]float64): Costs grouped by service type
- `by_region` (map[string]float64): Costs grouped by region
- `by_tag` (map[string]map[string]float64): Costs grouped by tag key/value
- `timestamp` (string): ISO 8601 timestamp of analysis
- `metadata` (CostMetadata): Additional context

**Validation Rules**:

- `total_monthly` >= 0
- `currency` must be valid ISO 4217 code
- `timestamp` in ISO 8601 format

**Relationships**:

- Contains multiple `ResourceCost` entities
- Contains `CostMetadata`

### ResourceCost

**Purpose**: Cost breakdown for a single resource

**Fields**:

- `urn` (string): Pulumi resource URN
- `type` (string): Resource type (e.g., aws:ec2/instance:Instance)
- `name` (string): Resource name
- `provider` (string): Cloud provider (aws, azure, gcp)
- `service` (string): Service category (compute, storage, network)
- `region` (string): Deployment region
- `monthly_cost` (float64): Estimated monthly cost
- `hourly_cost` (float64): Estimated hourly cost
- `tags` (map[string]string): Resource tags
- `dependencies` (string[]): URNs of dependent resources

**Validation Rules**:

- `monthly_cost` >= 0
- `hourly_cost` >= 0
- `urn` matches Pulumi URN format
- `provider` one of: aws, azure, gcp, kubernetes

**Relationships**:

- Part of `CostResult`
- May reference other `ResourceCost` via `dependencies`

### Plugin

**Purpose**: Represents a cost source plugin with metadata

**Fields**:

- `name` (string): Plugin identifier
- `version` (string): Semantic version
- `description` (string): Plugin purpose
- `capabilities` (PluginCapabilities): Supported features
- `health_status` (HealthStatus): Current health state
- `grpc_address` (string): gRPC endpoint
- `metadata` (map[string]string): Additional plugin info

**Validation Rules**:

- `name` non-empty, alphanumeric with hyphens
- `version` follows semver format (vX.Y.Z)
- `grpc_address` valid host:port format

**State Transitions**:

- UNKNOWN → HEALTHY (on successful health check)
- HEALTHY → UNHEALTHY (on failed health check)
- UNHEALTHY → HEALTHY (on recovery)
- ANY → UNKNOWN (on connection loss)

### PluginCapabilities

**Purpose**: Describes what a plugin can do

**Fields**:

- `supports_projected` (bool): Can calculate projected costs
- `supports_actual` (bool): Can retrieve actual historical costs
- `supports_providers` (string[]): Supported cloud providers
- `supports_resources` (string[]): Supported resource types

**Validation Rules**:

- At least one of `supports_projected` or `supports_actual` must be true
- `supports_providers` non-empty

### HealthStatus

**Purpose**: Plugin health information

**Fields**:

- `status` (string): One of: HEALTHY, UNHEALTHY, UNKNOWN
- `last_check` (string): ISO 8601 timestamp
- `latency_ms` (int64): Response time in milliseconds
- `error_message` (string, optional): Last error if unhealthy

**Validation Rules**:

- `status` must be one of: HEALTHY, UNHEALTHY, UNKNOWN
- `latency_ms` >= 0

### PluginValidationReport

**Purpose**: Results from conformance testing

**Fields**:

- `plugin_name` (string): Plugin being validated
- `conformance_level` (string): BASIC, STANDARD, FULL
- `passed` (bool): Overall pass/fail
- `test_results` (ValidationTest[]): Individual test outcomes
- `timestamp` (string): ISO 8601 timestamp of validation
- `spec_version` (string): pulumicost-spec version used

**Validation Rules**:

- `conformance_level` one of: BASIC, STANDARD, FULL
- `timestamp` ISO 8601 format

**Relationships**:

- Contains multiple `ValidationTest` entities

### ValidationTest

**Purpose**: Single conformance test result

**Fields**:

- `name` (string): Test identifier
- `passed` (bool): Test outcome
- `error_message` (string, optional): Failure reason
- `duration_ms` (int64): Test execution time

**Validation Rules**:

- `duration_ms` >= 0

### Recommendation

**Purpose**: Cost optimization suggestion

**Fields**:

- `id` (string): Unique recommendation ID
- `type` (string): RIGHTSIZING, RESERVED_INSTANCES, SPOT_INSTANCES, etc.
- `resource_urn` (string): Affected resource
- `current_cost` (float64): Current monthly cost
- `projected_savings` (float64): Estimated monthly savings
- `confidence` (string): LOW, MEDIUM, HIGH
- `description` (string): Human-readable explanation
- `action_steps` (string[]): Implementation guidance

**Validation Rules**:

- `current_cost` >= 0
- `projected_savings` >= 0
- `confidence` one of: LOW, MEDIUM, HIGH
- `type` one of defined recommendation types

### Anomaly

**Purpose**: Detected cost irregularity

**Fields**:

- `id` (string): Unique anomaly ID
- `timestamp` (string): ISO 8601 when detected
- `resource_urns` (string[]): Affected resources
- `severity` (string): LOW, MEDIUM, HIGH, CRITICAL
- `current_cost` (float64): Current spending
- `baseline_cost` (float64): Expected spending
- `deviation_percent` (float64): Percentage difference
- `potential_causes` (string[]): Possible explanations

**Validation Rules**:

- `severity` one of: LOW, MEDIUM, HIGH, CRITICAL
- `deviation_percent` can be negative or positive
- `timestamp` ISO 8601 format

### Forecast

**Purpose**: Projected future costs

**Fields**:

- `stack_name` (string): Stack being forecasted
- `forecast_period` (TimeRange): Future time range
- `data_points` (ForecastPoint[]): Time-series predictions
- `confidence_level` (float64): 0.0-1.0 confidence
- `methodology` (string): Forecasting approach used
- `generated_at` (string): ISO 8601 timestamp

**Validation Rules**:

- `confidence_level` between 0.0 and 1.0
- `data_points` non-empty
- `generated_at` ISO 8601 format

**Relationships**:

- Contains multiple `ForecastPoint` entities

### ForecastPoint

**Purpose**: Single point in forecast time series

**Fields**:

- `timestamp` (string): ISO 8601 date
- `predicted_cost` (float64): Expected cost
- `lower_bound` (float64): Confidence interval lower
- `upper_bound` (float64): Confidence interval upper

**Validation Rules**:

- `lower_bound` <= `predicted_cost` <= `upper_bound`
- All cost values >= 0

### Budget

**Purpose**: Budget definition and tracking

**Fields**:

- `id` (string): Unique budget ID
- `stack_name` (string): Stack this budget applies to
- `amount` (float64): Budget amount
- `currency` (string): Currency code
- `period` (string): DAILY, WEEKLY, MONTHLY, YEARLY
- `alert_thresholds` (float64[]): Percentage thresholds (e.g., [50, 80, 100])
- `current_spending` (float64): Actual spending so far
- `remaining` (float64): Budget remaining
- `burn_rate` (float64): Daily spending rate
- `projected_end_date` (string): ISO 8601 when budget exhausts
- `status` (string): OK, WARNING, EXCEEDED

**Validation Rules**:

- `amount` > 0
- `alert_thresholds` sorted ascending, values 0-100
- `current_spending` >= 0
- `period` one of: DAILY, WEEKLY, MONTHLY, YEARLY
- `status` one of: OK, WARNING, EXCEEDED

**State Transitions**:

- OK → WARNING (current_spending >= first threshold)
- WARNING → EXCEEDED (current_spending > amount)
- Can reset to OK at period boundary

## Supporting Types

### ResourceFilter

**Fields**:

- `provider` (string, optional): Filter by provider
- `resource_type` (string, optional): Filter by type
- `region` (string, optional): Filter by region

### TagFilter

**Fields**:

- `key` (string): Tag key to filter on
- `values` (string[]): Acceptable tag values

### TimeRange

**Fields**:

- `start` (string): ISO 8601 start time
- `end` (string): ISO 8601 end time

**Validation Rules**:

- `start` < `end`

### CostMetadata

**Fields**:

- `source` (string): Data source (projected, actual, estimated)
- `confidence` (string): Data quality (high, medium, low)
- `notes` (string[]): Additional context or warnings

## Entity Relationships

```text
CostResult
  ├── contains many ResourceCost
  └── contains one CostMetadata

Plugin
  ├── has one PluginCapabilities
  └── has one HealthStatus

PluginValidationReport
  └── contains many ValidationTest

Forecast
  └── contains many ForecastPoint

ResourceCost
  └── references other ResourceCost via dependencies
```

## AI Methodology

### Algorithmic Approach

The "AI-powered" features in PulumiCost MCP Server use statistical analysis,
pattern recognition, and time-series forecasting techniques. These are
algorithmic approaches rather than machine learning models, optimized for:

- **Deterministic results**: Same input always produces same output
- **Explainability**: Clear reasoning for recommendations
- **Low latency**: Sub-second analysis without model inference
- **No training data required**: Works immediately on new stacks

### Recommendation Engine

**Approach**: Rule-based pattern analysis with heuristics

**Methodology**:

1. **Rightsizing**: Compare resource utilization metrics (if available from
   cloud provider APIs) against provisioned capacity. Flag resources with <50%
   utilization consistently for 30+ days.
2. **Reserved Instances**: Identify long-running resources (>90 days uptime)
   with stable usage patterns. Calculate savings from reserved capacity.
3. **Spot Instances**: Detect stateless, fault-tolerant workloads (containers,
   batch jobs) suitable for spot/preemptible instances.
4. **Storage Optimization**: Find infrequently accessed storage (last access
   >90 days) for lifecycle policies or archival tiers.

**Confidence Levels**:

- HIGH: Strong signal (e.g., 95%+ utilization match, 180+ day stable pattern)
- MEDIUM: Moderate signal (70-95% match, 90-180 day pattern)
- LOW: Weak signal (<70% match, <90 day pattern)

### Anomaly Detection

**Approach**: Statistical baseline comparison with standard deviation thresholds

**Methodology**:

1. **Baseline Calculation**: Compute 30-day rolling average and standard
   deviation for each resource's daily cost
2. **Deviation Detection**: Flag costs exceeding baseline ± (N × stddev) where
   N is sensitivity (LOW=3, MEDIUM=2, HIGH=1)
3. **Severity Assignment**:
   - CRITICAL: >3σ deviation, cost increase >100%
   - HIGH: 2-3σ deviation, cost increase 50-100%
   - MEDIUM: 1-2σ deviation, cost increase 25-50%
   - LOW: <1σ deviation, cost increase <25%

4. **Root Cause Inference**: Pattern matching against known causes (scaling
   events, new resources, region changes, pricing updates)

### Cost Forecasting

**Approach**: Time-series analysis with linear regression and confidence
intervals

**Methodology**:

1. **Trend Analysis**: Fit linear regression to 90-day historical cost data
2. **Seasonality Detection**: Identify weekly/monthly patterns using moving
   averages
3. **Confidence Intervals**: Calculate prediction intervals using regression
   standard error
4. **Extrapolation**: Project costs forward with uncertainty bounds widening
   over time

**Confidence Levels**:

- 0.95 (default): 95% confidence interval (2σ)
- 0.90: 90% confidence interval (1.65σ)
- 0.99: 99% confidence interval (2.58σ)

**Limitations**:

- Assumes past patterns continue (no major architecture changes)
- Accuracy decreases for forecasts >90 days
- Requires minimum 30 days historical data

### Future Enhancements

Potential ML integration points (not in MVP):

- Neural network forecasting for complex seasonal patterns
- Anomaly detection with autoencoders for multivariate analysis
- NLP-based recommendation explanations using LLMs
- Transfer learning from similar infrastructure patterns


## Type Mapping to Goa DSL

All entities will be defined in `design/types.go` using Goa DSL:

```go
var CostQuery = Type("CostQuery", func() {
    Attribute("stack_name", String)
    Attribute("pulumi_json", String)
    Attribute("resource_filters", ArrayOf(ResourceFilter))
    Attribute("tag_filters", ArrayOf(TagFilter))
    Attribute("time_range", TimeRange)
    Attribute("granularity", String)
    Attribute("group_by", ArrayOf(String))
})
```

Generated Go structs will provide JSON marshaling and validation.
