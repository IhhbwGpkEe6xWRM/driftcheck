package drift

import (
	"bytes"
	"strings"
	"testing"
)

func makeDriftResult(kind DriftKind, resourceKey string) DriftResult {
	return DriftResult{
		Kind:        kind,
		ResourceKey: resourceKey,
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize(nil)
	if s.Total != 0 || s.Missing != 0 || s.Changed != 0 || s.Extra != 0 {
		t.Errorf("expected zero summary, got %+v", s)
	}
}

func TestSummarize_MixedKinds(t *testing.T) {
	results := []DriftResult{
		makeDriftResult(KindMissing, "aws_instance.web"),
		makeDriftResult(KindChanged, "aws_instance.api"),
		makeDriftResult(KindChanged, "aws_s3_bucket.data"),
		makeDriftResult(KindExtra, "aws_s3_bucket.logs"),
	}
	s := Summarize(results)

	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Missing != 1 {
		t.Errorf("expected Missing=1, got %d", s.Missing)
	}
	if s.Changed != 2 {
		t.Errorf("expected Changed=2, got %d", s.Changed)
	}
	if s.Extra != 1 {
		t.Errorf("expected Extra=1, got %d", s.Extra)
	}
}

func TestSummarize_ByType(t *testing.T) {
	results := []DriftResult{
		makeDriftResult(KindMissing, "aws_instance.web"),
		makeDriftResult(KindChanged, "aws_instance.api"),
		makeDriftResult(KindExtra, "aws_s3_bucket.logs"),
	}
	s := Summarize(results)

	if s.ByType["aws_instance"] != 2 {
		t.Errorf("expected aws_instance=2, got %d", s.ByType["aws_instance"])
	}
	if s.ByType["aws_s3_bucket"] != 1 {
		t.Errorf("expected aws_s3_bucket=1, got %d", s.ByType["aws_s3_bucket"])
	}
}

func TestPrintSummary_ContainsKeyFields(t *testing.T) {
	results := []DriftResult{
		makeDriftResult(KindMissing, "aws_instance.web"),
		makeDriftResult(KindChanged, "aws_s3_bucket.data"),
	}
	s := Summarize(results)

	var buf bytes.Buffer
	PrintSummary(&buf, s)
	out := buf.String()

	for _, want := range []string{"Total", "Missing", "Changed", "Extra", "aws_instance", "aws_s3_bucket"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}
