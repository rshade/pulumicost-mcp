package service

import (
	"context"
	"testing"

	"github.com/rshade/pulumicost-mcp/internal/adapter"
	cost "github.com/rshade/pulumicost-mcp/gen/cost"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T023: TestAnalyzeProjected - RED test for FR-001
func TestAnalyzeProjected(t *testing.T) {
	// Arrange
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.AnalyzeProjectedPayload{
		PulumiJSON: `{
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
		}`,
	}

	// Act
	result, err := service.AnalyzeProjected(ctx, payload)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.TotalMonthly, 0.0, "Total monthly cost should be greater than 0")
	assert.Equal(t, "USD", result.Currency)
	assert.NotEmpty(t, result.Resources, "Resources breakdown should not be empty")
}

// T023: TestAnalyzeProjected_InvalidJSON
func TestAnalyzeProjected_InvalidJSON(t *testing.T) {
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.AnalyzeProjectedPayload{
		PulumiJSON: "invalid json",
	}

	result, err := service.AnalyzeProjected(ctx, payload)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// T023: TestAnalyzeProjected_WithFilters
func TestAnalyzeProjected_WithFilters(t *testing.T) {
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	provider := "aws"
	payload := &cost.AnalyzeProjectedPayload{
		PulumiJSON: `{"resources": []}`,
		Filters: &cost.ResourceFilter{
			Provider: &provider,
		},
	}

	result, err := service.AnalyzeProjected(ctx, payload)

	require.NoError(t, err)
	assert.NotNil(t, result)
}

// T024: TestGetActual - RED test for FR-002
func TestGetActual(t *testing.T) {
	// Arrange
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.GetActualPayload{
		StackName: "myapp-dev",
		TimeRange: &cost.TimeRange{
			Start: "2024-01-01T00:00:00Z",
			End:   "2024-01-31T23:59:59Z",
		},
		Granularity: stringPtr("daily"),
	}

	// Act
	result, err := service.GetActual(ctx, payload)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Greater(t, result.TotalMonthly, 0.0)
	assert.Equal(t, "USD", result.Currency)
	assert.NotEmpty(t, result.Resources)
}

// T024: TestGetActual_InvalidTimeRange
func TestGetActual_InvalidTimeRange(t *testing.T) {
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.GetActualPayload{
		StackName: "myapp-dev",
		TimeRange: &cost.TimeRange{
			Start: "invalid",
			End:   "also-invalid",
		},
	}

	result, err := service.GetActual(ctx, payload)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// T025: TestCompareCosts - RED test for FR-003
func TestCompareCosts(t *testing.T) {
	// Arrange
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.CompareCostsPayload{
		Baseline: &struct {
			StackName  *string
			PulumiJSON *string
			Filters    *cost.ResourceFilter
		}{
			StackName: stringPtr("myapp-dev"),
		},
		Target: &struct {
			StackName  *string
			PulumiJSON *string
			Filters    *cost.ResourceFilter
		}{
			StackName: stringPtr("myapp-prod"),
		},
		ComparisonType: stringPtr("both"),
	}

	// Act
	result, err := service.CompareCosts(ctx, payload)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.BaselineCost)
	assert.NotNil(t, result.TargetCost)
	assert.NotNil(t, result.Difference)
	assert.NotNil(t, result.DifferencePercent)
}

// T025: TestCompareCosts_WithPulumiJSON
func TestCompareCosts_WithPulumiJSON(t *testing.T) {
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	baselineJSON := `{"resources": []}`
	targetJSON := `{"resources": []}`

	payload := &cost.CompareCostsPayload{
		Baseline: &struct {
			StackName  *string
			PulumiJSON *string
			Filters    *cost.ResourceFilter
		}{
			PulumiJSON: &baselineJSON,
		},
		Target: &struct {
			StackName  *string
			PulumiJSON *string
			Filters    *cost.ResourceFilter
		}{
			PulumiJSON: &targetJSON,
		},
	}

	result, err := service.CompareCosts(ctx, payload)

	require.NoError(t, err)
	assert.NotNil(t, result)
}

// T026: TestAnalyzeResource - RED test for FR-004
func TestAnalyzeResource(t *testing.T) {
	// Arrange
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.AnalyzeResourcePayload{
		ResourceUrn:         "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::web-server",
		IncludeDependencies: true,
	}

	// Act
	result, err := service.AnalyzeResource(ctx, payload)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.Resource)
	assert.Equal(t, payload.ResourceUrn, result.Resource.Urn)
	assert.Greater(t, result.Resource.MonthlyCost, 0.0)
}

// T026: TestAnalyzeResource_WithDependencies
func TestAnalyzeResource_WithDependencies(t *testing.T) {
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.AnalyzeResourcePayload{
		ResourceUrn:         "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::web-server",
		IncludeDependencies: true,
	}

	result, err := service.AnalyzeResource(ctx, payload)

	require.NoError(t, err)
	assert.NotNil(t, result)
	// Dependencies may be empty or populated depending on the resource
}

// T027: TestQueryByTags - RED test for FR-005
func TestQueryByTags(t *testing.T) {
	// Arrange
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.QueryByTagsPayload{
		StackName: "myapp-dev",
		TagKeys:   []string{"environment", "team"},
	}

	// Act
	result, err := service.QueryByTags(ctx, payload)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.ByTag)
	assert.NotEmpty(t, result.ByTag, "Should have tag-based cost groupings")
}

// T027: TestQueryByTags_WithFilters
func TestQueryByTags_WithFilters(t *testing.T) {
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.QueryByTagsPayload{
		StackName: "myapp-dev",
		TagKeys:   []string{"environment"},
		Filters: &struct {
			Key    *string
			Values []string
		}{
			Key:    stringPtr("team"),
			Values: []string{"platform", "backend"},
		},
	}

	result, err := service.QueryByTags(ctx, payload)

	require.NoError(t, err)
	assert.NotNil(t, result)
}

// T028: TestAnalyzeStack - RED test for FR-006 (streaming)
func TestAnalyzeStack(t *testing.T) {
	// Arrange
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.AnalyzeStackPayload{
		StackName:              "myapp-dev",
		IncludeRecommendations: false,
	}

	// Mock stream that captures events
	mockStream := &mockAnalyzeStackStream{
		events: make([]cost.AnalyzeStackEvent, 0),
	}

	// Act
	err := service.AnalyzeStack(ctx, payload, mockStream)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, mockStream.events, "Should have received streaming events")

	// Verify we got progress updates
	hasProgress := false
	for _, event := range mockStream.events {
		if result, ok := event.(*cost.AnalyzeStackResult); ok {
			if result.Progress != nil {
				hasProgress = true
				assert.GreaterOrEqual(t, *result.Progress, 0.0)
				assert.LessOrEqual(t, *result.Progress, 100.0)
			}
		}
	}
	assert.True(t, hasProgress, "Should have received progress updates")
}

// T028: TestAnalyzeStack_WithRecommendations
func TestAnalyzeStack_WithRecommendations(t *testing.T) {
	mockAdapter := adapter.NewPulumiCostAdapter("../adapter/testdata/mock_pulumicost.sh")
	service := NewCostService(mockAdapter, nil)
	ctx := context.Background()

	payload := &cost.AnalyzeStackPayload{
		StackName:              "myapp-dev",
		IncludeRecommendations: true,
	}

	mockStream := &mockAnalyzeStackStream{
		events: make([]cost.AnalyzeStackEvent, 0),
	}

	err := service.AnalyzeStack(ctx, payload, mockStream)

	require.NoError(t, err)
	assert.NotEmpty(t, mockStream.events)
}

// Mock stream for testing
type mockAnalyzeStackStream struct {
	events []cost.AnalyzeStackEvent
}

func (m *mockAnalyzeStackStream) Send(ctx context.Context, event cost.AnalyzeStackEvent) error {
	m.events = append(m.events, event)
	return nil
}

func (m *mockAnalyzeStackStream) SendAndClose(ctx context.Context, event cost.AnalyzeStackEvent) error {
	m.events = append(m.events, event)
	return nil
}

func (m *mockAnalyzeStackStream) SendError(ctx context.Context, id string, err error) error {
	return err
}
