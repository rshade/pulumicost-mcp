# Example Queries for Simple AWS Stack

This document contains example queries you can use with Claude Desktop or the PulumiCost MCP Server to analyze this stack's costs.

## Basic Cost Queries

### Total Projected Cost

```
I'm about to deploy this Pulumi stack. What will be the total monthly cost?
```

**Expected Response**: Breakdown showing ~$50-75/month total, with detailed costs for:
- EC2 instance (t3.micro): ~$7.50/month
- RDS PostgreSQL (db.t3.small): ~$25-30/month
- S3 storage: ~$0.50-2/month (depending on usage)
- Data transfer: ~$5-10/month
- Other resources (VPC, security groups): minimal/free

### Resource-by-Resource Breakdown

```
Break down the costs by resource type for this stack.
```

**Expected Response**: Table or list showing:
1. Database (RDS) - largest cost component
2. Compute (EC2) - second largest
3. Storage (S3) - variable based on usage
4. Network (data transfer, NAT) - usage-based
5. Free tier resources (VPC, security groups, etc.)

### Top Cost Drivers

```
What are the top 3 most expensive resources in this stack?
```

**Expected Response**:
1. RDS PostgreSQL instance (~45% of total)
2. EC2 instance (~15% of total)
3. Data transfer costs (~10-15% of total)

## Optimization Queries

### Right-Sizing Recommendations

```
Are there any opportunities to reduce costs through right-sizing?
```

**Expected Response**: Recommendations such as:
- Consider t3.nano for development workloads if lower traffic
- Use Aurora Serverless for variable database workloads
- Enable S3 Intelligent-Tiering for lifecycle management
- Consider Reserved Instances for 1+ year commitments

### Reserved Instance Savings

```
How much could I save with Reserved Instances?
```

**Expected Response**:
- 1-year no-upfront: ~30% savings on EC2 and RDS
- 3-year all-upfront: ~60% savings on EC2 and RDS
- Specific dollar amounts for this stack

### Serverless Alternative

```
What would it cost to run this as a serverless architecture instead?
```

**Expected Response**: Comparison showing:
- Lambda + API Gateway + Aurora Serverless
- Cost savings at low-medium traffic
- Break-even point analysis
- Trade-offs (cold starts, complexity)

## Comparison Queries

### Regional Cost Comparison

```
Compare the cost of running this stack in us-east-1 vs us-west-2
```

**Expected Response**:
- Side-by-side cost comparison
- Price differences by resource type
- Total monthly cost delta
- Recommendation based on requirements

### Instance Type Comparison

```
Compare costs between t3.micro, t3.small, and t3.medium for the web server
```

**Expected Response**:
- t3.micro: $7.50/month
- t3.small: $15/month (2x cost, 2x capacity)
- t3.medium: $30/month (4x cost, 4x capacity)
- Cost per vCPU and GB RAM analysis

### Database Alternatives

```
Compare RDS PostgreSQL vs Aurora PostgreSQL vs DynamoDB for this use case
```

**Expected Response**:
- RDS: Fixed cost, traditional relational
- Aurora: Higher cost, better performance/scalability
- DynamoDB: Pay-per-request, NoSQL trade-offs
- Recommendation based on workload

## Forecasting Queries

### Growth Projection

```
What would costs look like if traffic doubled over the next 6 months?
```

**Expected Response**:
- EC2: May need larger instance (~2x cost)
- RDS: May need larger instance or read replicas (~2-3x cost)
- Data transfer: ~2x cost
- S3: ~2x storage cost
- Total projected: ~$120-150/month

### Budget Tracking

```
Track our spending against a $100/month budget for this stack
```

**Expected Response**:
- Current spend: $65/month
- Remaining budget: $35
- Days until budget exhausted: N/A (under budget)
- Burn rate: ~$2.15/day
- Status: Within budget (65% used)

## Detailed Analysis

### Resource Dependencies

```
Analyze the cost of the web server including all its dependencies
```

**Expected Response**:
- Web server (EC2): $7.50/month
- Security group: Free
- Public subnet: Free
- Internet gateway: Free
- Data transfer: ~$5/month
- Total stack cost: ~$12.50/month

### Tag-Based Analysis

```
Show me costs grouped by the Environment tag
```

**Expected Response**:
- Environment=dev: $65/month (all resources)
- Breakdown by resource within dev environment
- Percentage of total by tag value

### Storage Analysis

```
Analyze S3 costs with current lifecycle policies
```

**Expected Response**:
- Standard storage (0-30 days): $X/month
- Standard-IA (30-90 days): $Y/month
- Glacier (90+ days): $Z/month
- Total estimated: Based on 100GB example
- Savings from lifecycle policies: XX%

## Time-Based Queries

### Historical Trends

```
Show me cost trends for this stack over the past 30 days
```

**Expected Response** (after deployment):
- Daily cost breakdown
- Trend line (increasing/decreasing/stable)
- Anomalies detected
- Cost drivers that changed

### Anomaly Detection

```
Are there any unusual cost spikes in this stack?
```

**Expected Response** (after deployment):
- List of detected anomalies
- Date/time of spike
- Resource causing spike
- Percentage deviation from normal
- Possible causes

## What-If Scenarios

### Scaling Scenario

```
What if I add 2 more EC2 instances behind a load balancer?
```

**Expected Response**:
- Additional EC2 costs: $15/month (2x t3.micro)
- ALB cost: ~$16-20/month
- Data transfer increase: ~$5/month
- Total increase: ~$36-40/month
- New total: ~$100-110/month

### Multi-AZ Scenario

```
What's the cost impact of making the RDS instance Multi-AZ?
```

**Expected Response**:
- Current single-AZ: $25/month
- Multi-AZ: $50/month (2x)
- Additional data transfer: ~$2/month
- Total increase: ~$27/month
- Benefit: High availability, automatic failover

### Development vs Production

```
Compare this dev stack to a production configuration
```

**Expected Response**:
Production changes:
- t3.medium (vs t3.micro): +$22.50/month
- db.m5.large Multi-AZ (vs db.t3.small): +$175/month
- Load balancer: +$20/month
- Additional storage: +$10/month
- Total prod stack: ~$250-300/month

## Integration Testing

### MCP Tool Discovery

```
What cost analysis tools are available?
```

**Expected Response**: List of MCP tools:
- get_projected_cost
- get_actual_cost
- compare_costs
- analyze_resource_cost
- query_cost_by_tags
- get_recommendations
- detect_anomalies
- forecast_costs
- track_budget

### Streaming Analysis

```
Analyze this entire stack and stream the progress
```

**Expected Response**: Streaming updates showing:
- Analyzed VPC (1/10 resources)
- Analyzed subnet (2/10 resources)
- ... progress updates ...
- Complete analysis with summary
