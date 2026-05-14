// Package tfstate provides types and utilities for parsing and working with
// Terraform state files.
package tfstate

// SupportedVersion is the Terraform state format version this package supports.
const SupportedVersion = 4

// State represents the top-level structure of a Terraform state file.
type State struct {
	// Version is the state file format version.
	Version int `json:"version"`

	// TerraformVersion is the version of Terraform that wrote this state.
	TerraformVersion string `json:"terraform_version"`

	// Serial is incremented on each state write.
	Serial int `json:"serial"`

	// Lineage is a unique ID assigned to a state when it is created.
	Lineage string `json:"lineage"`

	// Outputs contains the root module output values.
	Outputs map[string]Output `json:"outputs"`

	// Resources is the list of resources tracked in this state.
	Resources []Resource `json:"resources"`
}

// Output represents a single Terraform output value.
type Output struct {
	// Value is the output value. It can be any JSON-compatible type.
	Value interface{} `json:"value"`

	// Type is the Terraform type of the output value.
	Type interface{} `json:"type"`

	// Sensitive indicates whether the output value is marked as sensitive.
	Sensitive bool `json:"sensitive"`
}

// Resource represents a single resource tracked in Terraform state.
type Resource struct {
	// Module is the module path where this resource is defined (empty for root).
	Module string `json:"module,omitempty"`

	// Mode is either "managed" or "data".
	Mode string `json:"mode"`

	// Type is the resource type, e.g. "aws_instance".
	Type string `json:"type"`

	// Name is the resource name as defined in configuration.
	Name string `json:"name"`

	// Provider is the provider source address, e.g. "provider[\"registry.terraform.io/hashicorp/aws\"]".
	Provider string `json:"provider"`

	// Instances holds the individual resource instances (one per count/for_each key).
	Instances []Instance `json:"instances"`
}

// Address returns a human-readable address for the resource, including the
// module prefix when the resource is not in the root module.
func (r Resource) Address() string {
	addr := r.Type + "." + r.Name
	if r.Module != "" {
		return r.Module + "." + addr
	}
	return addr
}

// Instance represents a single instance of a Terraform resource.
type Instance struct {
	// IndexKey is the count index or for_each key for this instance.
	// It is nil for resources without count or for_each.
	IndexKey interface{} `json:"index_key,omitempty"`

	// Status is the instance status, typically "tainted" or empty.
	Status string `json:"status,omitempty"`

	// Attributes holds the current attribute values for this instance.
	Attributes map[string]interface{} `json:"attributes"`

	// Dependencies lists the addresses of resources this instance depends on.
	Dependencies []string `json:"dependencies,omitempty"`

	// SchemaVersion is the schema version for the resource type at the time
	// the state was written.
	SchemaVersion int `json:"schema_version"`
}
