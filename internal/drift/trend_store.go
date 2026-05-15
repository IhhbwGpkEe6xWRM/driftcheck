package drift

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const defaultMaxPoints = 90

// SaveTrend persists a Trend to the given file path as JSON.
func SaveTrend(path string, trend *Trend) error {
	data, err := json.MarshalIndent(trend, "", "  ")
	if err != nil {
		return fmt.Errorf("trend: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("trend: write %s: %w", path, err)
	}
	return nil
}

// LoadTrend reads a Trend from the given file path.
// Returns an empty Trend if the file does not exist.
func LoadTrend(path string) (*Trend, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Trend{}, nil
		}
		return nil, fmt.Errorf("trend: read %s: %w", path, err)
	}
	var trend Trend
	if err := json.Unmarshal(data, &trend); err != nil {
		return nil, fmt.Errorf("trend: unmarshal: %w", err)
	}
	return &trend, nil
}

// AppendSnapshot loads the trend from path, adds a new snapshot, trims old
// points beyond maxPoints, and saves it back. If maxPoints <= 0 the default
// of 90 is used.
func AppendSnapshot(path string, results []Result, scores []ScoredDrift, maxPoints int) error {
	if maxPoints <= 0 {
		maxPoints = defaultMaxPoints
	}
	trend, err := LoadTrend(path)
	if err != nil {
		return err
	}
	trend.AddSnapshot(results, scores, time.Now().UTC())
	if len(trend.Points) > maxPoints {
		trend.Points = trend.Points[len(trend.Points)-maxPoints:]
	}
	return SaveTrend(path, trend)
}
