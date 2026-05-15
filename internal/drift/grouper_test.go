package drift

import (
	"testing"
)

func TestGroupByType_Empty(t *testing.T) {
	g := GroupByType(nil)
	if len(g.ByType) != 0 {
		t.Errorf("expected empty map, got %v", g.ByType)
	}
	if len(g.Ordered) != 0 {
		t.Errorf("expected empty ordered slice, got %v", g.Ordered)
	}
}

func TestGroupByType_SingleType(t *testing.T) {
	results := []DriftResult{
		makeDriftResult("aws_instance.web", KindChanged),
		makeDriftResult("aws_instance.api", KindMissing),
	}
	g := GroupByType(results)

	if len(g.ByType) != 1 {
		t.Fatalf("expected 1 type, got %d", len(g.ByType))
	}
	if items, ok := g.ByType["aws_instance"]; !ok || len(items) != 2 {
		t.Errorf("expected 2 aws_instance entries, got %v", items)
	}
}

func TestGroupByType_MultipleTypes_OrderedAlphabetically(t *testing.T) {
	results := []DriftResult{
		makeDriftResult("aws_s3_bucket.data", KindExtra),
		makeDriftResult("aws_instance.web", KindChanged),
		makeDriftResult("aws_lambda_function.handler", KindMissing),
	}
	g := GroupByType(results)

	if len(g.Ordered) != 3 {
		t.Fatalf("expected 3 types, got %d", len(g.Ordered))
	}
	expected := []string{"aws_instance", "aws_lambda_function", "aws_s3_bucket"}
	for i, name := range expected {
		if g.Ordered[i] != name {
			t.Errorf("Ordered[%d]: want %q, got %q", i, name, g.Ordered[i])
		}
	}
}

func TestGroupByKind_Mixed(t *testing.T) {
	results := []DriftResult{
		makeDriftResult("aws_instance.a", KindChanged),
		makeDriftResult("aws_instance.b", KindMissing),
		makeDriftResult("aws_instance.c", KindChanged),
		makeDriftResult("aws_s3_bucket.x", KindExtra),
	}
	byKind := GroupByKind(results)

	if len(byKind[KindChanged]) != 2 {
		t.Errorf("expected 2 changed, got %d", len(byKind[KindChanged]))
	}
	if len(byKind[KindMissing]) != 1 {
		t.Errorf("expected 1 missing, got %d", len(byKind[KindMissing]))
	}
	if len(byKind[KindExtra]) != 1 {
		t.Errorf("expected 1 extra, got %d", len(byKind[KindExtra]))
	}
}

func TestTopDriftedTypes_ReturnsTopN(t *testing.T) {
	results := []DriftResult{
		makeDriftResult("aws_instance.a", KindChanged),
		makeDriftResult("aws_instance.b", KindChanged),
		makeDriftResult("aws_instance.c", KindChanged),
		makeDriftResult("aws_s3_bucket.x", KindMissing),
		makeDriftResult("aws_s3_bucket.y", KindMissing),
		makeDriftResult("aws_lambda_function.fn", KindExtra),
	}
	g := GroupByType(results)
	top := TopDriftedTypes(g, 2)

	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
	if top[0] != "aws_instance" {
		t.Errorf("expected aws_instance first, got %q", top[0])
	}
	if top[1] != "aws_s3_bucket" {
		t.Errorf("expected aws_s3_bucket second, got %q", top[1])
	}
}

func TestTopDriftedTypes_NLargerThanAvailable(t *testing.T) {
	results := []DriftResult{
		makeDriftResult("aws_instance.a", KindChanged),
	}
	g := GroupByType(results)
	top := TopDriftedTypes(g, 10)
	if len(top) != 1 {
		t.Errorf("expected 1 result, got %d", len(top))
	}
}
