# Example Pulumi Stacks

This directory contains example Pulumi projects that you can use to test the PulumiCost MCP Server's cost analysis capabilities.

## Available Examples

### 1. Simple AWS Stack (`simple-aws/`)

A basic AWS infrastructure with common resources:
- EC2 instance (t3.micro)
- S3 bucket
- RDS PostgreSQL instance (db.t3.small)
- VPC with subnets

**Estimated Monthly Cost**: ~$50-75

**Use Cases**:
- Testing basic cost projection
- Learning Pulumi syntax
- Quick validation of MCP server

### 2. Multi-Tier Application (`multi-tier-app/`)

A production-like 3-tier application:
- Application Load Balancer
- Auto Scaling Group (3x t3.medium instances)
- RDS Multi-AZ PostgreSQL (db.m5.large)
- ElastiCache Redis cluster
- S3 buckets for static assets and logs
- CloudFront distribution

**Estimated Monthly Cost**: ~$400-500

**Use Cases**:
- Testing cost breakdown by resource type
- Comparing costs across regions
- Optimization recommendations

### 3. Kubernetes Cluster (`k8s-cluster/`)

EKS cluster with supporting resources:
- EKS control plane
- Managed node groups (3x t3.large)
- VPC with private/public subnets
- NAT Gateways
- Load balancers

**Estimated Monthly Cost**: ~$200-250

**Use Cases**:
- Kubernetes cost analysis
- Comparing container vs VM costs
- Network cost analysis

### 4. Serverless Application (`serverless-app/`)

Event-driven serverless architecture:
- Lambda functions (10+ functions)
- API Gateway
- DynamoDB tables
- S3 buckets with lifecycle policies
- EventBridge rules
- SQS queues

**Estimated Monthly Cost**: ~$10-30 (highly variable based on usage)

**Use Cases**:
- Pay-per-use cost analysis
- Comparing serverless vs traditional
- Usage-based forecasting

## Quick Start

### Prerequisites

- Pulumi CLI installed
- AWS credentials configured
- PulumiCost MCP Server running

### Deploy an Example Stack

```bash
# Navigate to an example
cd examples/pulumi-stacks/simple-aws

# Install dependencies
npm install

# Preview the deployment
pulumi preview --json > preview.json

# Get projected costs
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d "{
    \"jsonrpc\": \"2.0\",
    \"method\": \"cost.get_projected_cost\",
    \"params\": {
      \"preview_data\": $(cat preview.json)
    },
    \"id\": 1
  }"

# Deploy the stack
pulumi up

# Export stack state for actual cost tracking
pulumi stack export > stack-state.json

# Get actual costs (after resources are deployed)
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d "{
    \"jsonrpc\": \"2.0\",
    \"method\": \"cost.get_actual_cost\",
    \"params\": {
      \"stack_name\": \"dev\"
    },
    \"id\": 1
  }"
```

### Using with Claude Desktop

With PulumiCost MCP server configured in Claude Desktop, you can ask:

```
I'm about to deploy this Pulumi stack (paste preview.json).
What will it cost me per month?
```

```
Compare the costs of my dev and prod stacks
```

```
What optimization opportunities exist in my infrastructure?
```

## Example Queries

Each example directory contains a `queries.md` file with suggested cost analysis queries specific to that stack architecture.

### General Queries

- "What's the total projected cost for this stack?"
- "Break down costs by resource type"
- "Show me the top 5 most expensive resources"
- "Compare costs between us-east-1 and us-west-2"
- "What would happen to costs if I scaled up by 50%?"
- "Recommend cost optimizations"

### Stack-Specific Queries

See `queries.md` in each example directory for architecture-specific questions.

## Customizing Examples

Feel free to modify these examples:

1. **Change regions**: Update `region` in Pulumi config
2. **Scale resources**: Modify instance counts, sizes
3. **Add resources**: Extend with additional AWS services
4. **Multi-cloud**: Add Azure or GCP resources

```bash
# Change region
pulumi config set aws:region us-west-2

# Change instance type
pulumi config set instanceType t3.small

# Preview cost impact
pulumi preview --json | pulumicost analyze
```

## Cost Estimation Accuracy

**Important Notes**:

- Costs are estimates based on AWS pricing as of the last update
- Actual costs may vary based on:
  - Data transfer volumes
  - Request counts
  - Storage growth
  - Reserved Instance or Savings Plans usage
  - Regional pricing differences
- Always check actual costs in AWS Cost Explorer

## Cleaning Up

To avoid ongoing charges, destroy stacks when done testing:

```bash
# Destroy the stack
pulumi destroy

# Remove the stack from Pulumi state
pulumi stack rm dev
```

## Contributing Examples

Have a useful example stack? Contributions welcome!

1. Create a new directory under `examples/pulumi-stacks/`
2. Include:
   - `README.md` with description and cost estimate
   - `queries.md` with example analysis queries
   - Complete Pulumi project files
   - Example preview output
3. Submit a PR

See `CONTRIBUTING.md` for guidelines.

## Support

- **Issues**: Report problems with examples in GitHub Issues
- **Questions**: Ask in GitHub Discussions
- **Pulumi Help**: See [Pulumi Documentation](https://www.pulumi.com/docs/)
