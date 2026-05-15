package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeTrendResult(kind DriftKind, resKey string) Result {
	return Result{Kind: kind, ResourceKey: resKey}
}

func TestAddSnapshot_PopulatesFields(t *testing.T) {
	var trend Trend
	results := []Result{
		makeTrendResult(DriftKindMissing, "aws_instance.web"),
		makeTrendResult(DriftKindChanged, "aws_s3_bucket.logs"),
	}
	scores := []ScoredDrift{
		{Result: results[0], Score: 10},
		{Result: results[1], Score: 4},
	}
	ts := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	trend.AddSnapshot(results, scores, ts)

	if len(trend.Points) != 1 {
		t.Fatalf("expected 1 point, got %d", len(trend.Points))
	}
	p := trend.Points[0]
	if p.TotalDrifts != 2 {
		t.Errorf("expected TotalDrifts=2, got %d", p.TotalDrifts)
	}
	if p.AverageScore != 7.0 {
		t.Errorf("expected AverageScore=7.0, got %f", p.AverageScore)
	}
	if p.ByKind[string(DriftKindMissing)] != 1 {
		t.Errorf("expected ByKind[missing]=1")
	}
}

func TestDelta_LessThanTwoPoints(t *testing.T) {
	var trend Trend
	if trend.Delta() != 0 {
		t.Error("expected delta 0 for empty trend")
	}
}

func TestDelta_TwoPoints(t *testing.T) {
	var trend Trend
	ts := time.Now().UTC()
	trend.AddSnapshot([]Result{makeTrendResult(DriftKindMissing, "aws_instance.a")}, nil, ts)
	trend.AddSnapshot([]Result{
		makeTrendResult(DriftKindMissing, "aws_instance.a"),
		makeTrendResult(DriftKindChanged, "aws_s3_bucket.b"),
	}, nil, ts.Add(time.Hour))
	if trend.Delta() != 1 {
		t.Errorf("expected delta 1, got %d", trend.Delta())
	}
}

func TestSaveAndLoadTrend(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "trend.json")

	var trend Trend
	trend.AddSnapshot([]Result{makeTrendResult(DriftKindExtra, "aws_instance.x")}, nil, time.Now().UTC())

	if err := SaveTrend(path, &trend); err != nil {
		t.Fatalf("SaveTrend: %v", err)
	}
	loaded, err := LoadTrend(path)
	if err != nil {
		t.Fatalf("LoadTrend: %v", err)
	}
	if len(loaded.Points) != 1 {
		t.Errorf("expected 1 point after reload, got %d", len(loaded.Points))
	}
}

func TestLoadTrend_MissingFile(t *testing.T) {
	trend, err := LoadTrend("/nonexistent/trend.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if trend == nil || len(trend.Points) != 0 {
		t.Error("expected empty trend for missing file")
	}
}

func TestAppendSnapshot_TrimsOldPoints(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "trend.json")

	for i := 0; i < 5; i++ {
		if err := AppendSnapshot(path, []Result{makeTrendResult(DriftKindMissing, "aws_instance.a")}, nil, 3); err != nil {
			t.Fatalf("AppendSnapshot iteration %d: %v", i, err)
		}
	}
	data, _ := os.ReadFile(path)
	var trend Trend
	_ = json.Unmarshal(data, &trend)
	if len(trend.Points) != 3 {
		t.Errorf("expected 3 points after trimming, got %d", len(trend.Points))
	}
}
