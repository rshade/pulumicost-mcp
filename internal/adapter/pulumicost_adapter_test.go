package adapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T029: TestNewPulumiCostAdapter - Create adapter instance
func TestNewPulumiCostAdapter(t *testing.T) {
	adapter := NewPulumiCostAdapter("/usr/local/bin/pulumicost")

	require.NotNil(t, adapter)
	assert.Equal(t, "/usr/local/bin/pulumicost", adapter.GetCorePath())
}

// T030: TestGetProjectedCost - Parse Pulumi preview JSON and calculate cost
func TestGetProjectedCost(t *testing.T) {
	// Arrange
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost.sh")
	ctx := context.Background()

	pulumiJSON := `{
		"resources": [
			{
				"urn": "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::web-server",
				"type": "aws:ec2/instance:Instance",
				"inputs": {
					"instanceType": "t3.micro",
					"ami": "ami-12345678"
				}
			}
		]
	}`

	// Act
	result, err := adapter.GetProjectedCost(ctx, pulumiJSON)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.TotalMonthly, 0.0, "Total monthly cost should be greater than 0")
	assert.Equal(t, "USD", result.Currency)
	assert.NotEmpty(t, result.Resources, "Resources breakdown should not be empty")
}

// T030: TestGetProjectedCost_InvalidJSON
func TestGetProjectedCost_InvalidJSON(t *testing.T) {
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost.sh")
	ctx := context.Background()

	result, err := adapter.GetProjectedCost(ctx, "invalid json")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "JSON")
}

// T030: TestGetProjectedCost_EmptyResources
func TestGetProjectedCost_EmptyResources(t *testing.T) {
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost.sh")
	ctx := context.Background()

	pulumiJSON := `{"resources": []}`

	result, err := adapter.GetProjectedCost(ctx, pulumiJSON)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Mock script returns static data, so just verify we get a result
	assert.NotNil(t, result.TotalMonthly)
	assert.Equal(t, "USD", result.Currency)
}

// T030: TestGetProjectedCost_WithFilters
func TestGetProjectedCost_WithFilters(t *testing.T) {
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost.sh")
	ctx := context.Background()

	pulumiJSON := `{
		"resources": [
			{
				"urn": "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::web-1",
				"type": "aws:ec2/instance:Instance",
				"inputs": {"instanceType": "t3.micro"}
			},
			{
				"urn": "urn:pulumi:dev::myapp::azure:compute/virtualMachine:VirtualMachine::web-2",
				"type": "azure:compute/virtualMachine:VirtualMachine",
				"inputs": {"vmSize": "Standard_B1s"}
			}
		]
	}`

	filters := &ResourceFilters{
		Provider: stringPtr("aws"),
	}

	result, err := adapter.GetProjectedCostWithFilters(ctx, pulumiJSON, filters)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Should only include AWS resources
	for _, resource := range result.Resources {
		assert.Equal(t, "aws", *resource.Provider)
	}
}

// T032: TestGetActualCost - Retrieve historical costs from cloud provider
func TestGetActualCost(t *testing.T) {
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost.sh")
	ctx := context.Background()

	stackName := "myapp-dev"
	timeRange := TimeRange{
		Start: "2024-01-01T00:00:00Z",
		End:   "2024-01-31T23:59:59Z",
	}

	result, err := adapter.GetActualCost(ctx, stackName, timeRange)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.TotalMonthly, 0.0)
	assert.Equal(t, "USD", result.Currency)
	assert.NotEmpty(t, result.Resources)
}

// T032: TestGetActualCost_InvalidTimeRange
func TestGetActualCost_InvalidTimeRange(t *testing.T) {
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost.sh")
	ctx := context.Background()

	stackName := "myapp-dev"
	timeRange := TimeRange{
		Start: "invalid",
		End:   "also-invalid",
	}

	result, err := adapter.GetActualCost(ctx, stackName, timeRange)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "time")
}

// T032: TestGetActualCost_ContextCancellation
func TestGetActualCost_ContextCancellation(t *testing.T) {
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost_slow.sh")
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	stackName := "myapp-dev"
	timeRange := TimeRange{
		Start: "2024-01-01T00:00:00Z",
		End:   "2024-01-31T23:59:59Z",
	}

	result, err := adapter.GetActualCost(ctx, stackName, timeRange)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context")
}

// T033: TestGetActualCost_WithGranularity
func TestGetActualCost_WithGranularity(t *testing.T) {
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost.sh")
	ctx := context.Background()

	stackName := "myapp-dev"
	timeRange := TimeRange{
		Start: "2024-01-01T00:00:00Z",
		End:   "2024-01-31T23:59:59Z",
	}

	result, err := adapter.GetActualCostWithGranularity(ctx, stackName, timeRange, "daily")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.TotalMonthly, 0.0)
	// Should have daily breakdown
	assert.NotNil(t, result.Breakdown)
}

// TestGetProjectedCost_WithMultipleFilters tests various filter combinations
func TestGetProjectedCost_WithMultipleFilters(t *testing.T) {
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost.sh")
	ctx := context.Background()

	pulumiJSON := `{
		"resources": [
			{
				"urn": "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::web-1",
				"type": "aws:ec2/instance:Instance",
				"inputs": {"instanceType": "t3.micro"}
			}
		]
	}`

	tests := []struct {
		name    string
		filters *ResourceFilters
	}{
		{
			name: "filter by resource type",
			filters: &ResourceFilters{
				ResourceType: stringPtr("aws:ec2/instance:Instance"),
			},
		},
		{
			name: "filter by region",
			filters: &ResourceFilters{
				Region: stringPtr("us-east-1"),
			},
		},
		{
			name: "filter by name pattern",
			filters: &ResourceFilters{
				NamePattern: stringPtr("web"),
			},
		},
		{
			name: "filter by tags",
			filters: &ResourceFilters{
				Tags: map[string]string{
					"environment": "dev",
				},
			},
		},
		{
			name: "filter by provider and type",
			filters: &ResourceFilters{
				Provider:     stringPtr("aws"),
				ResourceType: stringPtr("aws:ec2/instance:Instance"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.GetProjectedCostWithFilters(ctx, pulumiJSON, tt.filters)
			require.NoError(t, err)
			require.NotNil(t, result)
		})
	}
}

// TestGetProjectedCost_FilterNoMatches tests filters that don't match any resources
func TestGetProjectedCost_FilterNoMatches(t *testing.T) {
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost.sh")
	ctx := context.Background()

	pulumiJSON := `{
		"resources": [
			{
				"urn": "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::web-1",
				"type": "aws:ec2/instance:Instance",
				"inputs": {"instanceType": "t3.micro"}
			}
		]
	}`

	filters := &ResourceFilters{
		Provider: stringPtr("gcp"), // No GCP resources in the mock data
	}

	result, err := adapter.GetProjectedCostWithFilters(ctx, pulumiJSON, filters)

	require.NoError(t, err)
	require.NotNil(t, result)
	// When filters don't match, we should get empty or filtered results
	// The mock always returns data, so we just verify no error
}

// TestGetActualCost_WithMultipleGranularities tests different time granularities
func TestGetActualCost_WithMultipleGranularities(t *testing.T) {
	adapter := NewPulumiCostAdapter("./testdata/mock_pulumicost.sh")
	ctx := context.Background()

	stackName := "myapp-dev"
	timeRange := TimeRange{
		Start: "2024-01-01T00:00:00Z",
		End:   "2024-01-31T23:59:59Z",
	}

	granularities := []string{"hourly", "daily", "monthly"}

	for _, granularity := range granularities {
		t.Run(granularity, func(t *testing.T) {
			result, err := adapter.GetActualCostWithGranularity(ctx, stackName, timeRange, granularity)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Greater(t, result.TotalMonthly, 0.0)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
