package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/user/driftcheck/internal/drift"
)

func makeAnnotated(key string, kind drift.DriftKind, attrs []drift.AttributeDiff) drift.AnnotatedDrift {
	results := drift.Annotate([]drift.DriftResult{{
		ResourceKey: key,
		Kind:        kind,
		Attributes:  attrs,
	}})
	return results[0]
}

func TestWriteAnnotatedText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	err := WriteAnnotatedText(&buf, []drift.AnnotatedDrift{})
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "No drift detected")
}

func TestWriteAnnotatedText_MissingResource(t *testing.T) {
	items := []drift.AnnotatedDrift{
		makeAnnotated("aws_instance.web", drift.DriftKindMissing, nil),
	}
	var buf bytes.Buffer
	err := WriteAnnotatedText(&buf, items)
	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "[CRITICAL]")
	assert.Contains(t, out, "aws_instance.web")
	assert.Contains(t, out, "Hint:")
}

func TestWriteAnnotatedText_ChangedAttributes(t *testing.T) {
	attrs := []drift.AttributeDiff{
		{Key: "instance_type", StateVal: "t2.micro", CloudVal: "t3.medium"},
	}
	items := []drift.AnnotatedDrift{
		makeAnnotated("aws_instance.api", drift.DriftKindChanged, attrs),
	}
	var buf bytes.Buffer
	err := WriteAnnotatedText(&buf, items)
	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "instance_type")
	assert.Contains(t, out, "t2.micro")
	assert.Contains(t, out, "t3.medium")
}

func TestWriteAnnotatedJSON_ValidOutput(t *testing.T) {
	items := []drift.AnnotatedDrift{
		makeAnnotated("aws_s3_bucket.data", drift.DriftKindExtra, nil),
	}
	var buf bytes.Buffer
	err := WriteAnnotatedJSON(&buf, items)
	require.NoError(t, err)

	var out []drift.AnnotatedDrift
	err = json.Unmarshal(buf.Bytes(), &out)
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "aws_s3_bucket.data", out[0].Result.ResourceKey)
	assert.Equal(t, drift.SeverityInfo, out[0].Annotation.Severity)
}

func TestWriteAnnotatedJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := WriteAnnotatedJSON(&buf, []drift.AnnotatedDrift{})
	require.NoError(t, err)
	assert.True(t, strings.TrimSpace(buf.String()) == "[]")
}
