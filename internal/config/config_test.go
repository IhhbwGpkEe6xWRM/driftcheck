package config

import (
	"testing"
)

const sampleConfigYAML = `
statefile: terraform.tfstate
region: us-west-2
profile: dev
ignore:
  - resource: aws_instance.web
    attribute: tags
  - resource: "*"
    attribute: last_modified
output:
  format: json
  file: drift-report.json
`

func TestParse_ValidConfig(t *testing.T) {
	cfg, err := parse([]byte(sampleConfigYAML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Statefile != "terraform.tfstate" {
		t.Errorf("statefile: got %q, want %q", cfg.Statefile, "terraform.tfstate")
	}
	if cfg.Region != "us-west-2" {
		t.Errorf("region: got %q, want %q", cfg.Region, "us-west-2")
	}
	if cfg.Profile != "dev" {
		t.Errorf("profile: got %q, want %q", cfg.Profile, "dev")
	}
	if len(cfg.Ignore) != 2 {
		t.Fatalf("ignore rules: got %d, want 2", len(cfg.Ignore))
	}
	if cfg.Ignore[0].Resource != "aws_instance.web" {
		t.Errorf("ignore[0].resource: got %q", cfg.Ignore[0].Resource)
	}
	if cfg.Output.Format != "json" {
		t.Errorf("output.format: got %q, want json", cfg.Output.Format)
	}
	if cfg.Output.File != "drift-report.json" {
		t.Errorf("output.file: got %q", cfg.Output.File)
	}
}

func TestParse_DefaultsApplied(t *testing.T) {
	cfg, err := parse([]byte("statefile: state.tfstate\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Output.Format != "text" {
		t.Errorf("default format: got %q, want text", cfg.Output.Format)
	}
	if cfg.Region != "us-east-1" {
		t.Errorf("default region: got %q, want us-east-1", cfg.Region)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	_, err := parse([]byte(":::invalid yaml:::"))
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/driftcheck.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
