# Cost Accuracy Validation (SC-006)

## Requirement

Cost estimates must be within ±5% of actual cloud provider billing statements for the same resources over the same time period.

## Validation Methodology

### Phase 1: Test Data Collection

1. **Deploy Reference Infrastructure**
   ```bash
   # Deploy a known Pulumi stack with well-defined resources
   cd test-fixtures/cost-validation-stack
   pulumi up
   ```

2. **Capture Projected Costs**
   ```bash
   # Get projected costs from MCP server
   pulumi preview --json > projected.json

   # Query MCP server for cost estimate
   curl -X POST http://localhost:8080/rpc \
     -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"analyze_projected_costs","arguments":{"pulumi_json":"$(cat projected.json)"}},"id":1}' \
     > projected_estimate.json
   ```

3. **Wait for Billing Period**
   - Run infrastructure for full billing cycle (30 days recommended)
   - Ensure no manual changes to deployed resources
   - Monitor for any unexpected usage patterns

4. **Collect Actual Costs**
   ```bash
   # AWS example
   aws ce get-cost-and-usage \
     --time-period Start=2025-01-01,End=2025-01-31 \
     --granularity MONTHLY \
     --metrics BlendedCost \
     --filter file://resource-filter.json \
     > actual_costs.json

   # Azure example
   az consumption usage list \
     --start-date 2025-01-01 \
     --end-date 2025-01-31 \
     > actual_costs_azure.json
   ```

### Phase 2: Comparison Analysis

```go
// test/validation/cost_accuracy_test.go
func TestCostAccuracy(t *testing.T) {
    // Load projected costs
    projected := loadProjectedCosts("projected_estimate.json")

    // Load actual costs from billing
    actual := loadActualCosts("actual_costs.json")

    // Compare by resource
    for urn, projectedCost := range projected.Resources {
        actualCost := actual.Resources[urn]

        difference := math.Abs(actualCost - projectedCost)
        percentDiff := (difference / actualCost) * 100

        if percentDiff > 5.0 {
            t.Errorf("Resource %s: %.1f%% deviation (projected: $%.2f, actual: $%.2f)",
                urn, percentDiff, projectedCost, actualCost)
        }
    }

    // Compare total stack cost
    totalDiff := math.Abs(actual.Total - projected.Total)
    totalPercentDiff := (totalDiff / actual.Total) * 100

    if totalPercentDiff > 5.0 {
        t.Errorf("Total stack: %.1f%% deviation (projected: $%.2f, actual: $%.2f)",
            totalPercentDiff, projected.Total, actual.Total)
    } else {
        t.Logf("✓ Cost accuracy within ±5%% (deviation: %.2f%%)", totalPercentDiff)
    }
}
```

### Phase 3: Multi-Provider Validation

Test across different cloud providers:

1. **AWS Resources**
   - EC2 instances (on-demand, reserved, spot)
   - RDS databases (various instance types)
   - S3 storage (standard, infrequent access, glacier)
   - Lambda functions
   - ELB/ALB

2. **Azure Resources**
   - Virtual Machines
   - Azure SQL Database
   - Storage Accounts
   - Functions
   - Application Gateway

3. **GCP Resources**
   - Compute Engine
   - Cloud SQL
   - Cloud Storage
   - Cloud Functions
   - Load Balancers

### Phase 4: Edge Cases

Test accuracy under various conditions:

1. **Reserved Instances**
   - Verify RI pricing correctly applied
   - Test RI utilization calculations
   - Validate upfront payment amortization

2. **Spot Instances**
   - Historical spot pricing accuracy
   - Interruption rate impact
   - Fallback to on-demand costs

3. **Data Transfer**
   - Inter-region transfer costs
   - Egress charges
   - CDN costs

4. **Storage Tiers**
   - Lifecycle policy costs
   - Archive storage pricing
   - Request costs (GET, PUT, LIST)

## Acceptance Criteria

- [ ] ≥90% of individual resources within ±5%
- [ ] Total stack cost within ±5%
- [ ] Tested across AWS, Azure, GCP
- [ ] Edge cases documented with known variances
- [ ] Continuous monitoring established

## Known Limitations

1. **Spot Instance Pricing**
   - Highly variable market prices
   - Best-effort estimation only
   - May exceed ±5% during price spikes

2. **Data Transfer Costs**
   - Actual usage patterns vary
   - Geographic distribution affects pricing
   - May require usage-based multipliers

3. **Reserved Instance Allocation**
   - Complex RI sharing rules
   - Account-level RI pool management
   - May show variance in multi-account setups

4. **Free Tier Credits**
   - Not reflected in cost estimates
   - Requires manual adjustment
   - Affects new account validation

## Continuous Validation

```yaml
# .github/workflows/cost-accuracy.yml
name: Monthly Cost Accuracy Check

on:
  schedule:
    - cron: '0 0 1 * *'  # First day of each month
  workflow_dispatch:

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - name: Fetch actual costs
        run: |
          # Query billing APIs

      - name: Compare with projections
        run: |
          go test -v ./test/validation -run TestCostAccuracy

      - name: Report deviations
        if: failure()
        run: |
          # Create GitHub issue if >5% deviation
```

## Validation Report Template

```markdown
## Cost Accuracy Validation Report

**Date**: 2025-01-09
**Period**: 2024-12-01 to 2024-12-31
**Stack**: production-web-app

### Results

| Resource Type | Projected | Actual | Deviation | Status |
|---------------|-----------|--------|-----------|--------|
| EC2 t3.large  | $234.50   | $229.80| -2.0%     | ✓ PASS |
| RDS db.r5.xlarge | $450.00 | $447.20| -0.6%   | ✓ PASS |
| S3 Standard   | $12.30    | $14.50 | +17.9%    | ✗ FAIL |
| **Total**     | **$696.80** | **$691.50** | **-0.8%** | **✓ PASS** |

### Analysis

- **S3 Deviation**: Higher than expected due to increased GET requests not accounted for in static estimate
- **Overall Accuracy**: 99.2% - well within ±5% requirement

### Action Items

- [ ] Improve S3 request cost modeling
- [ ] Add API call volume estimates to S3 projections
```

## References

- AWS Cost Explorer API: https://docs.aws.amazon.com/cost-management/
- Azure Cost Management: https://docs.microsoft.com/azure/cost-management-billing/
- GCP Billing API: https://cloud.google.com/billing/docs/apis
