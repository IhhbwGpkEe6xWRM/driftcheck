package tfstate

import (
	"encoding/json"
	"testing"
)

func sampleStateJSON(version int) []byte {
	raw := tfStateFile{
		Version: version,
		Resources: []tfStateResource{
			{
				Type:     "aws_instance",
				Name:     "web",
				Provider: "provider[\"registry.terraform.io/hashicorp/aws\"]",
				Instances: []tfStateInstance{
					{
						Attributes: map[string]interface{}{
							"id":            "i-0abc123",
							"instance_type": "t3.micro",
							"ami":           "ami-0deadbeef",
						},
					},
				},
			},
		},
	}
	b, _ := json.Marshal(raw)
	return b
}

func TestParse_ValidState(t *testing.T) {
	state, err := Parse(sampleStateJSON(4))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if state.Version != 4 {
		t.Errorf("expected version 4, got %d", state.Version)
	}
	if len(state.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(state.Resources))
	}
	res := state.Resources[0]
	if res.Type != "aws_instance" {
		t.Errorf("expected type aws_instance, got %s", res.Type)
	}
	if res.Name != "web" {
		t.Errorf("expected name web, got %s", res.Name)
	}
	if res.Attributes["instance_type"] != "t3.micro" {
		t.Errorf("unexpected instance_type: %v", res.Attributes["instance_type"])
	}
}

func TestParse_UnsupportedVersion(t *testing.T) {
	_, err := Parse(sampleStateJSON(3))
	if err == nil {
		t.Fatal("expected error for unsupported version, got nil")
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := Parse([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParse_ResourceWithNoInstances(t *testing.T) {
	raw := tfStateFile{
		Version: 4,
		Resources: []tfStateResource{
			{Type: "aws_s3_bucket", Name: "empty", Instances: []tfStateInstance{}},
		},
	}
	b, _ := json.Marshal(raw)
	state, err := Parse(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(state.Resources) != 0 {
		t.Errorf("expected 0 resources for empty instances, got %d", len(state.Resources))
	}
}
