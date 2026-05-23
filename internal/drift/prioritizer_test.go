package drift

import (
	"testing"
)

func makePrioritizerResult(key string, kind DriftKind, severity string) AnnotatedResult {
	return AnnotatedResult{
		Result: DriftResult{
			ResourceKey: key,
			Kind:        kind,
		},
		Severity: severity,
	}
}

func TestPrioritize_MissingResourceIsHigh(t *testing.T) {
	input := []AnnotatedResult{
		makePrioritizerResult("aws_instance.web", DriftMissing, SeverityLow),
	}
	out := Prioritize(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Priority != PriorityHigh {
		t.Errorf("expected high priority for missing resource, got %s", out[0].Priority)
	}
}

func TestPrioritize_CriticalSeverityIsHigh(t *testing.T) {
	input := []AnnotatedResult{
		makePrioritizerResult("aws_s3_bucket.data", DriftChanged, SeverityCritical),
	}
	out := Prioritize(input)
	if out[0].Priority != PriorityHigh {
		t.Errorf("expected high priority for critical severity, got %s", out[0].Priority)
	}
}

func TestPrioritize_MediumSeverityIsMedium(t *testing.T) {
	input := []AnnotatedResult{
		makePrioritizerResult("aws_sg.allow_http", DriftChanged, SeverityMedium),
	}
	out := Prioritize(input)
	if out[0].Priority != PriorityMedium {
		t.Errorf("expected medium priority, got %s", out[0].Priority)
	}
}

func TestPrioritize_ExtraResourceIsLow(t *testing.T) {
	input := []AnnotatedResult{
		makePrioritizerResult("aws_instance.old", DriftExtra, SeverityLow),
	}
	out := Prioritize(input)
	if out[0].Priority != PriorityLow {
		t.Errorf("expected low priority for extra resource, got %s", out[0].Priority)
	}
}

func TestPrioritize_SortedHighestFirst(t *testing.T) {
	input := []AnnotatedResult{
		makePrioritizerResult("aws_instance.a", DriftExtra, SeverityLow),
		makePrioritizerResult("aws_instance.b", DriftMissing, SeverityLow),
		makePrioritizerResult("aws_sg.c", DriftChanged, SeverityMedium),
	}
	out := Prioritize(input)
	if out[0].Priority != PriorityHigh {
		t.Errorf("first result should be high, got %s", out[0].Priority)
	}
	if out[1].Priority != PriorityMedium {
		t.Errorf("second result should be medium, got %s", out[1].Priority)
	}
	if out[2].Priority != PriorityLow {
		t.Errorf("third result should be low, got %s", out[2].Priority)
	}
}

func TestFilterByPriority_ReturnsMatchingOnly(t *testing.T) {
	input := []AnnotatedResult{
		makePrioritizerResult("aws_instance.a", DriftMissing, SeverityLow),
		makePrioritizerResult("aws_instance.b", DriftExtra, SeverityLow),
		makePrioritizerResult("aws_sg.c", DriftChanged, SeverityMedium),
	}
	all := Prioritize(input)
	high := FilterByPriority(all, PriorityHigh)
	if len(high) != 1 {
		t.Errorf("expected 1 high-priority result, got %d", len(high))
	}
	medium := FilterByPriority(all, PriorityMedium)
	if len(medium) != 1 {
		t.Errorf("expected 1 medium-priority result, got %d", len(medium))
	}
}
