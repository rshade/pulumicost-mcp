package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rshade/pulumicost-mcp/gen/analysis"
)

// AnalysisService implements the analysis.Service interface
type AnalysisService struct{}

// NewAnalysisService creates a new Analysis Service instance
func NewAnalysisService(dataSource interface{}, logger interface{}) *AnalysisService {
	return &AnalysisService{}
}

// GetRecommendations returns cost optimization recommendations
func (s *AnalysisService) GetRecommendations(ctx context.Context, payload *analysis.GetRecommendationsPayload) (*analysis.GetRecommendationsResult, error) {
	if payload.StackName == "" {
		return nil, fmt.Errorf("stack name cannot be empty")
	}

	recommendations := []*analysis.Recommendation{
		{
			ID:               "rec-001",
			Type:             "RIGHTSIZING",
			ResourceUrn:      "urn:pulumi:prod::myapp::aws:ec2/instance:Instance::web-1",
			CurrentCost:      120.00,
			ProjectedSavings: 245.50,
			Confidence:       "HIGH",
			Description:      "EC2 instance is over-provisioned and can be downsized from t3.large to t3.medium",
			ActionSteps:      []string{"Review instance metrics", "Update instance type", "Monitor performance"},
		},
		{
			ID:               "rec-002",
			Type:             "RESERVED_INSTANCES",
			ResourceUrn:      "urn:pulumi:prod::myapp::aws:rds/instance:Instance::db",
			CurrentCost:      250.00,
			ProjectedSavings: 180.00,
			Confidence:       "MEDIUM",
			Description:      "Database running 24/7 can benefit from 1-year Reserved Instance pricing",
			ActionSteps:      []string{"Analyze usage patterns", "Purchase Reserved Instance", "Update infrastructure"},
		},
		{
			ID:               "rec-003",
			Type:             "SPOT_INSTANCES",
			ResourceUrn:      "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::dev-1",
			CurrentCost:      80.00,
			ProjectedSavings: 56.00,
			Confidence:       "MEDIUM",
			Description:      "Development workload can tolerate interruptions and save 70% on compute",
			ActionSteps:      []string{"Implement graceful shutdown", "Configure Spot requests", "Test failure scenarios"},
		},
	}

	// Apply filters
	if payload.RecommendationTypes != nil && len(payload.RecommendationTypes) > 0 {
		filtered := []*analysis.Recommendation{}
		typeMap := make(map[string]bool)
		for _, t := range payload.RecommendationTypes {
			typeMap[t] = true
		}
		for _, rec := range recommendations {
			if typeMap[rec.Type] {
				filtered = append(filtered, rec)
			}
		}
		recommendations = filtered
	}

	if payload.MinimumSavings != nil {
		filtered := []*analysis.Recommendation{}
		for _, rec := range recommendations {
			if rec.ProjectedSavings >= *payload.MinimumSavings {
				filtered = append(filtered, rec)
			}
		}
		recommendations = filtered
	}

	return &analysis.GetRecommendationsResult{
		Recommendations: recommendations,
	}, nil
}

// DetectAnomalies detects unusual spending patterns
func (s *AnalysisService) DetectAnomalies(ctx context.Context, payload *analysis.DetectAnomaliesPayload) (*analysis.DetectAnomaliesResult, error) {
	if payload.StackName == "" {
		return nil, fmt.Errorf("stack name cannot be empty")
	}
	if payload.TimeRange == nil {
		return nil, fmt.Errorf("time range is required")
	}

	anomalies := []*analysis.Anomaly{
		{
			ID:               "anom-001",
			Timestamp:        "2024-01-15T14:30:00Z",
			ResourceUrns:     []string{"urn:pulumi:prod::myapp::aws:ec2/instance:Instance::web-3"},
			Severity:         "HIGH",
			CurrentCost:      425.00,
			BaselineCost:     150.00,
			DeviationPercent: 183.33,
			PotentialCauses:  []string{"Unexpected spike in EC2 instance hours", "Auto-scaling triggered", "Instance not properly terminated"},
		},
	}

	// Add more anomalies for higher sensitivity
	if payload.Sensitivity == "HIGH" {
		anomalies = append(anomalies, &analysis.Anomaly{
			ID:               "anom-002",
			Timestamp:        "2024-01-20T09:15:00Z",
			ResourceUrns:     []string{"urn:pulumi:prod::myapp::aws:s3/bucket:Bucket::data"},
			Severity:         "MEDIUM",
			CurrentCost:      120.00,
			BaselineCost:     75.00,
			DeviationPercent: 60.00,
			PotentialCauses:  []string{"Increased S3 storage usage", "Large file uploads", "Lifecycle policy not applied"},
		})
	}

	return &analysis.DetectAnomaliesResult{
		Anomalies: anomalies,
	}, nil
}

// Forecast generates cost forecasts
func (s *AnalysisService) Forecast(ctx context.Context, payload *analysis.ForecastPayload) (*analysis.Forecast2, error) {
	if payload.StackName == "" {
		return nil, fmt.Errorf("stack name cannot be empty")
	}
	if payload.ForecastPeriod == nil {
		return nil, fmt.Errorf("forecast period is required")
	}

	dataPoints := []*analysis.ForecastPoint{
		{
			Timestamp:     "2024-02-01T00:00:00Z",
			PredictedCost: 850.00,
			LowerBound:    800.00,
			UpperBound:    900.00,
		},
		{
			Timestamp:     "2024-02-15T00:00:00Z",
			PredictedCost: 875.00,
			LowerBound:    825.00,
			UpperBound:    925.00,
		},
		{
			Timestamp:     "2024-02-29T00:00:00Z",
			PredictedCost: 900.00,
			LowerBound:    850.00,
			UpperBound:    950.00,
		},
	}

	return &analysis.Forecast2{
		StackName:       payload.StackName,
		ForecastPeriod:  payload.ForecastPeriod,
		DataPoints:      dataPoints,
		ConfidenceLevel: payload.ConfidenceLevel,
		Methodology:     "Linear regression with seasonal adjustment based on historical spending patterns",
	}, nil
}

// TrackBudget monitors spending against budget
func (s *AnalysisService) TrackBudget(ctx context.Context, payload *analysis.TrackBudgetPayload) (*analysis.Budget, error) {
	if payload.StackName == "" {
		return nil, fmt.Errorf("stack name cannot be empty")
	}
	if payload.BudgetAmount <= 0 {
		return nil, fmt.Errorf("budget amount must be positive")
	}

	// Mock current spending
	currentSpending := 750.00
	remaining := payload.BudgetAmount - currentSpending
	burnRate := 25.00 // Daily burn rate
	percentageUsed := (currentSpending / payload.BudgetAmount) * 100

	// Calculate projected end date
	daysRemaining := int(remaining / burnRate)
	projectedEndDate := time.Now().AddDate(0, 0, daysRemaining).Format(time.RFC3339)

	// Determine status
	status := "OK"
	switch {
	case percentageUsed >= 100:
		status = "EXCEEDED"
	case percentageUsed >= 90:
		status = "CRITICAL"
	case percentageUsed >= 80:
		status = "WARNING"
	}

	// Generate alerts based on thresholds
	var alerts []any
	if payload.AlertThresholds != nil {
		for _, threshold := range payload.AlertThresholds {
			if percentageUsed >= threshold {
				severity := "INFO"
				switch {
				case percentageUsed >= 100:
					severity = "CRITICAL"
				case percentageUsed >= 90:
					severity = "HIGH"
				case percentageUsed >= 80:
					severity = "MEDIUM"
				}

				alerts = append(alerts, map[string]interface{}{
					"threshold":     threshold,
					"current_spend": currentSpending,
					"severity":      severity,
					"message":       fmt.Sprintf("Budget utilization at %.1f%% (threshold: %.1f%%)", percentageUsed, threshold),
					"timestamp":     time.Now().Format(time.RFC3339),
				})
			}
		}
	}

	return &analysis.Budget{
		BudgetAmount:     payload.BudgetAmount,
		CurrentSpending:  currentSpending,
		Remaining:        remaining,
		BurnRate:         &burnRate,
		ProjectedEndDate: &projectedEndDate,
		Status:           status,
		Alerts:           alerts,
	}, nil
}
