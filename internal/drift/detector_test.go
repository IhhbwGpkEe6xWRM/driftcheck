package drift_test

import (
	"testing"

	"github.com/driftcheck/internal/drift"
	"github.com/driftcheck/internal/tfstate"
)

func makeResource(rType, rName string, attrs map[string]interface{}) tfstate.Resource {
	return tfstate.Resource{
		Type: rType,
		Name: rName,
		Instances: []tfstate.ResourceInstance{
			{Attributes: attrs},
		},
	}
}

func TestDetect_NoDrift(t *testing.T) {
	resources := []tfstate.Resource{
		makeResource("aws_instance", "web", map[string]interface{}{"instance_type": "t2.micro"}),
	}
	live := drift.LiveResourceMap{
		"aws_instance.web": {"instance_type": "t2.micro"},
	}

	results := drift.Detect(resources, live)
	if len(results) != 0 {
		t.Errorf("expected no drift, got %d result(s): %v", len(results), results)
	}
}

func TestDetect_MissingResource(t *testing.T) {
	resources := []tfstate.Resource{
		makeResource("aws_s3_bucket", "logs", map[string]interface{}{"bucket": "my-logs"}),
	}
	live := drift.LiveResourceMap{}

	results := drift.Detect(resources, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].DriftType != drift.DriftTypeMissing {
		t.Errorf("expected DriftTypeMissing, got %s", results[0].DriftType)
	}
}

func TestDetect_ChangedAttribute(t *testing.T) {
	resources := []tfstate.Resource{
		makeResource("aws_instance", "app", map[string]interface{}{"instance_type": "t2.micro"}),
	}
	live := drift.LiveResourceMap{
		"aws_instance.app": {"instance_type": "t2.large"},
	}

	results := drift.Detect(resources, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].DriftType != drift.DriftTypeChanged {
		t.Errorf("expected DriftTypeChanged, got %s", results[0].DriftType)
	}
	if results[0].Attribute != "instance_type" {
		t.Errorf("expected attribute 'instance_type', got '%s'", results[0].Attribute)
	}
}

func TestDetect_ExtraResource(t *testing.T) {
	resources := []tfstate.Resource{}
	live := drift.LiveResourceMap{
		"aws_security_group.allow_ssh": {"name": "allow_ssh"},
	}

	results := drift.Detect(resources, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].DriftType != drift.DriftTypeExtra {
		t.Errorf("expected DriftTypeExtra, got %s", results[0].DriftType)
	}
}

func TestDriftResult_String(t *testing.T) {
	d := drift.DriftResult{
		ResourceType: "aws_instance",
		ResourceName: "web",
		DriftType:    drift.DriftTypeChanged,
		Attribute:    "ami",
		StateValue:   "ami-abc",
		LiveValue:    "ami-xyz",
	}
	s := d.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
