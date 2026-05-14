package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved snapshot of drift results used for comparison
// across runs to detect new or resolved drift.
type Baseline struct {
	CreatedAt time.Time       `json:"created_at"`
	Results   []DriftResult   `json:"results"`
	Meta      BaselineMeta    `json:"meta"`
}

type BaselineMeta struct {
	StateFile string `json:"state_file"`
	Region    string `json:"region"`
}

// SaveBaseline writes drift results to a JSON file at the given path.
func SaveBaseline(path string, results []DriftResult, meta BaselineMeta) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Results:   results,
		Meta:      meta,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal baseline: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write baseline file: %w", err)
	}
	return nil
}

// LoadBaseline reads a previously saved baseline from disk.
func LoadBaseline(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read baseline file: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("parse baseline file: %w", err)
	}
	return &b, nil
}

// DiffBaseline compares current results against a baseline and returns
// newly introduced drift and drift that has been resolved since the baseline.
func DiffBaseline(baseline *Baseline, current []DriftResult) (newDrift []DriftResult, resolved []DriftResult) {
	baselineKeys := make(map[string]struct{}, len(baseline.Results))
	for _, r := range baseline.Results {
		baselineKeys[driftResultKey(r)] = struct{}{}
	}
	currentKeys := make(map[string]struct{}, len(current))
	for _, r := range current {
		currentKeys[driftResultKey(r)] = struct{}{}
		if _, found := baselineKeys[driftResultKey(r)]; !found {
			newDrift = append(newDrift, r)
		}
	}
	for _, r := range baseline.Results {
		if _, found := currentKeys[driftResultKey(r)]; !found {
			resolved = append(resolved, r)
		}
	}
	return newDrift, resolved
}

func driftResultKey(r DriftResult) string {
	return fmt.Sprintf("%s|%s|%s", r.ResourceKey, r.Attribute, string(r.Kind))
}
