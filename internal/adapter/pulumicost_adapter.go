package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// PulumiCostAdapter provides integration with pulumicost-core binary
type PulumiCostAdapter interface {
	GetProjectedCost(ctx context.Context, pulumiJSON string) (*CostResult, error)
	GetProjectedCostWithFilters(ctx context.Context, pulumiJSON string, filters *ResourceFilters) (*CostResult, error)
	GetActualCost(ctx context.Context, stackName string, timeRange TimeRange) (*CostResult, error)
	GetActualCostWithGranularity(ctx context.Context, stackName string, timeRange TimeRange, granularity string) (*CostResult, error)
	GetCorePath() string
}

// pulumiCostAdapter is the concrete implementation
type pulumiCostAdapter struct {
	corePath string
}

// NewPulumiCostAdapter creates a new PulumiCost adapter instance
func NewPulumiCostAdapter(corePath string) PulumiCostAdapter {
	return &pulumiCostAdapter{
		corePath: corePath,
	}
}

// GetCorePath returns the path to the pulumicost-core binary
func (a *pulumiCostAdapter) GetCorePath() string {
	return a.corePath
}

// GetProjectedCost calculates projected costs from Pulumi preview JSON
func (a *pulumiCostAdapter) GetProjectedCost(ctx context.Context, pulumiJSON string) (*CostResult, error) {
	return a.GetProjectedCostWithFilters(ctx, pulumiJSON, nil)
}

// GetProjectedCostWithFilters calculates projected costs with resource filters
func (a *pulumiCostAdapter) GetProjectedCostWithFilters(ctx context.Context, pulumiJSON string, filters *ResourceFilters) (*CostResult, error) {
	// Validate JSON first
	var previewData map[string]interface{}
	if err := json.Unmarshal([]byte(pulumiJSON), &previewData); err != nil {
		return nil, fmt.Errorf("invalid Pulumi JSON: %w", err)
	}

	// Prepare command with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, a.corePath, "analyze", "--projected")

	// Pass Pulumi JSON via stdin
	cmd.Stdin = strings.NewReader(pulumiJSON)

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	if err := cmd.Run(); err != nil {
		if cmdCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("pulumicost timeout: %w", cmdCtx.Err())
		}
		if cmdCtx.Err() == context.Canceled {
			return nil, fmt.Errorf("context canceled: %w", cmdCtx.Err())
		}
		return nil, fmt.Errorf("pulumicost execution failed: %w (stderr: %s)", err, stderr.String())
	}

	// Parse output
	var result CostResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse pulumicost output: %w", err)
	}

	// Apply filters if provided
	if filters != nil {
		result.Resources = applyFilters(result.Resources, filters)
		// Recalculate total
		result.TotalMonthly = calculateTotal(result.Resources)
	}

	return &result, nil
}

// GetActualCost retrieves historical costs from cloud providers
func (a *pulumiCostAdapter) GetActualCost(ctx context.Context, stackName string, timeRange TimeRange) (*CostResult, error) {
	return a.GetActualCostWithGranularity(ctx, stackName, timeRange, "")
}

// GetActualCostWithGranularity retrieves historical costs with specific time granularity
func (a *pulumiCostAdapter) GetActualCostWithGranularity(ctx context.Context, stackName string, timeRange TimeRange, granularity string) (*CostResult, error) {
	// Validate time range
	if _, err := time.Parse(time.RFC3339, timeRange.Start); err != nil {
		return nil, fmt.Errorf("invalid start time format: %w", err)
	}
	if _, err := time.Parse(time.RFC3339, timeRange.End); err != nil {
		return nil, fmt.Errorf("invalid end time format: %w", err)
	}

	// Prepare command with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	args := []string{"analyze", "--actual", "--stack", stackName, "--start", timeRange.Start, "--end", timeRange.End}
	if granularity != "" {
		args = append(args, "--granularity", granularity)
	}

	cmd := exec.CommandContext(cmdCtx, a.corePath, args...)

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	if err := cmd.Run(); err != nil {
		if cmdCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("pulumicost timeout: %w", cmdCtx.Err())
		}
		if cmdCtx.Err() == context.Canceled {
			return nil, fmt.Errorf("context canceled: %w", cmdCtx.Err())
		}
		return nil, fmt.Errorf("pulumicost execution failed: %w (stderr: %s)", err, stderr.String())
	}

	// Parse output
	var result CostResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse pulumicost output: %w", err)
	}

	return &result, nil
}

// CostResult represents the result of a cost analysis
type CostResult struct {
	TotalMonthly float64         `json:"total_monthly"`
	TotalHourly  *float64        `json:"total_hourly,omitempty"`
	Currency     string          `json:"currency"`
	Resources    []ResourceCost  `json:"resources"`
	Breakdown    *CostBreakdown  `json:"breakdown,omitempty"`
}

// ResourceCost represents the cost of a single resource
type ResourceCost struct {
	Urn         string            `json:"urn"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Provider    *string           `json:"provider,omitempty"`
	MonthlyCost float64           `json:"monthly_cost"`
	HourlyCost  *float64          `json:"hourly_cost,omitempty"`
	Region      *string           `json:"region,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// ResourceFilters specifies criteria for filtering resources
type ResourceFilters struct {
	Provider     *string
	ResourceType *string
	Region       *string
	Tags         map[string]string
	NamePattern  *string
}

// TimeRange represents a time period for cost queries
type TimeRange struct {
	Start string
	End   string
}

// CostBreakdown provides detailed cost breakdown by time period
type CostBreakdown struct {
	Daily   []DailyCost   `json:"daily,omitempty"`
	Weekly  []WeeklyCost  `json:"weekly,omitempty"`
	Monthly []MonthlyCost `json:"monthly,omitempty"`
}

// DailyCost represents daily cost data
type DailyCost struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

// WeeklyCost represents weekly cost data
type WeeklyCost struct {
	Week   string  `json:"week"`
	Amount float64 `json:"amount"`
}

// MonthlyCost represents monthly cost data
type MonthlyCost struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
}

// applyFilters filters resources based on the provided criteria
func applyFilters(resources []ResourceCost, filters *ResourceFilters) []ResourceCost {
	if filters == nil {
		return resources
	}

	var filtered []ResourceCost
	for _, resource := range resources {
		if !matchesFilter(resource, filters) {
			continue
		}
		filtered = append(filtered, resource)
	}
	return filtered
}

// matchesFilter checks if a resource matches the filter criteria
// nolint:gocognit // linear filter checks - refactoring would reduce readability
func matchesFilter(resource ResourceCost, filters *ResourceFilters) bool {
	// Check provider filter
	if filters.Provider != nil {
		if resource.Provider == nil || *resource.Provider != *filters.Provider {
			return false
		}
	}

	// Check resource type filter
	if filters.ResourceType != nil {
		if resource.Type != *filters.ResourceType {
			return false
		}
	}

	// Check region filter
	if filters.Region != nil {
		if resource.Region == nil || *resource.Region != *filters.Region {
			return false
		}
	}

	// Check tag filters
	if len(filters.Tags) > 0 {
		for key, value := range filters.Tags {
			resourceValue, ok := resource.Tags[key]
			if !ok || resourceValue != value {
				return false
			}
		}
	}

	// Check name pattern filter
	if filters.NamePattern != nil {
		if !strings.Contains(resource.Name, *filters.NamePattern) {
			return false
		}
	}

	return true
}

// calculateTotal sums up the monthly costs of all resources
func calculateTotal(resources []ResourceCost) float64 {
	total := 0.0
	for _, resource := range resources {
		total += resource.MonthlyCost
	}
	return total
}
