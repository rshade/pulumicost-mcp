# Cost Optimization & Analytics Query Examples

This document provides example questions for AI-powered cost optimization, anomaly detection, forecasting, and budget tracking.

## Cost Optimization Recommendations

1. **General Recommendations**
   ```
   What are cost optimization recommendations for my infrastructure?
   ```

2. **Rightsizing Opportunities**
   ```
   Which resources are over-provisioned and could be downsized?
   ```

3. **Reserved Instance Recommendations**
   ```
   Should I purchase reserved instances? Where would I save the most?
   ```

4. **Spot Instance Opportunities**
   ```
   Which workloads could run on spot instances to reduce costs?
   ```

5. **Storage Optimization**
   ```
   What storage optimization opportunities exist (archival, lifecycle policies)?
   ```

6. **High-Confidence Recommendations**
   ```
   Show me only HIGH confidence optimization recommendations.
   ```

7. **Savings Potential**
   ```
   What's my total potential savings from all recommendations?
   ```

## Anomaly Detection

8. **Recent Anomalies**
   ```
   Show me cost anomalies from the last 30 days.
   ```

9. **Critical Anomalies**
   ```
   Are there any CRITICAL cost anomalies I should investigate immediately?
   ```

10. **Anomaly Root Causes**
    ```
    What caused the cost spike on January 15th?
    ```

11. **Resource-Specific Anomalies**
    ```
    Has this EC2 instance shown any unusual cost patterns?
    ```

12. **Trend Analysis**
    ```
    Are my costs trending upward abnormally?
    ```

13. **Baseline Comparison**
    ```
    How does this month's spending compare to my 30-day baseline?
    ```

## Cost Forecasting

14. **Monthly Forecast**
    ```
    What are my projected costs for next month?
    ```

15. **Quarterly Forecast**
    ```
    Forecast my cloud spending for Q2 2024 with confidence intervals.
    ```

16. **Trend-Based Projection**
    ```
    If my current growth rate continues, what will costs look like in 6 months?
    ```

17. **Scenario Planning**
    ```
    What would happen to my forecast if I add these new resources?
    ```

18. **Confidence Intervals**
    ```
    Give me 90% and 99% confidence intervals for next month's costs.
    ```

## Budget Tracking

19. **Budget Status**
    ```
    How am I tracking against my $5000/month budget?
    ```

20. **Burn Rate**
    ```
    What's my current daily burn rate? When will I exhaust my budget?
    ```

21. **Budget Alerts**
    ```
    Have I crossed any budget threshold alerts?
    ```

22. **Monthly Budget Health**
    ```
    Am I on track to stay within budget this month?
    ```

23. **Remaining Budget**
    ```
    How much of my budget is remaining?
    ```

24. **Budget vs Actual**
    ```
    Compare my actual spending to my planned budget for each month.
    ```

25. **Alert Threshold Setup**
    ```
    Set budget alerts at 50%, 80%, and 100% of my $10K monthly budget.
    ```

## Combined Analytics

26. **Comprehensive Cost Health**
    ```
    Give me a complete cost health report: recommendations, anomalies, forecast, and budget status.
    ```

27. **Cost Intelligence Dashboard**
    ```
    Show me the most important cost insights right now.
    ```

28. **Executive Summary**
    ```
    Summarize my cloud costs for an executive report: trends, risks, opportunities.
    ```

## Example Recommendation Response

When you ask "What are cost optimization recommendations?", you might get:

```
Cost Optimization Recommendations

Total Potential Savings: $847/month (23% reduction)

HIGH CONFIDENCE:
1. Rightsize EC2 Instance: web-server-prod
   Current: t3.large ($60.74/month)
   Recommended: t3.medium ($30.37/month)
   Savings: $30.37/month
   Reason: CPU utilization consistently <25% for 90 days

2. Reserved Instance: rds-postgres-prod
   Current: On-Demand ($175/month)
   Recommended: 1-year Reserved Instance ($105/month)
   Savings: $70/month
   Reason: 180+ days continuous uptime, stable usage pattern

MEDIUM CONFIDENCE:
3. Move to Spot: batch-processor
   Current: t3.xlarge On-Demand ($121.47/month)
   Recommended: t3.xlarge Spot ($36.44/month)
   Savings: $85.03/month
   Reason: Stateless batch workload, fault-tolerant

4. S3 Lifecycle Policy: logs-bucket
   Savings: ~$120/month
   Reason: 85% of objects not accessed in 90+ days
```

## Example Anomaly Response

When you ask "Show me cost anomalies from the last 30 days", you might get:

```
Cost Anomalies Detected (Last 30 Days)

CRITICAL:
- January 15, 2024: +$1,247 (+156% above baseline)
  Affected Resources:
    - urn:...:ec2-autoscaling-group (new instances launched)
    - urn:...:rds-instance (increased IOPS)
  Potential Cause: Traffic spike, autoscaling event
  Severity: CRITICAL

HIGH:
- January 22, 2024: +$385 (+48% above baseline)
  Affected Resources:
    - urn:...:s3-data-transfer
  Potential Cause: Large data egress, region transfer
  Severity: HIGH

Current Baseline: $800/day (30-day rolling average)
Standard Deviation: $120/day
```

## Example Forecast Response

When you ask "Forecast my costs for next month", you might get:

```
Cost Forecast: February 2024

Predicted Cost: $24,500
Confidence Level: 95%

Confidence Intervals:
- 90% CI: $23,100 - $25,900
- 95% CI: $22,800 - $26,200
- 99% CI: $22,100 - $26,900

Methodology: Linear regression on 90-day historical data
Trend: Costs increasing 8% month-over-month
Seasonality: Slight weekly pattern detected (higher on weekdays)

Assumptions:
- No major architecture changes
- Current usage patterns continue
- No reserved instance purchases

If current growth rate continues: $31,200 by June 2024
```

## Example Budget Status

When you ask "How am I tracking against budget?", you might get:

```
Budget Status: January 2024

Budget: $10,000/month
Current Spending: $7,850 (78.5%)
Remaining: $2,150 (21.5%)
Days Remaining: 7 days

Burn Rate: $350/day
Projected End Date: January 29 (2 days before month end)

Status: ⚠️ WARNING

Alert Thresholds:
✓ 50% threshold passed (Jan 15)
✓ 80% threshold passed (Jan 25) ← Current
○ 100% threshold (projected Jan 29)

Recommendation: You're on track to exceed budget by ~$350.
Consider implementing cost optimization recommendations to stay within budget.
```

## Tips for Analytics Queries

- **Be specific with time ranges**: "last 30 days", "Q1 2024", etc.
- **Combine analytics**: Ask for recommendations + anomalies in one query
- **Set thresholds**: Specify confidence levels (HIGH only, etc.)
- **Track trends**: Regular forecasts help predict budget issues
- **Act on recommendations**: Implement HIGH confidence recommendations first
- **Investigate anomalies**: Don't ignore cost spikes - find root causes
