package service

import (
	"context"
	"testing"

	"github.com/rshade/pulumicost-mcp/gen/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetRecommendations tests getting cost optimization recommendations
func TestGetRecommendations(t *testing.T) {
	service := NewAnalysisService(nil, nil)
	ctx := context.Background()

	payload := &analysis.GetRecommendationsPayload{
		StackName: "my-stack",
	}

	result, err := service.GetRecommendations(ctx, payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Recommendations)

	// Verify recommendation structure
	rec := result.Recommendations[0]
	assert.NotEmpty(t, rec.ID)
	assert.NotEmpty(t, rec.Type)
	assert.Greater(t, rec.ProjectedSavings, 0.0)
	assert.NotEmpty(t, rec.Description)
}

// TestGetRecommendations_WithFilters tests filtering recommendations
func TestGetRecommendations_WithFilters(t *testing.T) {
	service := NewAnalysisService(nil, nil)
	ctx := context.Background()

	minSavings := 100.0
	payload := &analysis.GetRecommendationsPayload{
		StackName:           "my-stack",
		RecommendationTypes: []string{"RIGHTSIZING"},
		MinimumSavings:      &minSavings,
	}

	result, err := service.GetRecommendations(ctx, payload)

	require.NoError(t, err)
	assert.NotEmpty(t, result.Recommendations)
}

// TestDetectAnomalies tests anomaly detection
func TestDetectAnomalies(t *testing.T) {
	service := NewAnalysisService(nil, nil)
	ctx := context.Background()

	payload := &analysis.DetectAnomaliesPayload{
		StackName: "my-stack",
		TimeRange: &analysis.TimeRange{
			Start: "2024-01-01T00:00:00Z",
			End:   "2024-01-31T23:59:59Z",
		},
		Sensitivity: "MEDIUM",
	}

	result, err := service.DetectAnomalies(ctx, payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, result.Anomalies)
}

// TestForecast tests cost forecasting
func TestForecast(t *testing.T) {
	service := NewAnalysisService(nil, nil)
	ctx := context.Background()

	payload := &analysis.ForecastPayload{
		StackName: "my-stack",
		ForecastPeriod: &analysis.TimeRange{
			Start: "2024-02-01T00:00:00Z",
			End:   "2024-02-29T23:59:59Z",
		},
		ConfidenceLevel: 0.95,
	}

	result, err := service.Forecast(ctx, payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "my-stack", result.StackName)
	assert.NotEmpty(t, result.DataPoints)
	assert.Greater(t, result.DataPoints[0].PredictedCost, 0.0)
}

// TestTrackBudget tests budget tracking
func TestTrackBudget(t *testing.T) {
	service := NewAnalysisService(nil, nil)
	ctx := context.Background()

	payload := &analysis.TrackBudgetPayload{
		StackName:    "my-stack",
		BudgetAmount: 1000.0,
		Period:       "MONTHLY",
	}

	result, err := service.TrackBudget(ctx, payload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1000.0, result.BudgetAmount)
	assert.Greater(t, result.CurrentSpending, 0.0)
	assert.NotEmpty(t, result.Status)
}

// TestTrackBudget_WithAlerts tests budget with alert thresholds
func TestTrackBudget_WithAlerts(t *testing.T) {
	service := NewAnalysisService(nil, nil)
	ctx := context.Background()

	payload := &analysis.TrackBudgetPayload{
		StackName:       "my-stack",
		BudgetAmount:    1000.0,
		Period:          "MONTHLY",
		AlertThresholds: []float64{50.0, 80.0, 100.0},
	}

	result, err := service.TrackBudget(ctx, payload)

	require.NoError(t, err)
	assert.NotNil(t, result.Alerts)
}
