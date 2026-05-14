package drift

// DriftKind classifies the type of drift detected for a resource.
type DriftKind string

const (
	// KindMissing indicates the resource exists in Terraform state but not in the cloud.
	KindMissing DriftKind = "missing"
	// KindChanged indicates the resource exists in both but attributes differ.
	KindChanged DriftKind = "changed"
	// KindExtra indicates the resource exists in the cloud but not in Terraform state.
	KindExtra DriftKind = "extra"
)

// AttributeDiff records the expected vs actual value for a single attribute.
type AttributeDiff struct {
	Attribute string
	Expected  interface{}
	Actual    interface{}
}

// DriftResult represents a single drift finding for a resource.
type DriftResult struct {
	// ResourceKey is "resource_type.resource_name", e.g. "aws_instance.web".
	ResourceKey string
	// Kind classifies the nature of the drift.
	Kind DriftKind
	// Diffs holds per-attribute differences; only populated for KindChanged.
	Diffs []AttributeDiff
}

// HasDiffs returns true when the result contains attribute-level differences.
func (r DriftResult) HasDiffs() bool {
	return len(r.Diffs) > 0
}
