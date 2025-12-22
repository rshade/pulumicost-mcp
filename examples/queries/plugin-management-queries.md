# Plugin Management Query Examples

This document provides example questions for managing and validating cost source plugins.

## Plugin Discovery

1. **List All Plugins**
   ```
   What cost plugins are available?
   ```

2. **Plugin Status**
   ```
   Show me the health status of all installed plugins.
   ```

3. **Healthy Plugins Only**
   ```
   Which plugins are currently healthy and ready to use?
   ```

## Plugin Information

4. **Plugin Details**
   ```
   Tell me about the kubecost plugin - what can it do?
   ```

5. **Plugin Capabilities**
   ```
   What cloud providers does the AWS Cost Explorer plugin support?
   ```

6. **Plugin Version**
   ```
   What version of the infracost plugin is installed?
   ```

7. **Supported Resources**
   ```
   Which resource types can the GCP billing plugin analyze?
   ```

## Plugin Validation

8. **Conformance Testing**
   ```
   Validate the kubecost plugin against the pulumicost-spec.
   ```

9. **Detailed Validation Report**
   ```
   Run conformance tests on the AWS plugin and show me the detailed results.
   ```

10. **Multiple Plugin Validation**
    ```
    Validate all installed plugins and show which ones pass conformance tests.
    ```

11. **Validation Level**
    ```
    What conformance level does the Azure plugin achieve (BASIC, STANDARD, FULL)?
    ```

## Plugin Health Monitoring

12. **Health Check**
    ```
    Check the health of the infracost plugin and report latency.
    ```

13. **Plugin Connectivity**
    ```
    Are all plugins responding? Show me any connection issues.
    ```

14. **Latency Monitoring**
    ```
    Which plugins have the lowest response time?
    ```

15. **Troubleshooting**
    ```
    The kubecost plugin isn't working - what's wrong?
    ```

## Plugin Comparison

16. **Feature Comparison**
    ```
    Compare the capabilities of the AWS and Azure cost plugins.
    ```

17. **Provider Coverage**
    ```
    Which plugins support Kubernetes cost analysis?
    ```

18. **Best Plugin for Use Case**
    ```
    I need historical cost data for AWS - which plugin should I use?
    ```

## Plugin Configuration

19. **Configuration Requirements**
    ```
    What configuration does the kubecost plugin need?
    ```

20. **Setup Instructions**
    ```
    How do I configure the AWS Cost Explorer plugin?
    ```

## Example Plugin Discovery Response

When you ask "What cost plugins are available?", you might get:

```
Available Cost Plugins:

1. **infracost** (v0.10.30)
   - Status: HEALTHY
   - Latency: 45ms
   - Supports: AWS, Azure, GCP
   - Capabilities: Projected costs, resource-level pricing

2. **kubecost** (v1.108.0)
   - Status: HEALTHY
   - Latency: 120ms
   - Supports: Kubernetes
   - Capabilities: Actual costs, pod-level attribution

3. **aws-cost-explorer** (v2.1.0)
   - Status: UNHEALTHY (connection timeout)
   - Last Error: Failed to reach AWS Cost Explorer API
   - Supports: AWS
   - Capabilities: Actual historical costs, forecasting
```

## Plugin Validation Example

When you ask "Validate the kubecost plugin", you might get:

```
Kubecost Plugin Validation Report

Conformance Level: STANDARD
Overall Status: PASS

Test Results:
✓ GetActualCost - PASS (125ms)
✓ Health Check - PASS (45ms)
✓ Metadata Response - PASS
✓ Error Handling - PASS
✓ Timeout Handling - PASS
✗ GetProjectedCost - FAIL (not supported)

The kubecost plugin successfully implements STANDARD conformance level.
It supports actual cost retrieval but not projected cost estimation.
```

## Tips for Plugin Management

- **Regular Health Checks**: Monitor plugin status periodically
- **Validate After Updates**: Run conformance tests when plugins are updated
- **Check Capabilities**: Verify plugins support your required cloud providers
- **Monitor Latency**: Watch for performance degradation
- **Read Validation Reports**: Understand what features each plugin supports
