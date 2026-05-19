package drift

import (
	"testing"
	"time"
)

func makeSuppressResult(key string, kind DriftKind) Result {
	return Result{
		ResourceKey: key,
		Kind:        kind,
	}
}

func TestIsSuppressed_MatchingKey(t *testing.T) {
	sup := NewSuppressor([]Suppression{
		{ResourceKey: "aws_instance.web", Reason: "planned maintenance"},
	})
	r := makeSuppressResult("aws_instance.web", DriftKindChanged)
	if !sup.IsSuppressed(r) {
		t.Error("expected result to be suppressed")
	}
}

func TestIsSuppressed_NoMatch(t *testing.T) {
	sup := NewSuppressor([]Suppression{
		{ResourceKey: "aws_instance.web", Reason: "planned maintenance"},
	})
	r := makeSuppressResult("aws_s3_bucket.data", DriftKindMissing)
	if sup.IsSuppressed(r) {
		t.Error("expected result not to be suppressed")
	}
}

func TestIsSuppressed_ExpiredSuppression(t *testing.T) {
	sup := NewSuppressor([]Suppression{
		{
			ResourceKey: "aws_instance.web",
			Reason:      "old window",
			ExpiresAt:   time.Now().Add(-1 * time.Hour),
		},
	})
	r := makeSuppressResult("aws_instance.web", DriftKindChanged)
	if sup.IsSuppressed(r) {
		t.Error("expected expired suppression to not suppress the result")
	}
}

func TestIsSuppressed_ActiveExpiry(t *testing.T) {
	sup := NewSuppressor([]Suppression{
		{
			ResourceKey: "aws_instance.web",
			Reason:      "active window",
			ExpiresAt:   time.Now().Add(1 * time.Hour),
		},
	})
	r := makeSuppressResult("aws_instance.web", DriftKindChanged)
	if !sup.IsSuppressed(r) {
		t.Error("expected result to be suppressed within active window")
	}
}

func TestApply_FiltersAndCounts(t *testing.T) {
	sup := NewSuppressor([]Suppression{
		{ResourceKey: "aws_instance.web", Reason: "maintenance"},
	})
	results := []Result{
		makeSuppressResult("aws_instance.web", DriftKindChanged),
		makeSuppressResult("aws_s3_bucket.data", DriftKindMissing),
		makeSuppressResult("aws_instance.db", DriftKindExtra),
	}
	kept, count := sup.Apply(results)
	if count != 1 {
		t.Errorf("expected 1 suppressed, got %d", count)
	}
	if len(kept) != 2 {
		t.Errorf("expected 2 kept results, got %d", len(kept))
	}
}

func TestApply_NoSuppressions(t *testing.T) {
	sup := NewSuppressor(nil)
	results := []Result{
		makeSuppressResult("aws_instance.web", DriftKindChanged),
	}
	kept, count := sup.Apply(results)
	if count != 0 {
		t.Errorf("expected 0 suppressed, got %d", count)
	}
	if len(kept) != 1 {
		t.Errorf("expected 1 kept result, got %d", len(kept))
	}
}
