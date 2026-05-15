package drift

import "fmt"

// Severity levels for drift annotations.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// Annotation attaches human-readable context to a DriftResult.
type Annotation struct {
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Hint     string   `json:"hint,omitempty"`
}

// AnnotatedDrift pairs a DriftResult with its computed Annotation.
type AnnotatedDrift struct {
	Result     DriftResult `json:"result"`
	Annotation Annotation  `json:"annotation"`
}

// Annotate enriches a slice of DriftResults with severity and hints.
func Annotate(results []DriftResult) []AnnotatedDrift {
	annotated := make([]AnnotatedDrift, 0, len(results))
	for _, r := range results {
		annotated = append(annotated, AnnotatedDrift{
			Result:     r,
			Annotation: annotate(r),
		})
	}
	return annotated
}

func annotate(r DriftResult) Annotation {
	switch r.Kind {
	case DriftKindMissing:
		return Annotation{
			Severity: SeverityCritical,
			Message:  fmt.Sprintf("Resource %q exists in state but was not found in the cloud.", r.ResourceKey),
			Hint:     "The resource may have been deleted outside of Terraform. Run 'terraform refresh' or import the resource.",
		}
	case DriftKindExtra:
		return Annotation{
			Severity: SeverityInfo,
			Message:  fmt.Sprintf("Resource %q exists in the cloud but is not tracked in state.", r.ResourceKey),
			Hint:     "Consider importing this resource with 'terraform import'.",
		}
	case DriftKindChanged:
		sev := severityForAttributes(r.Attributes)
		return Annotation{
			Severity: sev,
			Message:  fmt.Sprintf("Resource %q has %d changed attribute(s).", r.ResourceKey, len(r.Attributes)),
			Hint:     "Review the attribute differences and apply or update your Terraform configuration.",
		}
	default:
		return Annotation{
			Severity: SeverityInfo,
			Message:  fmt.Sprintf("Resource %q has an unknown drift kind.", r.ResourceKey),
		}
	}
}

func severityForAttributes(attrs []AttributeDiff) Severity {
	for _, a := range attrs {
		if isCriticalAttribute(a.Key) {
			return SeverityCritical
		}
	}
	if len(attrs) > 3 {
		return SeverityHigh
	}
	if len(attrs) > 1 {
		return SeverityMedium
	}
	return SeverityLow
}

func isCriticalAttribute(key string) bool {
	critical := []string{"iam_instance_profile", "role_arn", "kms_key_id", "encryption", "policy"}
	for _, c := range critical {
		if key == c {
			return true
		}
	}
	return false
}
