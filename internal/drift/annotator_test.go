package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeAnnotatorResult(key string, kind DriftKind, attrs []AttributeDiff) DriftResult {
	return DriftResult{
		ResourceKey: key,
		Kind:        kind,
		Attributes:  attrs,
	}
}

func TestAnnotate_MissingResource(t *testing.T) {
	results := []DriftResult{
		makeAnnotatorResult("aws_instance.web", DriftKindMissing, nil),
	}
	annotated := Annotate(results)
	require.Len(t, annotated, 1)
	assert.Equal(t, SeverityCritical, annotated[0].Annotation.Severity)
	assert.Contains(t, annotated[0].Annotation.Message, "aws_instance.web")
	assert.NotEmpty(t, annotated[0].Annotation.Hint)
}

func TestAnnotate_ExtraResource(t *testing.T) {
	results := []DriftResult{
		makeAnnotatorResult("aws_s3_bucket.logs", DriftKindExtra, nil),
	}
	annotated := Annotate(results)
	require.Len(t, annotated, 1)
	assert.Equal(t, SeverityInfo, annotated[0].Annotation.Severity)
	assert.Contains(t, annotated[0].Annotation.Message, "aws_s3_bucket.logs")
}

func TestAnnotate_ChangedSingleAttribute(t *testing.T) {
	attrs := []AttributeDiff{{Key: "instance_type", StateVal: "t2.micro", CloudVal: "t3.medium"}}
	results := []DriftResult{
		makeAnnotatorResult("aws_instance.api", DriftKindChanged, attrs),
	}
	annotated := Annotate(results)
	require.Len(t, annotated, 1)
	assert.Equal(t, SeverityLow, annotated[0].Annotation.Severity)
}

func TestAnnotate_ChangedCriticalAttribute(t *testing.T) {
	attrs := []AttributeDiff{{Key: "kms_key_id", StateVal: "old-key", CloudVal: "new-key"}}
	results := []DriftResult{
		makeAnnotatorResult("aws_s3_bucket.data", DriftKindChanged, attrs),
	}
	annotated := Annotate(results)
	require.Len(t, annotated, 1)
	assert.Equal(t, SeverityCritical, annotated[0].Annotation.Severity)
}

func TestAnnotate_ChangedManyAttributes(t *testing.T) {
	attrs := []AttributeDiff{
		{Key: "tags", StateVal: "a", CloudVal: "b"},
		{Key: "description", StateVal: "old", CloudVal: "new"},
		{Key: "name", StateVal: "x", CloudVal: "y"},
		{Key: "size", StateVal: "10", CloudVal: "20"},
	}
	results := []DriftResult{
		makeAnnotatorResult("aws_security_group.main", DriftKindChanged, attrs),
	}
	annotated := Annotate(results)
	require.Len(t, annotated, 1)
	assert.Equal(t, SeverityHigh, annotated[0].Annotation.Severity)
}

func TestAnnotate_Empty(t *testing.T) {
	annotated := Annotate([]DriftResult{})
	assert.Empty(t, annotated)
}

func TestAnnotate_PreservesResult(t *testing.T) {
	r := makeAnnotatorResult("aws_instance.web", DriftKindMissing, nil)
	annotated := Annotate([]DriftResult{r})
	require.Len(t, annotated, 1)
	assert.Equal(t, r, annotated[0].Result)
}
