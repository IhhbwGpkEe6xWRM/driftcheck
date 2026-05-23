package drift

import (
	"strings"
	"testing"
	"time"
)

func makeSnapResult(resourceKey, kind string, attrs []string) DriftResult {
	r := makeTrendResult(resourceKey, kind)
	if len(attrs) > 0 {
		r.ChangedAttributes = attrs
	}
	return r
}

func makeSnap(t time.Time, results []DriftResult) TrendSnapshot {
	return TrendSnapshot{Timestamp: t, Results: results}
}

var (
	t1 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 = time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
)

func TestDiffSnapshots_AllNew(t *testing.T) {
	older := makeSnap(t1, nil)
	newer := makeSnap(t2, []DriftResult{
		makeSnapResult("aws_instance.web", "changed", []string{"ami"}),
	})

	diff := DiffSnapshots(older, newer)

	if len(diff.AddedDrifts) != 1 {
		t.Errorf("expected 1 added, got %d", len(diff.AddedDrifts))
	}
	if len(diff.RemovedDrifts) != 0 {
		t.Errorf("expected 0 removed, got %d", len(diff.RemovedDrifts))
	}
	if diff.DeltaTotal != 1 {
		t.Errorf("expected delta +1, got %d", diff.DeltaTotal)
	}
}

func TestDiffSnapshots_AllResolved(t *testing.T) {
	older := makeSnap(t1, []DriftResult{
		makeSnapResult("aws_s3_bucket.logs", "missing", nil),
	})
	newer := makeSnap(t2, nil)

	diff := DiffSnapshots(older, newer)

	if len(diff.RemovedDrifts) != 1 {
		t.Errorf("expected 1 removed, got %d", len(diff.RemovedDrifts))
	}
	if len(diff.AddedDrifts) != 0 {
		t.Errorf("expected 0 added, got %d", len(diff.AddedDrifts))
	}
	if diff.DeltaTotal != -1 {
		t.Errorf("expected delta -1, got %d", diff.DeltaTotal)
	}
}

func TestDiffSnapshots_Persisted(t *testing.T) {
	res := makeSnapResult("aws_instance.web", "changed", []string{"instance_type"})
	older := makeSnap(t1, []DriftResult{res})
	newer := makeSnap(t2, []DriftResult{res})

	diff := DiffSnapshots(older, newer)

	if len(diff.Persisted) != 1 {
		t.Errorf("expected 1 persisted, got %d", len(diff.Persisted))
	}
	if len(diff.AddedDrifts) != 0 || len(diff.RemovedDrifts) != 0 {
		t.Error("expected no added or removed drifts")
	}
	if diff.DeltaTotal != 0 {
		t.Errorf("expected delta 0, got %d", diff.DeltaTotal)
	}
}

func TestDiffSnapshots_Summary_ContainsKeyFields(t *testing.T) {
	older := makeSnap(t1, nil)
	newer := makeSnap(t2, []DriftResult{
		makeSnapResult("aws_instance.web", "changed", []string{"ami"}),
	})

	diff := DiffSnapshots(older, newer)
	summary := diff.Summary()

	for _, want := range []string{"+1 added", "-0 resolved", "delta"} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary missing %q: %s", want, summary)
		}
	}
}
