package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeResult(key, attr string, kind DriftKind) DriftResult {
	return DriftResult{
		ResourceKey: key,
		Attribute:   attr,
		Kind:        kind,
	}
}

func TestSaveAndLoadBaseline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	results := []DriftResult{
		makeResult("aws_instance.web", "instance_type", DriftKindChanged),
		makeResult("aws_s3_bucket.data", "", DriftKindMissing),
	}
	meta := BaselineMeta{StateFile: "terraform.tfstate", Region: "us-east-1"}

	if err := SaveBaseline(path, results, meta); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}

	if len(loaded.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(loaded.Results))
	}
	if loaded.Meta.StateFile != "terraform.tfstate" {
		t.Errorf("unexpected state file: %s", loaded.Meta.StateFile)
	}
	if loaded.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0644)

	_, err := LoadBaseline(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiffBaseline_NewAndResolved(t *testing.T) {
	baseline := &Baseline{
		CreatedAt: time.Now(),
		Results: []DriftResult{
			makeResult("aws_instance.web", "instance_type", DriftKindChanged),
			makeResult("aws_s3_bucket.old", "", DriftKindMissing),
		},
	}
	current := []DriftResult{
		makeResult("aws_instance.web", "instance_type", DriftKindChanged),
		makeResult("aws_instance.api", "ami", DriftKindChanged),
	}

	newDrift, resolved := DiffBaseline(baseline, current)

	if len(newDrift) != 1 || newDrift[0].ResourceKey != "aws_instance.api" {
		t.Errorf("expected 1 new drift for aws_instance.api, got %+v", newDrift)
	}
	if len(resolved) != 1 || resolved[0].ResourceKey != "aws_s3_bucket.old" {
		t.Errorf("expected 1 resolved drift for aws_s3_bucket.old, got %+v", resolved)
	}
}

func TestDiffBaseline_NoDiff(t *testing.T) {
	results := []DriftResult{
		makeResult("aws_instance.web", "instance_type", DriftKindChanged),
	}
	baseline := &Baseline{CreatedAt: time.Now(), Results: results}

	newDrift, resolved := DiffBaseline(baseline, results)

	if len(newDrift) != 0 || len(resolved) != 0 {
		t.Errorf("expected no diff, got new=%d resolved=%d", len(newDrift), len(resolved))
	}
}
