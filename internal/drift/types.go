package drift

// DriftKind describes the category of drift detected for a resource or attribute.
type DriftKind string

const (
	// DriftKindMissing indicates the resource exists in Terraform state but
	// was not found in the live cloud environment.
	DriftKindMissing DriftKind = "missing"

	// DriftKindChanged indicates one or more attributes differ between the
	// Terraform state and the live cloud resource.
	DriftKindChanged DriftKind = "changed"

	// DriftKindExtra indicates a resource exists in the cloud but is not
	// tracked in Terraform state.
	DriftKindExtra DriftKind = "extra"
)

// DriftResult captures a single unit of detected drift.
type DriftResult struct {
	// ResourceKey is the fully qualified Terraform resource address,
	// e.g. "aws_instance.web".
	ResourceKey string `json:"resource_key"`

	// ResourceType is the Terraform resource type, e.g. "aws_instance".
	ResourceType string `json:"resource_type"`

	// Kind describes the category of drift.
	Kind DriftKind `json:"kind"`

	// Attribute is the specific attribute that drifted (empty for missing/extra).
	Attribute string `json:"attribute,omitempty"`

	// StateValue is the value recorded in Terraform state.
	StateValue interface{} `json:"state_value,omitempty"`

	// LiveValue is the value observed in the live cloud resource.
	LiveValue interface{} `json:"live_value,omitempty"`
}
