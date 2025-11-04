# Cost Analyst - PulumiCost MCP Server

## Role Context

You are a Cloud Cost Analyst (FinOps Practitioner) using the PulumiCost MCP Server to analyze, optimize, and manage cloud infrastructure costs. Your focus is on cost visibility, attribution, optimization, and helping teams make cost-effective decisions.

## Key Responsibilities

- **Cost Analysis**: Analyze current and projected infrastructure costs
- **Cost Attribution**: Allocate costs to teams, projects, and environments
- **Optimization**: Identify cost-saving opportunities
- **Reporting**: Generate cost reports for stakeholders
- **Budgeting**: Track spending against budgets
- **Forecasting**: Project future costs based on trends
- **Education**: Help teams understand and reduce costs

## Using PulumiCost MCP via AI Assistant

### Setting Up

Once your AI assistant (e.g., Claude Desktop) is configured with the PulumiCost MCP server, you can interact with it using natural language.

### Example Conversations

#### Basic Cost Queries

**Projected Costs:**
```
You: What will my production infrastructure cost next month?

AI: [Uses analyze_projected_costs tool]
Based on your current Pulumi stack, projected monthly costs:

Total: $4,234.50/month

Breakdown:
- Compute (EC2): $2,456.00 (58%)
- Database (RDS): $1,234.00 (29%)
- Storage (S3): $345.50 (8%)
- Network: $199.00 (5%)

Top 5 Expensive Resources:
1. prod-db-primary (RDS db.r5.xlarge): $734/mo
2. prod-app-cluster (EKS): $584/mo
3. prod-web-lb (ALB): $456/mo
4. prod-cache (ElastiCache): $389/mo
5. prod-db-read-replica (RDS db.r5.large): $367/mo
```

**Actual Costs:**
```
You: Show me actual costs for the last 30 days

AI: [Uses get_actual_costs tool]
Actual costs (Last 30 days): $4,567.89

Daily average: $152.26
Trend: ↑ 8% vs previous period

Breakdown by service:
- EC2: $2,678.45 (59%) ↑ 12%
- RDS: $1,245.67 (27%) ↔ 0%
- S3: $389.23 (9%) ↓ 5%
- Other: $254.54 (5%) ↑ 3%

Notable changes:
- EC2 costs increased 12% due to new instances launched on 1/15
- S3 costs decreased 5% after cleanup on 1/10
```

**Cost Comparison:**
```
You: Compare projected vs actual costs

AI: [Uses compare_costs tool]
Projected vs Actual Comparison:

Total:
- Projected: $4,234.50
- Actual: $4,567.89
- Variance: +$333.39 (+7.9%) ⚠️

Over-budget items:
1. EC2: +$222.45 (+9.1%)
   Reason: Additional instances for load handling
2. S3: +$43.73 (+12.7%)
   Reason: Higher than expected data transfer
3. Network: +$67.21 (+33.8%) ⚠️
   Reason: Unexpected cross-region traffic

On-budget items:
- RDS: +$11.67 (+0.9%) ✓

Recommendation: Investigate network costs and EC2 usage patterns.
```

#### Cost Attribution

**By Team:**
```
You: Break down costs by team

AI: [Uses analyze_costs with grouping]
Cost attribution by team (January):

1. Engineering Team: $2,345.67 (51%)
   - prod environments: $1,678.90
   - staging: $456.77
   - dev: $210.00

2. Data Team: $1,234.00 (27%)
   - data warehouse: $890.00
   - analytics: $344.00

3. Platform Team: $789.22 (17%)
   - shared services: $567.00
   - monitoring: $222.22

4. Unallocated: $199.00 (4%)
   - Resources without team tags

Action items:
- Tag unallocated resources
- Review Engineering Team's prod spending
```

**By Environment:**
```
You: Compare costs across environments

AI: Cost by environment:

Production: $3,456.78 (76%)
- Expected for production workload ✓

Staging: $789.12 (17%)
- Higher than typical 10-15% of prod ⚠️
- Recommendation: Consider downsizing or shutting down when not in use

Development: $321.99 (7%)
- Within acceptable range ✓

Recommendation: Staging environment is oversized. Potential monthly savings: ~$250
```

#### Optimization Analysis

**Finding Savings:**
```
You: Find cost optimization opportunities

AI: [Uses get_optimization_recommendations tool]
Cost Optimization Opportunities:

High Impact (Est. savings: $567/month):

1. Right-size RDS instances ($234/month)
   Current: db.r5.xlarge (4 vCPU, 32 GB)
   Recommendation: db.r5.large (2 vCPU, 16 GB)
   Rationale: CPU utilization < 20%, memory < 40%
   Risk: Low - gradual scaling recommended

2. Convert to Reserved Instances ($198/month)
   Resources: 3x t3.large EC2 instances
   Current: On-Demand
   Recommendation: 1-year Reserved Instance
   Rationale: Running 24/7 for >6 months
   Risk: Very Low - commitment aligns with usage

3. Implement S3 Lifecycle Policies ($135/month)
   Bucket: prod-logs-bucket (1.2 TB)
   Current: S3 Standard
   Recommendation: Transition to S3-IA after 30 days, Glacier after 90 days
   Rationale: 80% of objects not accessed after 30 days
   Risk: Low - logs rarely accessed

Medium Impact (Est. savings: $289/month):

4. Schedule dev/staging shutdown ($156/month)
   Resources: Dev and staging EC2/RDS instances
   Recommendation: Stop outside business hours (6PM-8AM, weekends)
   Potential savings: 70% off-hours reduction
   Risk: Low - implement with monitoring

5. Clean up unused resources ($133/month)
   Found: 5 unused EBS volumes, 3 unattached EIPs, 2 idle load balancers
   Action: Audit and delete if confirmed unused
   Risk: Medium - verify with teams first

Total potential savings: $856/month ($10,272/year)
```

**Resource Analysis:**
```
You: Analyze costs for the web-server EC2 instance

AI: [Uses analyze_resource tool]
Resource Cost Analysis: prod-web-server-01

Current Month:
- Cost: $234.56
- Usage: 730 hours (100% uptime)
- Instance Type: t3.large
- Region: us-east-1

Cost Breakdown:
- Compute: $189.80 (81%)
- EBS Volume (100GB gp3): $32.00 (14%)
- Data Transfer: $12.76 (5%)

Historical Trend (90 days):
- Average: $228.45/month
- Min: $219.33 (December)
- Max: $245.78 (November)
- Trend: Stable ✓

Utilization:
- CPU: 45% average (healthy)
- Memory: 62% average (healthy)
- Network: 2.3 GB/day out

Recommendations:
- ✓ Size is appropriate for workload
- Consider Reserved Instance for 12% savings ($28/month)
- Monitor data transfer costs (increasing trend)
```

#### Budgeting and Alerts

**Budget Tracking:**
```
You: Are we on track with our monthly budget of $5,000?

AI: Budget Tracking (January 1-15):

Current spend: $2,456.78
Budget: $5,000.00
Used: 49% of budget, 50% of month elapsed

Projection for month: $4,913.56
Status: On track ✓

Daily average: $163.79
Remaining budget: $2,543.22
Days remaining: 15
Allowed daily spend: $169.55

Forecast: You'll end the month at 98.3% of budget with current usage patterns.
```

**Anomaly Detection:**
```
You: Are there any unusual cost spikes?

AI: Cost Anomaly Detection:

⚠️ ALERT: Network costs anomaly detected

Date: January 14, 2025
Normal: $6-8/day
Detected: $67.23/day (+840%)

Root cause analysis:
- Cross-region data transfer from us-east-1 to eu-west-1
- Volume: 245 GB (vs typical 3-5 GB)
- Likely cause: Unintentional cross-region replication or misconfigured backup

Cost impact: ~$60 extra on that day

Recommendation:
1. Review network configuration for cross-region transfers
2. Investigate why eu-west-1 traffic spiked
3. Consider using VPC endpoints to reduce transfer costs
4. Set up alerts for cross-region traffic > 10 GB/day

Additional findings:
- EC2 costs were 15% higher on Jan 15 (new instances)
- S3 costs show gradual increase (normal growth pattern)
```

#### Forecasting

**Trend Analysis:**
```
You: Project costs for next quarter based on current trends

AI: Cost Forecast (Q2 2025):

Current monthly average: $4,456.78
Growth rate (last 3 months): +5% MoM

Projected costs:
- April: $4,679.62
- May: $4,913.60
- June: $5,159.28
- Q2 Total: $14,752.50

Confidence: Medium (based on 3-month trend)

Assumptions:
- Current growth rate continues
- No major infrastructure changes
- Seasonal patterns remain consistent

Factors that could increase costs:
- New product features (+10-20%)
- Traffic growth (+5-15%)
- Additional environments (+5-10%)

Factors that could decrease costs:
- Optimization initiatives (-10-15%)
- Reserved Instance conversions (-10-12%)
- Resource cleanup (-3-5%)

Recommended actions:
1. Plan optimization efforts to offset growth
2. Budget $15,000 for Q2 with 15% buffer
3. Implement cost controls before growth accelerates
```

## Cost Analysis Workflows

### Monthly Cost Review

```
1. You: Show me the cost summary for last month

2. You: Compare it to the previous month

3. You: What drove the largest cost changes?

4. You: Break down by team and environment

5. You: Generate a monthly cost report
```

### Budget Planning

```
1. You: What were our costs for the last 6 months?

2. You: Show me the growth trend

3. You: Project costs for the next 6 months

4. You: What's a realistic budget for Q2?

5. You: What optimization opportunities exist?
```

### Optimization Sprint

```
1. You: Find all cost optimization opportunities

2. You: Prioritize by potential savings and risk

3. You: For the top 3, give me detailed analysis

4. You: What's the implementation plan for each?

5. You: Track savings after implementation
```

### New Project Cost Estimation

```
1. You: I'm planning to deploy a new microservice. Here are the resources: [paste Pulumi config]

2. You: What will the monthly cost be?

3. You: Compare costs across different instance sizes

4. You: What if we use spot instances?

5. You: What are the cost risks?
```

## FinOps Best Practices

### 1. Tagging Strategy

Ensure resources are tagged for proper attribution:

```
Required Tags:
- team: engineering|data|platform
- environment: prod|staging|dev
- project: project-name
- owner: email@company.com
- cost-center: CC-1234
```

```
You: Show me resources missing required tags

AI: Untagged resources found: 23 items

Missing 'team' tag:
- prod-web-lb ($456/month)
- staging-db-01 ($234/month)
- [21 more...]

Impact: $1,234/month cannot be attributed to teams

Recommendation: Implement automated tagging policy
```

### 2. Regular Cost Reviews

Schedule regular reviews:

**Weekly:**
- Check for anomalies
- Review daily spending trends
- Address urgent issues

**Monthly:**
- Full cost review
- Budget tracking
- Optimization planning
- Stakeholder reporting

**Quarterly:**
- Strategic planning
- Reserved Instance renewals
- Long-term forecasting
- ROI analysis

### 3. Cost Allocation

```
You: Create a chargeback report for January

AI: Cost Chargeback Report - January 2025

Engineering Team - $2,345.67
---
Production: $1,678.90
- web-servers: $789.00
- app-servers: $567.00
- databases: $322.90

Staging: $456.77
- staging-env: $456.77

Development: $210.00
- dev-environments: $210.00

Recommendations:
- Staging is 27% of production (target: 15%)
- Consider rightsizing staging environment

[Similar breakdowns for other teams...]
```

### 4. Cost Optimization Metrics

Track these KPIs:

```
You: Show me our FinOps KPIs

AI: FinOps Performance Metrics:

Cost Efficiency:
- Cost per customer: $2.34 (target: <$2.50) ✓
- Cost per transaction: $0.0045 (target: <$0.005) ✓
- Idle resource cost: 8% (target: <5%) ⚠️

Budget Management:
- Budget adherence: 98% (excellent) ✓
- Forecast accuracy: 94% (good) ✓
- Unexpected costs: $234 (target: <$500) ✓

Optimization:
- Reserved Instance coverage: 67% (target: >75%) ⚠️
- Spot Instance usage: 23% (target: >30%) ⚠️
- Savings realized YTD: $12,345 ✓

Action items:
1. Increase RI coverage to 75%
2. Reduce idle resources
3. Expand spot instance usage
```

## Reporting Templates

### Executive Summary

```
You: Create an executive cost summary for January

AI: Executive Cost Summary - January 2025

TOTAL SPEND: $4,567.89
vs Budget: $5,000 (91% utilized) ✓
vs Last Month: +7.8% ($330 increase)

KEY METRICS:
✓ Under budget by $432
✓ Cost per customer decreased 3%
⚠️ Total spend increased 8% MoM

TOP COST DRIVERS:
1. Compute (59%): $2,678 - New instances for growth
2. Database (27%): $1,246 - Expected steady state
3. Storage (9%): $389 - Decreased 5% (cleanup)

HIGHLIGHTS:
✓ Implemented S3 lifecycle policies (saving $135/mo)
✓ Right-sized 3 RDS instances (saving $234/mo)
⚠️ Network costs spike on 1/14 (investigation underway)

Q2 OUTLOOK:
- Projected: $15,000 (steady growth)
- Optimization target: $10,000 in savings
- At-risk: Budget if growth exceeds projections
```

### Team Cost Report

```
You: Create a detailed cost report for the Engineering team

AI: [Generates comprehensive team report with breakdowns, trends, and recommendations]
```

## Advanced Analysis

### What-If Analysis

```
You: What if we migrate our RDS databases to Aurora?

AI: Migration Cost Analysis: RDS → Aurora

Current RDS Cost: $1,234/month
Projected Aurora Cost: $1,089/month

Monthly savings: $145 (12% reduction)
Annual savings: $1,740

Breakdown:
- Compute: -$98/month (more efficient)
- Storage: -$67/month (pay-per-GB vs provisioned)
- I/O: +$20/month (Aurora charges for I/O)

Additional benefits:
- Better performance (5x faster)
- Automatic scaling
- Better availability
- Easier replication

Migration costs:
- Downtime: ~2 hours
- Testing effort: ~40 hours
- Risk: Medium

ROI: Payback in <1 month

Recommendation: Proceed with migration ✓
```

### Scenario Planning

```
You: Model costs for 3 scenarios: slow growth (5%), expected growth (15%), rapid growth (30%)

AI: [Creates detailed scenario analysis with projections, resource requirements, and budget recommendations for each scenario]
```

## Tools and Queries Reference

### Available MCP Tools

1. **analyze_projected_costs**: Estimate future costs
2. **get_actual_costs**: Historical spending data
3. **compare_costs**: Compare time periods or stacks
4. **analyze_resource**: Detailed resource analysis
5. **get_optimization_recommendations**: Find savings
6. **validate_budget**: Check budget status
7. **detect_anomalies**: Find cost spikes
8. **forecast_costs**: Project future spending
9. **generate_report**: Create formatted reports

### Query Patterns

**Time Ranges:**
```
- "last 7 days"
- "last month"
- "January 2025"
- "Q4 2024"
- "year to date"
- "between Jan 1 and Jan 31"
```

**Grouping:**
```
- "by service"
- "by resource type"
- "by team"
- "by environment"
- "by project"
- "by region"
```

**Filtering:**
```
- "only production"
- "exclude development"
- "team=engineering"
- "cost > $100"
- "AWS EC2 instances"
```

## Resources

- [FinOps Framework](https://www.finops.org/)
- [Cloud Cost Optimization Guide](../docs/guides/cost-optimization.md)
- [Tagging Best Practices](../docs/guides/tagging-strategy.md)
- [Cost Allocation Guide](../docs/guides/cost-allocation.md)

---

**Remember**: Good FinOps is about enabling teams to make cost-effective decisions, not just cutting costs. Focus on optimization, visibility, and education. Use data to drive conversations and decisions.
