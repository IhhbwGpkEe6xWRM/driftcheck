package drift

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func makeAnnotatedForExport(key, kind, attr, expected, actual, severity string, score int) AnnotatedResult {
	r := makeAnnotatorResult(key, DriftKind(kind))
	if attr != "" {
		r.ChangedAttributes = []AttributeChange{
			{Attribute: attr, Expected: expected, Actual: actual},
		}
	}
	return AnnotatedResult{
		Result:   r,
		Severity: Severity(severity),
		Score:    score,
	}
}

func TestExportResults_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := ExportResults(&buf, nil, ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("error should mention format, got: %v", err)
	}
}

func TestExportResults_CSV_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	err := ExportResults(&buf, []AnnotatedResult{}, ExportCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header row, got %d lines", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
}

func TestExportResults_CSV_WithChangedAttribute(t *testing.T) {
	results := []AnnotatedResult{
		makeAnnotatedForExport("aws_instance.web", "changed", "instance_type", "t2.micro", "t3.small", "medium", 50),
	}
	var buf bytes.Buffer
	err := ExportResults(&buf, results, ExportCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_instance.web") {
		t.Errorf("expected resource key in output")
	}
	if !strings.Contains(out, "instance_type") {
		t.Errorf("expected attribute name in output")
	}
	if !strings.Contains(out, "t3.small") {
		t.Errorf("expected actual value in output")
	}
}

func TestExportResults_JSON_ValidOutput(t *testing.T) {
	results := []AnnotatedResult{
		makeAnnotatedForExport("aws_s3_bucket.logs", "missing", "", "", "", "high", 80),
	}
	var buf bytes.Buffer
	err := ExportResults(&buf, results, ExportJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].ResourceKey != "aws_s3_bucket.logs" {
		t.Errorf("unexpected resource key: %s", records[0].ResourceKey)
	}
	if records[0].Score != 80 {
		t.Errorf("expected score 80, got %d", records[0].Score)
	}
}
