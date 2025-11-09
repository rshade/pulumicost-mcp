# MCP Tools Reference

Complete reference for all 14 MCP tools provided by PulumiCost MCP Server.

## Table of Contents

- [Cost Query Tools](#cost-query-tools)
  - [analyze_projected_cost](#analyze_projected_cost)
  - [get_actual_cost](#get_actual_cost)
  - [compare_costs](#compare_costs)
  - [analyze_resource_cost](#analyze_resource_cost)
  - [query_cost_by_tags](#query_cost_by_tags)
  - [analyze_stack](#analyze_stack)
- [Plugin Management Tools](#plugin-management-tools)
  - [list_plugins](#list_plugins)
  - [get_plugin_info](#get_plugin_info)
  - [validate_plugin](#validate_plugin)
  - [health_check](#health_check)
- [Analysis and Optimization Tools](#analysis-and-optimization-tools)
  - [get_recommendations](#get_recommendations)
  - [detect_anomalies](#detect_anomalies)
  - [forecast_costs](#forecast_costs)
  - [track_budget](#track_budget)

## Cost Query Tools

### analyze_projected_cost

Calculate estimated monthly costs before deploying infrastructure.

**Description**: Analyzes Pulumi preview data to estimate monthly costs for
infrastructure resources before they are deployed.

**Use Cases**:

- Pre-deployment cost validation
- Comparing infrastructure configuration options
- Budget planning for new projects

**Input Parameters**:

```json
{
  "pulumi_json": "string (required) - Pulumi stack preview JSON",
  "filters": {
    "provider": "string (optional) - Filter by cloud provider (aws, azure, gcp)",
    "region": "string (optional) - Filter by region",
    "resource_type": "string (optional) - Filter by resource type"
  },
  "grouping": {
    "level": "string (optional) - Grouping level (resource, type, provider)",
    "tag_key": "string (optional) - Tag key for grouping"
  }
}
```

**Output**:

```json
{
  "total_monthly": 1234.56,
  "currency": "USD",
  "resources": [
    {
      "urn": "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::web-server",
      "name": "web-server",
      "type": "aws:ec2/instance:Instance",
      "provider": "aws",
      "region": "us-east-1",
      "monthly_cost": 234.50,
      "cost_components": [
        {
          "name": "Instance usage",
          "unit": "hours",
          "monthly_quantity": 730.0,
          "price": 0.0416,
          "monthly_cost": 30.37
        }
      ]
    }
  ],
  "summary": {
    "by_provider": {
      "aws": 1234.56
    },
    "by_type": {
      "ec2/instance": 234.50,
      "rds/instance": 500.00
    }
  }
}
```

**Example Usage**:

```
User: What will my staging environment cost per month?

Claude: Let me analyze your Pulumi stack preview...
[Calls analyze_projected_cost with stack preview JSON]

Based on your infrastructure configuration:
- Total Monthly Cost: $1,234.56 USD
- AWS EC2 instances: $234.50
- AWS RDS database: $500.00
- AWS S3 storage: $50.06
- AWS ALB: $450.00
```

**Error Handling**:

- Returns `ValidationError` if pulumi_json is invalid
- Returns `BadRequestError` if JSON cannot be parsed
- Returns `InternalError` for cost calculation failures

---

### get_actual_cost

Retrieve historical spending data with detailed breakdowns.

**Description**: Query actual costs from cloud providers or cost source plugins
for a specific time range.

**Use Cases**:

- Monthly/quarterly cost reports
- Cost trend analysis
- Budget variance tracking

**Input Parameters**:

```json
{
  "stack_name": "string (required) - Pulumi stack name",
  "time_range": {
    "start": "string (required) - ISO 8601 datetime",
    "end": "string (required) - ISO 8601 datetime"
  },
  "granularity": "string (optional) - DAILY, WEEKLY, MONTHLY",
  "filters": {
    "provider": "string (optional)",
    "region": "string (optional)",
    "tags": {
      "key": "string",
      "values": ["string"]
    }
  }
}
```

**Output**:

```json
{
  "total_monthly": 1456.78,
  "currency": "USD",
  "time_series": [
    {
      "timestamp": "2024-01-01T00:00:00Z",
      "cost": 47.23,
      "breakdown": {
        "compute": 25.00,
        "storage": 12.00,
        "network": 10.23
      }
    }
  ],
  "resources": [
    {
      "urn": "urn:pulumi:prod::myapp::aws:ec2/instance:Instance::web-1",
      "total_cost": 234.50,
      "daily_average": 7.56
    }
  ]
}
```

**Example Usage**:

```
User: Show me actual costs for the last 30 days

Claude: [Calls get_actual_cost with 30-day range]

Here's your actual spending for the last 30 days:
- Total: $1,456.78 USD
- Daily Average: $48.56
- Highest Day: $67.89 (Jan 15)
- Lowest Day: $32.45 (Jan 7)

Top cost drivers:
1. EC2 instances: $456.78 (31%)
2. RDS databases: $389.12 (27%)
3. Data transfer: $289.45 (20%)
```

---

### compare_costs

Compare costs between configurations, time periods, or environments.

**Description**: Detailed cost comparison with difference calculations and
percentage changes.

**Use Cases**:

- Environment comparison (dev vs prod)
- Before/after configuration changes
- Historical period comparison

**Input Parameters**:

```json
{
  "baseline": {
    "stack_name": "string (optional)",
    "pulumi_json": "string (optional)",
    "filters": {}
  },
  "target": {
    "stack_name": "string (optional)",
    "pulumi_json": "string (optional)",
    "filters": {}
  },
  "comparison_type": "string (optional) - projected, actual, or both"
}
```

**Output**:

```json
{
  "baseline_cost": 1234.56,
  "target_cost": 1567.89,
  "difference": 333.33,
  "difference_percent": 27.0,
  "currency": "USD",
  "resource_comparison": [
    {
      "resource_type": "aws:ec2/instance:Instance",
      "baseline_cost": 234.50,
      "target_cost": 345.67,
      "difference": 111.17,
      "difference_percent": 47.4,
      "change_reason": "Instance type changed from t3.medium to t3.large"
    }
  ],
  "insights": [
    "Target configuration costs 27% more than baseline",
    "Largest increase: EC2 instances (+$111.17, +47.4%)"
  ]
}
```

**Example Usage**:

```
User: Compare costs between staging and production

Claude: [Calls compare_costs with both stack names]

Cost Comparison: Staging vs Production

Baseline (Staging): $1,234.56/month
Target (Production): $1,567.89/month
Difference: +$333.33 (+27.0%)

Key Differences:
• EC2 instances: +47.4% ($111.17 increase)
  - Prod uses larger instance types
• RDS: +15.2% ($76.00 increase)
  - Prod has Multi-AZ enabled
• S3: -8.5% ($4.25 decrease)
  - Similar storage usage
```

---

### analyze_resource_cost

Deep-dive cost analysis for specific resources.

**Description**: Detailed cost breakdown for a single resource including
dependencies and historical trends.

**Use Cases**:

- Resource-level cost investigation
- Dependency cost attribution
- Optimization opportunity identification

**Input Parameters**:

```json
{
  "resource_urn": "string (required) - Resource URN",
  "time_range": {
    "start": "string (optional)",
    "end": "string (optional)"
  },
  "include_dependencies": "boolean (optional) - Include dependent resources"
}
```

**Output**:

```json
{
  "resource": {
    "urn": "urn:pulumi:prod::myapp::aws:ec2/instance:Instance::web-1",
    "name": "web-1",
    "type": "aws:ec2/instance:Instance",
    "monthly_cost": 234.50,
    "cost_trend": "increasing",
    "cost_components": [
      {
        "name": "Instance hours",
        "monthly_cost": 182.50,
        "percentage": 77.8
      },
      {
        "name": "EBS volumes",
        "monthly_cost": 52.00,
        "percentage": 22.2
      }
    ]
  },
  "dependencies": [
    {
      "urn": "urn:pulumi:prod::myapp::aws:ebs/volume:Volume::web-1-data",
      "monthly_cost": 52.00,
      "relationship": "attached"
    }
  ],
  "optimization_hints": [
    "Consider Reserved Instance for 40% savings",
    "EBS volume is over-provisioned (65% unused capacity)"
  ]
}
```

**Example Usage**:

```
User: Why is web-server-1 so expensive?

Claude: [Calls analyze_resource_cost for web-server-1]

Analysis for web-server-1:

Monthly Cost: $234.50
Trend: Increasing (up 15% from last month)

Cost Breakdown:
• Instance hours: $182.50 (77.8%)
  - t3.large, 730 hours/month
• EBS volumes: $52.00 (22.2%)
  - 500GB gp3 volume, 65% unused

Optimization Opportunities:
1. Switch to Reserved Instance: Save $73/month (40%)
2. Right-size EBS volume: Save $18/month (reduce to 200GB)
3. Consider spot instances for non-critical workloads

Total Potential Savings: $91/month (39%)
```

---

### query_cost_by_tags

Group and analyze costs by resource tags.

**Description**: Tag-based cost attribution for multi-tenant, multi-team, or
multi-project environments.

**Use Cases**:

- Team/department cost allocation
- Project cost tracking
- Environment-based reporting

**Input Parameters**:

```json
{
  "stack_name": "string (required)",
  "tag_keys": ["string (required) - Tag keys to group by"],
  "filters": {
    "key": "string (optional) - Additional tag filter",
    "values": ["string"]
  },
  "time_range": {
    "start": "string (optional)",
    "end": "string (optional)"
  }
}
```

**Output**:

```json
{
  "total_cost": 5678.90,
  "currency": "USD",
  "by_tag": {
    "environment": {
      "production": {
        "cost": 3456.78,
        "percentage": 60.9,
        "resource_count": 45
      },
      "staging": {
        "cost": 1234.56,
        "percentage": 21.7,
        "resource_count": 23
      },
      "development": {
        "cost": 987.56,
        "percentage": 17.4,
        "resource_count": 12
      }
    },
    "team": {
      "platform": {
        "cost": 2345.67,
        "percentage": 41.3
      },
      "backend": {
        "cost": 1890.45,
        "percentage": 33.3
      }
    }
  },
  "untagged_cost": 123.45,
  "untagged_resources": 3
}
```

**Example Usage**:

```
User: Break down costs by team and environment

Claude: [Calls query_cost_by_tags with tag_keys: ["team", "environment"]]

Cost Allocation Report:

Total: $5,678.90

By Environment:
• Production: $3,456.78 (60.9%) - 45 resources
• Staging: $1,234.56 (21.7%) - 23 resources
• Development: $987.56 (17.4%) - 12 resources

By Team:
• Platform Team: $2,345.67 (41.3%)
• Backend Team: $1,890.45 (33.3%)
• Frontend Team: $1,319.33 (23.2%)

⚠️ Untagged Resources: $123.45 (2.2%) - 3 resources
Recommendation: Tag all resources for accurate cost attribution
```

---

### analyze_stack

Comprehensive stack analysis with real-time progress updates.

**Description**: Full stack cost analysis with streaming progress updates,
ideal for large infrastructure deployments.

**Use Cases**:

- Complete infrastructure cost audit
- Periodic cost reviews
- Pre-deployment validation

**Input Parameters**:

```json
{
  "stack_name": "string (required)",
  "include_recommendations": "boolean (optional) - Include optimization suggestions"
}
```

**Output** (streaming):

```json
{
  "progress": 25.0,
  "status": "Analyzing compute resources...",
  "partial_results": {
    "resources_analyzed": 12,
    "total_resources": 48
  }
}
```

**Final Output**:

```json
{
  "stack_name": "myapp-prod",
  "total_monthly_cost": 5678.90,
  "currency": "USD",
  "resource_count": 48,
  "provider_breakdown": {
    "aws": 5234.56,
    "datadog": 444.34
  },
  "category_breakdown": {
    "compute": 2345.67,
    "database": 1890.45,
    "storage": 567.89,
    "network": 874.89
  },
  "recommendations": [
    {
      "type": "RIGHTSIZING",
      "potential_savings": 456.78,
      "resources_affected": 5
    }
  ],
  "alerts": [
    {
      "severity": "WARNING",
      "message": "3 EC2 instances running 24/7 in dev environment"
    }
  ]
}
```

**Example Usage**:

```
User: Analyze all costs for the production stack

Claude: [Calls analyze_stack with streaming enabled]

Analyzing production stack...
⏳ Progress: 25% - Analyzing compute resources...
⏳ Progress: 50% - Analyzing databases...
⏳ Progress: 75% - Analyzing storage...
✅ Complete!

Production Stack Analysis:

Total Monthly Cost: $5,678.90

By Category:
• Compute: $2,345.67 (41.3%)
• Database: $1,890.45 (33.3%)
• Storage: $567.89 (10.0%)
• Network: $874.89 (15.4%)

Cost Optimization Opportunities:
1. Right-size 5 over-provisioned instances: Save $456.78/month
2. Enable S3 lifecycle policies: Save $67.89/month
3. Use Reserved Instances: Save $890.12/month

Total Potential Savings: $1,414.79/month (24.9%)
```

---

## Plugin Management Tools

### list_plugins

Discover and list all available cost source plugins.

**Description**: Returns all installed cost source plugins with their
capabilities and optional health status.

**Use Cases**:

- Plugin discovery
- Health monitoring
- Capability checking

**Input Parameters**:

```json
{
  "include_health": "boolean (optional) - Include health check data"
}
```

**Output**:

```json
{
  "plugins": [
    {
      "name": "aws-cost-source",
      "version": "v1.0.0",
      "description": "AWS Cost and Usage Report data source",
      "capabilities": {
        "supports_projected": false,
        "supports_actual": true,
        "supports_providers": ["aws"]
      },
      "health_status": {
        "status": "healthy",
        "last_check": "2024-01-08T10:30:00Z",
        "latency_ms": 12
      }
    },
    {
      "name": "infracost-plugin",
      "version": "v0.10.0",
      "description": "Infracost-based projected cost estimation",
      "capabilities": {
        "supports_projected": true,
        "supports_actual": false,
        "supports_providers": ["aws", "azure", "gcp"]
      },
      "health_status": {
        "status": "healthy",
        "last_check": "2024-01-08T10:30:00Z",
        "latency_ms": 25
      }
    }
  ]
}
```

**Example Usage**:

```
User: What cost plugins are available?

Claude: [Calls list_plugins with include_health: true]

Available Cost Plugins:

1. aws-cost-source (v1.0.0)
   • Type: Actual costs only
   • Provider: AWS
   • Status: ✅ Healthy (12ms latency)
   • Description: AWS Cost and Usage Report data source

2. azure-cost-source (v1.0.0)
   • Type: Actual costs only
   • Provider: Azure
   • Status: ✅ Healthy (18ms latency)
   • Description: Azure Cost Management data source

3. infracost-plugin (v0.10.0)
   • Type: Projected costs only
   • Providers: AWS, Azure, GCP
   • Status: ✅ Healthy (25ms latency)
   • Description: Infracost-based projected cost estimation

All plugins are operational.
```

---

### get_plugin_info

Get detailed information about a specific plugin.

**Description**: Returns complete plugin metadata, configuration, capabilities,
and health status.

**Use Cases**:

- Plugin configuration review
- Capability verification
- Troubleshooting

**Input Parameters**:

```json
{
  "plugin_name": "string (required)"
}
```

**Output**:

```json
{
  "name": "aws-cost-source",
  "version": "v1.0.0",
  "description": "AWS Cost and Usage Report data source",
  "capabilities": {
    "supports_projected": false,
    "supports_actual": true,
    "supports_providers": ["aws"]
  },
  "health_status": {
    "status": "healthy",
    "last_check": "2024-01-08T10:30:00Z",
    "latency_ms": 12
  },
  "grpc_address": "localhost:50051",
  "configuration": {
    "region": "us-east-1",
    "bucket": "my-cur-bucket",
    "report_name": "cost-usage-report",
    "poll_interval": "1h"
  }
}
```

**Example Usage**:

```
User: Tell me about the aws-cost-source plugin

Claude: [Calls get_plugin_info for "aws-cost-source"]

AWS Cost Source Plugin (v1.0.0)

Description: AWS Cost and Usage Report data source

Capabilities:
• Supports actual costs: Yes
• Supports projected costs: No
• Supported providers: AWS

Current Status: ✅ Healthy
• Last health check: 2 minutes ago
• Response latency: 12ms
• gRPC endpoint: localhost:50051

Configuration:
• Region: us-east-1
• S3 bucket: my-cur-bucket
• Report name: cost-usage-report
• Poll interval: 1 hour

The plugin is operational and ready to query actual AWS costs.
```

---

### validate_plugin

Validate plugin conformance to pulumicost-spec.

**Description**: Runs conformance tests against a plugin to verify it meets
the pulumicost-spec requirements.

**Use Cases**:

- Plugin development testing
- Quality assurance
- Conformance certification

**Input Parameters**:

```json
{
  "plugin_path": "string (required) - Path to plugin binary",
  "conformance_level": "string (required) - BASIC, STANDARD, or FULL"
}
```

**Output**:

```json
{
  "plugin_name": "/path/to/my-plugin",
  "conformance_level": "STANDARD",
  "passed": true,
  "test_results": [
    {
      "name": "gRPC Interface",
      "passed": true,
      "duration_ms": 15
    },
    {
      "name": "Health Check Endpoint",
      "passed": true,
      "duration_ms": 10
    },
    {
      "name": "Cost Query Basic",
      "passed": true,
      "duration_ms": 25
    },
    {
      "name": "Resource Filtering",
      "passed": true,
      "duration_ms": 20
    }
  ],
  "timestamp": "2024-01-08T10:30:00Z"
}
```

**Example Usage**:

```
User: Validate my custom cost plugin

Claude: [Calls validate_plugin with plugin path and STANDARD level]

Plugin Validation Report

Plugin: /usr/local/bin/my-cost-plugin
Conformance Level: STANDARD
Overall Result: ✅ PASSED

Test Results (4/4 passed):

✅ gRPC Interface (15ms)
   Verified plugin implements required gRPC services

✅ Health Check Endpoint (10ms)
   Plugin responds to health checks correctly

✅ Cost Query Basic (25ms)
   Plugin can execute basic cost queries

✅ Resource Filtering (20ms)
   Plugin correctly filters resources by criteria

Your plugin meets STANDARD conformance requirements and can
be used with the PulumiCost ecosystem.
```

---

### health_check

Check health and connectivity of a plugin.

**Description**: Performs a health check on a specific plugin and returns
current status.

**Use Cases**:

- Plugin monitoring
- Troubleshooting connectivity
- Pre-query validation

**Input Parameters**:

```json
{
  "plugin_name": "string (required)"
}
```

**Output**:

```json
{
  "status": "healthy",
  "last_check": "2024-01-08T10:30:00Z",
  "latency_ms": 12,
  "error_message": null
}
```

**Possible Statuses**:

- `healthy` - Plugin is operational
- `degraded` - Plugin is slow but functional
- `unhealthy` - Plugin is not responding

**Example Usage**:

```
User: Is the AWS cost plugin working?

Claude: [Calls health_check for "aws-cost-source"]

AWS Cost Source Plugin Health: ✅ Healthy

• Status: Operational
• Last checked: Just now (2024-01-08 10:30:00 UTC)
• Response time: 12ms
• No errors detected

The plugin is ready to serve cost queries.
```

**Error Example**:

```
User: Check the Azure plugin

Claude: [Calls health_check for "azure-cost-source"]

Azure Cost Source Plugin Health: ⚠️ Degraded

• Status: Degraded
• Last checked: Just now
• Response time: 500ms (high latency)
• Error: High latency detected

The plugin is functional but experiencing performance issues.
Consider checking network connectivity or plugin configuration.
```

---

## Analysis and Optimization Tools

### get_recommendations

AI-powered cost optimization recommendations.

**Description**: Analyzes infrastructure and generates actionable cost
optimization recommendations.

**Use Cases**:

- Cost optimization planning
- Regular cost reviews
- Budget reduction initiatives

**Input Parameters**:

```json
{
  "stack_name": "string (required)",
  "recommendation_types": ["string (optional) - Filter by type"],
  "minimum_savings": "number (optional) - Minimum monthly savings threshold"
}
```

**Recommendation Types**:

- `RIGHTSIZING` - Over/under-provisioned resources
- `RESERVED_INSTANCES` - RI opportunities
- `SPOT_INSTANCES` - Spot instance candidates
- `STORAGE_OPTIMIZATION` - Storage efficiency
- `IDLE_RESOURCES` - Unused resources

**Output**:

```json
{
  "recommendations": [
    {
      "id": "rec-001",
      "type": "RIGHTSIZING",
      "resource_urn": "urn:pulumi:prod::myapp::aws:ec2/instance:Instance::web-1",
      "current_cost": 120.00,
      "projected_savings": 245.50,
      "confidence": "HIGH",
      "description": "EC2 instance is over-provisioned",
      "action_steps": [
        "Review instance metrics",
        "Update instance type from t3.large to t3.medium",
        "Monitor performance after change"
      ]
    }
  ]
}
```

**Example Usage**:

```
User: What cost optimizations can I make?

Claude: [Calls get_recommendations]

Cost Optimization Recommendations:

Found 3 opportunities to save $481.50/month (18.7% reduction)

1. 🎯 Right-size EC2 Instance (HIGH confidence)
   Resource: web-1
   Current: $120.00/month
   Savings: $245.50/month

   Action: Downsize from t3.large to t3.medium
   • CPU utilization averages 35%
   • Memory usage averages 42%

   Steps:
   1. Review instance metrics over last 30 days
   2. Update instance type in Pulumi code
   3. Deploy and monitor performance

2. 💰 Purchase Reserved Instances (MEDIUM confidence)
   Resource: db-prod
   Current: $250.00/month
   Savings: $180.00/month

   Action: Buy 1-year Reserved Instance
   • Database runs 24/7
   • Consistent usage pattern

3. ☁️ Use Spot Instances (MEDIUM confidence)
   Resource: dev-worker-1
   Current: $80.00/month
   Savings: $56.00/month

   Action: Switch to Spot Instances
   • Development workload tolerates interruptions
   • 70% cost savings

Implement all recommendations to save $481.50/month
```

---

### detect_anomalies

Detect unusual cost patterns and spending anomalies.

**Description**: Uses statistical analysis to identify abnormal cost patterns
and spending spikes.

**Use Cases**:

- Cost spike investigation
- Budget overrun prevention
- Proactive cost monitoring

**Input Parameters**:

```json
{
  "stack_name": "string (required)",
  "time_range": {
    "start": "string (required)",
    "end": "string (required)"
  },
  "sensitivity": "string (optional) - LOW, MEDIUM, or HIGH"
}
```

**Output**:

```json
{
  "anomalies": [
    {
      "id": "anom-001",
      "timestamp": "2024-01-15T14:30:00Z",
      "resource_urns": [
        "urn:pulumi:prod::myapp::aws:ec2/instance:Instance::web-3"
      ],
      "severity": "HIGH",
      "current_cost": 425.00,
      "baseline_cost": 150.00,
      "deviation_percent": 183.33,
      "potential_causes": [
        "Unexpected spike in EC2 instance hours",
        "Auto-scaling triggered",
        "Instance not properly terminated"
      ]
    }
  ]
}
```

**Example Usage**:

```
User: Have there been any unusual cost spikes this month?

Claude: [Calls detect_anomalies with MEDIUM sensitivity]

Cost Anomaly Detection Results:

Found 2 anomalies in the last 30 days:

🚨 HIGH SEVERITY - Jan 15, 2:30 PM
   Resource: web-3 EC2 instance

   Anomaly: Cost spike of 183.3%
   • Baseline: $150.00/day
   • Actual: $425.00/day
   • Deviation: +$275.00 (+183%)

   Potential Causes:
   • Unexpected spike in instance hours
   • Auto-scaling may have triggered
   • Instance possibly not terminated properly

   Recommendation: Review instance lifecycle and scaling policies

⚠️ MEDIUM SEVERITY - Jan 20, 9:15 AM
   Resource: data S3 bucket

   Anomaly: Storage cost increase of 60%
   • Baseline: $75.00/day
   • Actual: $120.00/day
   • Deviation: +$45.00 (+60%)

   Potential Causes:
   • Large file uploads detected
   • Lifecycle policies may not be applied

   Recommendation: Review S3 lifecycle configuration

Action Required: Investigate high severity anomaly immediately
```

---

### forecast_costs

Forecast future costs based on historical trends.

**Description**: Generates cost forecasts using historical data and trend
analysis with confidence intervals.

**Use Cases**:

- Budget planning
- Quarterly projections
- Capacity planning

**Input Parameters**:

```json
{
  "stack_name": "string (required)",
  "forecast_period": {
    "start": "string (required)",
    "end": "string (required)"
  },
  "confidence_level": "number (optional) - 0.0 to 1.0, default 0.95"
}
```

**Output**:

```json
{
  "stack_name": "my-stack",
  "forecast_period": {
    "start": "2024-02-01T00:00:00Z",
    "end": "2024-02-29T23:59:59Z"
  },
  "data_points": [
    {
      "timestamp": "2024-02-01T00:00:00Z",
      "predicted_cost": 850.00,
      "lower_bound": 800.00,
      "upper_bound": 900.00
    },
    {
      "timestamp": "2024-02-15T00:00:00Z",
      "predicted_cost": 875.00,
      "lower_bound": 825.00,
      "upper_bound": 925.00
    }
  ],
  "confidence_level": 0.95,
  "methodology": "Linear regression with seasonal adjustment"
}
```

**Example Usage**:

```
User: What will my infrastructure cost next quarter?

Claude: [Calls forecast_costs for next quarter with 0.95 confidence]

Cost Forecast: Q2 2024 (Apr 1 - Jun 30)

Predicted Monthly Costs:
• April: $850.00 (range: $800-$900)
• May: $875.00 (range: $825-$925)
• June: $900.00 (range: $850-$950)

Total Forecast: $2,625.00
Confidence Level: 95%

Trend Analysis:
• Average monthly increase: 2.9%
• Growth driver: Increased EC2 utilization
• Seasonal pattern: Slight uptick in spring months

Methodology: Linear regression with seasonal adjustment
based on 90 days of historical data

💡 Budget Recommendation: Plan for $3,000 to account for
   upper confidence bound and unexpected growth
```

---

### track_budget

Track spending against defined budgets with alerts.

**Description**: Monitor current spending against budget targets with burn
rate analysis and alert thresholds.

**Use Cases**:

- Budget monitoring
- Cost governance
- Spending alerts

**Input Parameters**:

```json
{
  "stack_name": "string (required)",
  "budget_amount": "number (required) - Budget in currency units",
  "period": "string (required) - MONTHLY, QUARTERLY, ANNUALLY",
  "alert_thresholds": ["number (optional) - Alert at percentage thresholds"]
}
```

**Output**:

```json
{
  "budget_amount": 1000.00,
  "current_spending": 750.00,
  "remaining": 250.00,
  "burn_rate": 25.00,
  "projected_end_date": "2024-01-18T00:00:00Z",
  "status": "WARNING",
  "alerts": [
    {
      "threshold": 50.0,
      "current_spend": 750.00,
      "severity": "MEDIUM",
      "message": "Budget utilization at 75.0% (threshold: 50.0%)",
      "timestamp": "2024-01-08T10:30:00Z"
    }
  ]
}
```

**Status Values**:

- `OK` - Under 80% of budget
- `WARNING` - 80-90% of budget used
- `CRITICAL` - 90-100% of budget used
- `EXCEEDED` - Over budget

**Example Usage**:

```
User: Am I on track with my $1000 monthly budget?

Claude: [Calls track_budget with $1000 budget, MONTHLY period]

Budget Status: ⚠️ WARNING

Monthly Budget: $1,000.00
Current Spending: $750.00 (75.0%)
Remaining: $250.00

Burn Rate: $25.00/day
Projected End Date: Jan 18 (10 days early)

⚠️ Active Alerts:
• WARNING: Budget utilization at 75% (threshold: 50%)
  Triggered 2 hours ago

Forecast: At current burn rate, you'll exceed budget by
$250 before month end (Jan 31)

Recommendations:
1. Review recent cost increases
2. Consider implementing cost controls
3. Optimize high-cost resources (see get_recommendations)

Take action now to avoid budget overrun!
```

---

## Error Handling

All tools use consistent error types:

### ValidationError

Invalid input parameters.

```json
{
  "error": "validation_error",
  "message": "stack_name cannot be empty",
  "field": "stack_name"
}
```

### NotFoundError

Resource not found.

```json
{
  "error": "not_found",
  "message": "plugin 'aws-cost-source' not found",
  "resource": "aws-cost-source"
}
```

### BadRequestError

Invalid request format.

```json
{
  "error": "bad_request",
  "message": "invalid time range: start must be before end"
}
```

### InternalError

Server-side error.

```json
{
  "error": "internal_error",
  "message": "failed to connect to cost source plugin"
}
```

---

## Best Practices

### 1. Cost Query Efficiency

- Use filters to narrow results
- Request appropriate time granularity
- Cache frequently accessed data

### 2. Budget Management

- Set multiple alert thresholds (50%, 80%, 100%)
- Review budgets monthly
- Adjust forecasts based on business changes

### 3. Optimization Workflow

```
1. Run analyze_stack for overview
2. Call get_recommendations for opportunities
3. Use analyze_resource_cost for deep dives
4. Implement changes incrementally
5. Use track_budget to monitor savings
```

### 4. Plugin Health

- Enable health checks in list_plugins
- Monitor latency trends
- Validate plugins after updates

### 5. Anomaly Detection

- Use MEDIUM sensitivity for balanced detection
- Investigate HIGH severity anomalies immediately
- Set up regular anomaly scans (daily/weekly)

---

## Integration Examples

### Claude Desktop Workflow

```
# Morning cost review
User: Check my budget status
Claude: [track_budget] ✅ 45% of monthly budget used, on track

# Investigation
User: Any unusual costs yesterday?
Claude: [detect_anomalies] 🚨 Found spike in EC2 costs

# Deep dive
User: Analyze that EC2 instance
Claude: [analyze_resource_cost] Instance running 24/7 in dev

# Optimization
User: How can I reduce costs?
Claude: [get_recommendations] 3 opportunities: save $450/month

# Planning
User: Forecast next month
Claude: [forecast_costs] Projected: $1,250 (±$100)
```

### API Integration

```bash
# JSON-RPC call
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "cost.analyze_projected",
    "params": {
      "pulumi_json": "...",
      "filters": {
        "provider": "aws"
      }
    },
    "id": 1
  }'
```

---

## See Also

- [User Guide](user-guide.md) - Getting started
- [Developer Guide](developer-guide.md) - Development setup
- [Plugin Development](plugin-development.md) - Building plugins
- [API Reference](../api/) - Complete API documentation
