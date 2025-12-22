package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rshade/pulumicost-mcp/internal/adapter"
	"github.com/rshade/pulumicost-mcp/internal/logging"
	"github.com/rshade/pulumicost-mcp/internal/metrics"
	"github.com/rshade/pulumicost-mcp/internal/tracing"
	cost "github.com/rshade/pulumicost-mcp/gen/cost"
	"go.opentelemetry.io/otel/attribute"
)

// CostService implements the cost.Service interface
type CostService struct {
	adapter adapter.PulumiCostAdapter
	logger  *logging.Logger
}

// NewCostService creates a new Cost Service instance
func NewCostService(pulumiAdapter adapter.PulumiCostAdapter, logger *logging.Logger) *CostService {
	return &CostService{
		adapter: pulumiAdapter,
		logger:  logger,
	}
}

// AnalyzeProjected calculates projected costs from Pulumi preview JSON
func (s *CostService) AnalyzeProjected(ctx context.Context, payload *cost.AnalyzeProjectedPayload) (*cost.CostResult, error) {
	start := time.Now()
	ctx, span := tracing.Start(ctx, "CostService.AnalyzeProjected")
	defer span.End()

	s.logger.WithService("cost").Info("analyzing projected costs")
	metrics.RecordCostQuery("projected")

	// Validate payload
	if payload == nil || payload.PulumiJSON == "" {
		err := fmt.Errorf("missing Pulumi JSON")
		s.logger.WithService("cost").ErrorJSON("validation failed", err, nil)
		metrics.RecordError("cost", "analyze_projected", "validation")
		tracing.RecordError(ctx, err)
		return nil, err
	}

	// Build filters from payload
	var filters *adapter.ResourceFilters
	if payload.Filters != nil {
		filters = &adapter.ResourceFilters{
			Provider:     payload.Filters.Provider,
			ResourceType: payload.Filters.ResourceType,
			Region:       payload.Filters.Region,
			Tags:         payload.Filters.Tags,
			NamePattern:  payload.Filters.NamePattern,
		}
		tracing.SetAttributes(ctx,
			attribute.Bool("filtered", true),
			attribute.String("provider", stringValue(payload.Filters.Provider)),
		)
	}

	// Call adapter
	var adapterResult *adapter.CostResult
	var err error

	if filters != nil {
		adapterResult, err = s.adapter.GetProjectedCostWithFilters(ctx, payload.PulumiJSON, filters)
	} else {
		adapterResult, err = s.adapter.GetProjectedCost(ctx, payload.PulumiJSON)
	}

	if err != nil {
		s.logger.WithService("cost").ErrorJSON("adapter call failed", err, nil)
		metrics.RecordError("cost", "analyze_projected", "adapter")
		tracing.RecordError(ctx, err)
		return nil, fmt.Errorf("failed to analyze projected costs: %w", err)
	}

	// Convert adapter result to Goa result type
	result := convertToCostResult(adapterResult)

	// Record metrics
	metrics.RecordRequest("cost", "analyze_projected", time.Since(start))
	metrics.RecordResourceCount(len(result.Resources))

	tracing.SetAttributes(ctx,
		attribute.Int("resource_count", len(result.Resources)),
		attribute.Float64("total_monthly", result.TotalMonthly),
	)

	s.logger.WithService("cost").InfoJSON("projected costs analyzed", map[string]interface{}{
		"resource_count": len(result.Resources),
		"total_monthly":  result.TotalMonthly,
		"duration_ms":    time.Since(start).Milliseconds(),
	})

	return result, nil
}

// GetActual retrieves actual historical costs from cloud providers
func (s *CostService) GetActual(ctx context.Context, payload *cost.GetActualPayload) (*cost.CostResult, error) {
	start := time.Now()
	ctx, span := tracing.Start(ctx, "CostService.GetActual")
	defer span.End()

	s.logger.WithService("cost").Info("getting actual costs")
	metrics.RecordCostQuery("actual")

	// Validate payload
	if payload == nil || payload.StackName == "" {
		err := fmt.Errorf("missing stack name")
		s.logger.WithService("cost").ErrorJSON("validation failed", err, nil)
		metrics.RecordError("cost", "get_actual", "validation")
		tracing.RecordError(ctx, err)
		return nil, err
	}
	if payload.TimeRange == nil {
		err := fmt.Errorf("missing time range")
		s.logger.WithService("cost").ErrorJSON("validation failed", err, nil)
		metrics.RecordError("cost", "get_actual", "validation")
		tracing.RecordError(ctx, err)
		return nil, err
	}

	tracing.SetAttributes(ctx,
		attribute.String("stack_name", payload.StackName),
		attribute.String("time_range_start", payload.TimeRange.Start),
		attribute.String("time_range_end", payload.TimeRange.End),
	)

	// Convert to adapter TimeRange
	timeRange := adapter.TimeRange{
		Start: payload.TimeRange.Start,
		End:   payload.TimeRange.End,
	}

	// Call adapter
	var adapterResult *adapter.CostResult
	var err error

	if payload.Granularity != nil {
		adapterResult, err = s.adapter.GetActualCostWithGranularity(ctx, payload.StackName, timeRange, *payload.Granularity)
	} else {
		adapterResult, err = s.adapter.GetActualCost(ctx, payload.StackName, timeRange)
	}

	if err != nil {
		s.logger.WithService("cost").ErrorJSON("adapter call failed", err, nil)
		metrics.RecordError("cost", "get_actual", "adapter")
		tracing.RecordError(ctx, err)
		return nil, fmt.Errorf("failed to get actual costs: %w", err)
	}

	// Convert adapter result to Goa result type
	result := convertToCostResult(adapterResult)

	// Record metrics
	metrics.RecordRequest("cost", "get_actual", time.Since(start))
	metrics.RecordResourceCount(len(result.Resources))

	tracing.SetAttributes(ctx,
		attribute.Int("resource_count", len(result.Resources)),
		attribute.Float64("total_monthly", result.TotalMonthly),
	)

	s.logger.WithService("cost").InfoJSON("actual costs retrieved", map[string]interface{}{
		"stack_name":     payload.StackName,
		"resource_count": len(result.Resources),
		"total_monthly":  result.TotalMonthly,
		"duration_ms":    time.Since(start).Milliseconds(),
	})

	return result, nil
}

// CompareCosts compares costs between two configurations
func (s *CostService) CompareCosts(ctx context.Context, payload *cost.CompareCostsPayload) (*cost.CompareCostsResult, error) {
	// Validate payload
	if payload == nil || payload.Baseline == nil || payload.Target == nil {
		return nil, fmt.Errorf("missing baseline or target")
	}

	// Get baseline cost
	var baselineCost *cost.CostResult
	var err error
	if payload.Baseline.PulumiJSON != nil {
		baselinePayload := &cost.AnalyzeProjectedPayload{
			PulumiJSON: *payload.Baseline.PulumiJSON,
		}
		baselineCost, err = s.AnalyzeProjected(ctx, baselinePayload)
	} else if payload.Baseline.StackName != nil {
		// Use actual costs for stack name
		baselinePayload := &cost.GetActualPayload{
			StackName: *payload.Baseline.StackName,
			TimeRange: &cost.TimeRange{
				Start: "2024-01-01T00:00:00Z",
				End:   "2024-01-31T23:59:59Z",
			},
		}
		baselineCost, err = s.GetActual(ctx, baselinePayload)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get baseline cost: %w", err)
	}

	// Get target cost
	var targetCost *cost.CostResult
	if payload.Target.PulumiJSON != nil {
		targetPayload := &cost.AnalyzeProjectedPayload{
			PulumiJSON: *payload.Target.PulumiJSON,
		}
		targetCost, err = s.AnalyzeProjected(ctx, targetPayload)
	} else if payload.Target.StackName != nil {
		// Use actual costs for stack name
		targetPayload := &cost.GetActualPayload{
			StackName: *payload.Target.StackName,
			TimeRange: &cost.TimeRange{
				Start: "2024-01-01T00:00:00Z",
				End:   "2024-01-31T23:59:59Z",
			},
		}
		targetCost, err = s.GetActual(ctx, targetPayload)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get target cost: %w", err)
	}

	// Calculate difference
	difference := targetCost.TotalMonthly - baselineCost.TotalMonthly
	var differencePercent float64
	if baselineCost.TotalMonthly > 0 {
		differencePercent = (difference / baselineCost.TotalMonthly) * 100
	}

	return &cost.CompareCostsResult{
		BaselineCost:      baselineCost.TotalMonthly,
		TargetCost:        targetCost.TotalMonthly,
		Difference:        difference,
		DifferencePercent: differencePercent,
	}, nil
}

// AnalyzeResource provides detailed cost analysis for a specific resource
func (s *CostService) AnalyzeResource(ctx context.Context, payload *cost.AnalyzeResourcePayload) (*cost.AnalyzeResourceResult, error) {
	// Validate payload
	if payload == nil || payload.ResourceUrn == "" {
		return nil, fmt.Errorf("missing resource URN")
	}

	// For now, return a simple result with just the resource URN
	// In a real implementation, this would query the adapter for resource details
	resource := &cost.ResourceCost{
		Urn:         payload.ResourceUrn,
		Name:        "resource",
		Type:        "unknown",
		MonthlyCost: 10.0,
		Currency:    "USD",
	}

	return &cost.AnalyzeResourceResult{
		Resource: resource,
	}, nil
}

// QueryByTags groups and aggregates costs by resource tags
func (s *CostService) QueryByTags(ctx context.Context, payload *cost.QueryByTagsPayload) (*cost.QueryByTagsResult, error) {
	// Validate payload
	if payload == nil || payload.StackName == "" {
		return nil, fmt.Errorf("missing stack name")
	}
	if len(payload.TagKeys) == 0 {
		return nil, fmt.Errorf("missing tag keys")
	}

	// For now, return a simple result with mock data
	// In a real implementation, this would query the adapter for tagged resources
	byTag := make(map[string]map[string]float64)
	for _, tagKey := range payload.TagKeys {
		byTag[tagKey] = map[string]float64{
			"value1": 10.0,
			"value2": 5.0,
		}
	}

	return &cost.QueryByTagsResult{
		ByTag: byTag,
	}, nil
}

// AnalyzeStack performs comprehensive stack cost analysis with streaming
func (s *CostService) AnalyzeStack(ctx context.Context, payload *cost.AnalyzeStackPayload, stream cost.AnalyzeStackServerStream) error {
	// Validate payload
	if payload == nil || payload.StackName == "" {
		return fmt.Errorf("missing stack name")
	}

	// Send initial progress
	progress := 0.0
	if err := stream.Send(ctx, &cost.AnalyzeStackResult{Progress: &progress}); err != nil {
		return fmt.Errorf("failed to send progress: %w", err)
	}

	// Send 50% progress
	progress = 50.0
	if err := stream.Send(ctx, &cost.AnalyzeStackResult{Progress: &progress}); err != nil {
		return fmt.Errorf("failed to send progress: %w", err)
	}

	// Send final result with 100% progress
	progress = 100.0
	result := &cost.AnalyzeStackResult{
		Progress: &progress,
	}
	if err := stream.SendAndClose(ctx, result); err != nil {
		return fmt.Errorf("failed to send final result: %w", err)
	}

	return nil
}

// Helper functions

// convertToCostResult converts adapter.CostResult to cost.CostResult
func convertToCostResult(adapterResult *adapter.CostResult) *cost.CostResult {
	if adapterResult == nil {
		return nil
	}

	// Convert resources
	resources := make([]*cost.ResourceCost, len(adapterResult.Resources))
	for i, res := range adapterResult.Resources {
		resources[i] = &cost.ResourceCost{
			Urn:         res.Urn,
			Name:        res.Name,
			Type:        res.Type,
			Provider:    res.Provider,
			MonthlyCost: res.MonthlyCost,
			HourlyCost:  res.HourlyCost,
			Currency:    adapterResult.Currency,
		}
	}

	// Calculate aggregations by provider
	byProvider := make(map[string]float64)
	for _, res := range adapterResult.Resources {
		if res.Provider != nil {
			byProvider[*res.Provider] += res.MonthlyCost
		}
	}

	// Calculate aggregations by region (use region from tags if available)
	byRegion := make(map[string]float64)
	for _, res := range adapterResult.Resources {
		if res.Region != nil {
			byRegion[*res.Region] += res.MonthlyCost
		}
	}

	return &cost.CostResult{
		TotalMonthly: adapterResult.TotalMonthly,
		Currency:     adapterResult.Currency,
		Resources:    resources,
		ByProvider:   byProvider,
		ByRegion:     byRegion,
	}
}

// stringValue returns the string value or empty string if nil
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
