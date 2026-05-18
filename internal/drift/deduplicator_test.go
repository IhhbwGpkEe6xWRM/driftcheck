package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeDedupResult(rType, rName string, kind DriftKind, attrs map[string]AttributeDiff) DriftResult {
	return DriftResult{
		ResourceType:      rType,
		ResourceName:      rName,
		Kind:              kind,
		ChangedAttributes: attrs,
	}
}

func TestDeduplicateResults_NoDuplicates(t *testing.T) {
	input := []DriftResult{
		makeDedupResult("aws_s3_bucket", "my-bucket", DriftKindChanged, map[string]AttributeDiff{
			"region": {TFValue: "us-east-1", LiveValue: "us-west-2"},
		}),
		makeDedupResult("aws_s3_bucket", "other-bucket", DriftKindMissing, nil),
	}

	result := DeduplicateResults(input)
	assert.Len(t, result, 2)
}

func TestDeduplicateResults_RemovesDuplicates(t *testing.T) {
	entry := makeDedupResult("aws_s3_bucket", "my-bucket", DriftKindChanged, map[string]AttributeDiff{
		"region": {TFValue: "us-east-1", LiveValue: "us-west-2"},
	})
	input := []DriftResult{entry, entry, entry}

	result := DeduplicateResults(input)
	assert.Len(t, result, 1)
}

func TestDeduplicateResults_PreservesOrder(t *testing.T) {
	a := makeDedupResult("aws_s3_bucket", "alpha", DriftKindMissing, nil)
	b := makeDedupResult("aws_instance", "beta", DriftKindExtra, nil)
	input := []DriftResult{a, b, a}

	result := DeduplicateResults(input)
	assert.Len(t, result, 2)
	assert.Equal(t, "alpha", result[0].ResourceName)
	assert.Equal(t, "beta", result[1].ResourceName)
}

func TestDeduplicateByResource_MergesAttributes(t *testing.T) {
	r1 := makeDedupResult("aws_s3_bucket", "my-bucket", DriftKindChanged, map[string]AttributeDiff{
		"region": {TFValue: "us-east-1", LiveValue: "us-west-2"},
	})
	r2 := makeDedupResult("aws_s3_bucket", "my-bucket", DriftKindChanged, map[string]AttributeDiff{
		"acl": {TFValue: "private", LiveValue: "public-read"},
	})

	result := DeduplicateByResource([]DriftResult{r1, r2})
	assert.Len(t, result, 1)
	assert.Len(t, result[0].ChangedAttributes, 2)
	assert.Contains(t, result[0].ChangedAttributes, "region")
	assert.Contains(t, result[0].ChangedAttributes, "acl")
}

func TestDeduplicateByResource_DistinctResourcesKept(t *testing.T) {
	r1 := makeDedupResult("aws_s3_bucket", "bucket-a", DriftKindChanged, nil)
	r2 := makeDedupResult("aws_s3_bucket", "bucket-b", DriftKindChanged, nil)
	r3 := makeDedupResult("aws_instance", "bucket-a", DriftKindMissing, nil)

	result := DeduplicateByResource([]DriftResult{r1, r2, r3})
	assert.Len(t, result, 3)
}

func TestDeduplicateResults_EmptyInput(t *testing.T) {
	result := DeduplicateResults([]DriftResult{})
	assert.Empty(t, result)
}
