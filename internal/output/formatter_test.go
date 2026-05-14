package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/driftcheck/internal/drift"
	"github.com/yourorg/driftcheck/internal/output"
)

func makeDrift(key, kind, attr string, expected, actual interface{}) drift.DriftResult {
	return drift.DriftResult{
		ResourceKey: key,
		Kind:        drift.DriftKind(kind),
		Attribute:   attr,
		Expected:    expected,
		Actual:      actual,
	}
}

func TestWriteText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(output.FormatText, &buf)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteText_WithDrift(t *testing.T) {
	results := []drift.DriftResult{
		makeDrift("aws_instance.web", "changed", "instance_type", "t2.micro", "t3.small"),
	}
	var buf bytes.Buffer
	f := output.NewFormatter(output.FormatText, &buf)
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"aws_instance.web", "instance_type", "t2.micro", "t3.small"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWriteJSON_ValidOutput(t *testing.T) {
	results := []drift.DriftResult{
		makeDrift("aws_s3_bucket.logs", "missing", "", nil, nil),
	}
	var buf bytes.Buffer
	f := output.NewFormatter(output.FormatJSON, &buf)
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded []drift.DriftResult
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(decoded) != 1 || decoded[0].ResourceKey != "aws_s3_bucket.logs" {
		t.Errorf("unexpected decoded result: %+v", decoded)
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(output.Format("xml"), &buf)
	if err := f.Write(nil); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
