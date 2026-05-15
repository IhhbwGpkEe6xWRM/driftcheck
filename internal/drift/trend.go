package drift

import (
	"sort"
	"time"
)

// TrendPoint represents a snapshot of drift metrics at a point in time.
type TrendPoint struct {
	Timestamp    time.Time      `json:"timestamp"`
	TotalDrifts  int            `json:"total_drifts"`
	ByKind       map[string]int `json:"by_kind"`
	ByType       map[string]int `json:"by_type"`
	AverageScore float64        `json:"average_score"`
}

// Trend holds a series of TrendPoints ordered by time.
type Trend struct {
	Points []TrendPoint `json:"points"`
}

// AddSnapshot appends a new trend point derived from the given drift results.
func (t *Trend) AddSnapshot(results []Result, scores []ScoredDrift, ts time.Time) {
	byKind := make(map[string]int)
	byType := make(map[string]int)

	for _, r := range results {
		byKind[string(r.Kind)]++
		rType, _, _ := splitKey(r.ResourceKey)
		byType[rType]++
	}

	var totalScore float64
	for _, s := range scores {
		totalScore += s.Score
	}
	var avg float64
	if len(scores) > 0 {
		avg = totalScore / float64(len(scores))
	}

	t.Points = append(t.Points, TrendPoint{
		Timestamp:    ts,
		TotalDrifts:  len(results),
		ByKind:       byKind,
		ByType:       byType,
		AverageScore: avg,
	})

	sort.Slice(t.Points, func(i, j int) bool {
		return t.Points[i].Timestamp.Before(t.Points[j].Timestamp)
	})
}

// Delta returns the change in total drifts between the last two snapshots.
// Returns 0 if fewer than two points exist.
func (t *Trend) Delta() int {
	if len(t.Points) < 2 {
		return 0
	}
	last := t.Points[len(t.Points)-1]
	prev := t.Points[len(t.Points)-2]
	return last.TotalDrifts - prev.TotalDrifts
}

// Latest returns the most recent TrendPoint, or nil if empty.
func (t *Trend) Latest() *TrendPoint {
	if len(t.Points) == 0 {
		return nil
	}
	return &t.Points[len(t.Points)-1]
}
