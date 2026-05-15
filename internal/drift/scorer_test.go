package drift

import (
	"testing"
)

func TestScoreDrifts_MissingResource(t *testing.T) {
	results := []DriftResult{
		{Kind: DriftMissing, ResourceKey: "aws_instance.web"},
	}
	scored := ScoreDrifts(results, nil)
	if len(scored) != 1 {
		t.Fatalf("expected 1 scored drift, got %d", len(scored))
	}
	if scored[0].Severity != SeverityHigh {
		t.Errorf("expected High severity for missing resource, got %s", scored[0].Severity)
	}
}

func TestScoreDrifts_ExtraResource(t *testing.T) {
	results := []DriftResult{
		{Kind: DriftExtra, ResourceKey: "aws_instance.ghost"},
	}
	scored := ScoreDrifts(results, nil)
	if scored[0].Severity != SeverityLow {
		t.Errorf("expected Low severity for extra resource, got %s", scored[0].Severity)
	}
}

func TestScoreDrifts_CriticalAttribute(t *testing.T) {
	results := []DriftResult{
		{
			Kind:              DriftChanged,
			ResourceKey:       "aws_db_instance.main",
			ChangedAttributes: []string{"master_password"},
		},
	}
	scored := ScoreDrifts(results, nil)
	if scored[0].Severity != SeverityCritical {
		t.Errorf("expected Critical severity for password attribute, got %s", scored[0].Severity)
	}
}

func TestScoreDrifts_MediumAttribute(t *testing.T) {
	results := []DriftResult{
		{
			Kind:              DriftChanged,
			ResourceKey:       "aws_security_group.web",
			ChangedAttributes: []string{"ingress_port"},
		},
	}
	scored := ScoreDrifts(results, nil)
	if scored[0].Severity != SeverityMedium {
		t.Errorf("expected Medium severity for port attribute, got %s", scored[0].Severity)
	}
}

func TestScoreDrifts_LowAttribute(t *testing.T) {
	results := []DriftResult{
		{
			Kind:              DriftChanged,
			ResourceKey:       "aws_instance.web",
			ChangedAttributes: []string{"tags"},
		},
	}
	scored := ScoreDrifts(results, nil)
	if scored[0].Severity != SeverityLow {
		t.Errorf("expected Low severity for tags attribute, got %s", scored[0].Severity)
	}
}

func TestScoreDrifts_CustomRules(t *testing.T) {
	customRules := []ScoringRule{
		{AttributePattern: "tags", Severity: SeverityHigh},
	}
	results := []DriftResult{
		{
			Kind:              DriftChanged,
			ResourceKey:       "aws_instance.web",
			ChangedAttributes: []string{"tags"},
		},
	}
	scored := ScoreDrifts(results, customRules)
	if scored[0].Severity != SeverityHigh {
		t.Errorf("expected High severity with custom rule, got %s", scored[0].Severity)
	}
}

func TestScoreDrifts_Empty(t *testing.T) {
	scored := ScoreDrifts([]DriftResult{}, nil)
	if len(scored) != 0 {
		t.Errorf("expected empty result, got %d", len(scored))
	}
}
