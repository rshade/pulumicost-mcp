# Cost Analysis Query Examples

This document provides example questions you can ask Claude (or another AI assistant) when the PulumiCost MCP server is configured.

## Projected Cost Queries

1. **Basic Projected Cost**
   ```
   What are the projected costs for this Pulumi stack?
   [Provide Pulumi preview JSON]
   ```

2. **Cost Breakdown by Service**
   ```
   Show me the projected costs broken down by AWS service for my infrastructure.
   ```

3. **Provider Comparison**
   ```
   What would it cost to run this infrastructure on AWS vs Azure?
   ```

4. **Regional Cost Comparison**
   ```
   How much would it cost to deploy this in us-east-1 vs eu-west-1?
   ```

5. **Filtered Cost Analysis**
   ```
   What are the projected costs for only compute resources in my stack?
   ```

## Actual Cost Queries

6. **Historical Costs**
   ```
   What were my actual cloud costs for the last 30 days?
   ```

7. **Daily Cost Trend**
   ```
   Show me daily cost trends for January 2024.
   ```

8. **Monthly Comparison**
   ```
   Compare my cloud costs from last month to this month.
   ```

9. **Granular Cost Analysis**
   ```
   Break down my costs by hour for the last week.
   ```

## Cost Comparison

10. **Configuration Comparison**
    ```
    Compare the costs between my current infrastructure and this proposed change.
    [Provide two configurations]
    ```

11. **Percentage Change**
    ```
    What's the percentage increase if I upgrade from t3.micro to t3.medium instances?
    ```

12. **Baseline Comparison**
    ```
    How much more expensive is my production stack compared to staging?
    ```

## Resource-Specific Analysis

13. **Individual Resource Cost**
    ```
    How much does this specific EC2 instance cost per month?
    [Provide resource URN]
    ```

14. **Resource Dependencies**
    ```
    What are the costs of this RDS instance including its dependent resources?
    ```

15. **Cost Per Resource Type**
    ```
    How much am I spending on S3 storage vs EC2 compute?
    ```

## Tag-Based Queries

16. **Environment Costs**
    ```
    What are my costs grouped by environment tag (dev, staging, prod)?
    ```

17. **Team Attribution**
    ```
    Show me costs by team tag to understand each team's cloud spending.
    ```

18. **Project Costs**
    ```
    Break down costs by project tag for budget allocation.
    ```

## Comprehensive Analysis

19. **Full Stack Analysis**
    ```
    Analyze my entire Pulumi stack and give me a comprehensive cost report with trends.
    ```

20. **Large Stack with Progress**
    ```
    Analyze this 500+ resource stack and stream progress updates as you go.
    ```

## Example Pulumi JSON

For testing queries, you can use this sample Pulumi preview JSON:

```json
{
  "resources": [
    {
      "urn": "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::web-server",
      "type": "aws:ec2/instance:Instance",
      "inputs": {
        "instanceType": "t3.micro",
        "ami": "ami-12345678"
      }
    },
    {
      "urn": "urn:pulumi:dev::myapp::aws:rds/instance:Instance::db",
      "type": "aws:rds/instance:Instance",
      "inputs": {
        "instanceClass": "db.t3.micro",
        "engine": "postgres"
      }
    },
    {
      "urn": "urn:pulumi:dev::myapp::aws:s3/bucket:Bucket::assets",
      "type": "aws:s3/bucket:Bucket",
      "inputs": {}
    }
  ]
}
```

## Tips for Better Results

- **Be specific**: Include resource URNs, time ranges, or tag filters when possible
- **Use context**: Provide Pulumi JSON for projected costs, stack names for actual costs
- **Ask follow-ups**: Claude can drill down into details after initial analysis
- **Compare alternatives**: Ask about cost implications of different configurations
- **Check trends**: Request time-series data to understand cost patterns
