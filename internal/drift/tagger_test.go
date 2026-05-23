package drift

import (
	"testing"
)

func makeAnnotatedForTag(resource, kind string, attrs []string) AnnotatedResult {
	changes := make([]AnnotatedChange, 0, len(attrs))
	for _, a := range attrs {
		changes = append(changes, AnnotatedChange{Attribute: a, Severity: SeverityLow})
	}
	return AnnotatedResult{
		DriftResult: DriftResult{Resource: resource, Kind: kind},
		Changes:     changes,
	}
}

func TestApplyTags_NoRules(t *testing.T) {
	tagger := NewTagger(nil)
	results := []AnnotatedResult{makeAnnotatedForTag("aws_s3_bucket.logs", "changed", []string{"acl"})}
	tagged := tagger.Apply(results)
	if len(tagged) != 1 {
		t.Fatalf("expected 1 result, got %d", len(tagged))
	}
	if len(tagged[0].Tags) != 0 {
		t.Errorf("expected no tags, got %v", tagged[0].Tags)
	}
}

func TestApplyTags_ExactResourceMatch(t *testing.T) {
	rules := []TagRule{
		{ResourcePattern: "aws_s3_bucket.logs", TagKey: "team", TagValue: "platform"},
	}
	tagger := NewTagger(rules)
	results := []AnnotatedResult{makeAnnotatedForTag("aws_s3_bucket.logs", "changed", nil)}
	tagged := tagger.Apply(results)
	if len(tagged[0].Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tagged[0].Tags))
	}
	if tagged[0].Tags[0].Value != "platform" {
		t.Errorf("unexpected tag value: %s", tagged[0].Tags[0].Value)
	}
}

func TestApplyTags_WildcardResourceMatch(t *testing.T) {
	rules := []TagRule{
		{ResourcePattern: "aws_s3_bucket.*", TagKey: "service", TagValue: "storage"},
	}
	tagger := NewTagger(rules)
	results := []AnnotatedResult{
		makeAnnotatedForTag("aws_s3_bucket.logs", "changed", nil),
		makeAnnotatedForTag("aws_instance.web", "changed", nil),
	}
	tagged := tagger.Apply(results)
	if len(tagged[0].Tags) != 1 {
		t.Errorf("expected tag on s3 bucket, got %d tags", len(tagged[0].Tags))
	}
	if len(tagged[1].Tags) != 0 {
		t.Errorf("expected no tag on instance, got %d tags", len(tagged[1].Tags))
	}
}

func TestApplyTags_AttributeFilter(t *testing.T) {
	rules := []TagRule{
		{ResourcePattern: "*", Attribute: "ami", TagKey: "drift-type", TagValue: "image"},
	}
	tagger := NewTagger(rules)
	results := []AnnotatedResult{
		makeAnnotatedForTag("aws_instance.web", "changed", []string{"ami"}),
		makeAnnotatedForTag("aws_instance.db", "changed", []string{"instance_type"}),
	}
	tagged := tagger.Apply(results)
	if len(tagged[0].Tags) != 1 {
		t.Errorf("expected tag for ami change, got %d", len(tagged[0].Tags))
	}
	if len(tagged[1].Tags) != 0 {
		t.Errorf("expected no tag for instance_type change, got %d", len(tagged[1].Tags))
	}
}

func TestGroupByTag_Basic(t *testing.T) {
	results := []TaggedResult{
		{AnnotatedResult: makeAnnotatedForTag("aws_s3_bucket.a", "changed", nil), Tags: []Tag{{Key: "team", Value: "platform"}}},
		{AnnotatedResult: makeAnnotatedForTag("aws_s3_bucket.b", "changed", nil), Tags: []Tag{{Key: "team", Value: "platform"}}},
		{AnnotatedResult: makeAnnotatedForTag("aws_instance.x", "changed", nil), Tags: []Tag{{Key: "team", Value: "infra"}}},
		{AnnotatedResult: makeAnnotatedForTag("aws_rds.y", "missing", nil)},
	}
	groups := GroupByTag(results, "team")
	if len(groups["platform"]) != 2 {
		t.Errorf("expected 2 platform results, got %d", len(groups["platform"]))
	}
	if len(groups["infra"]) != 1 {
		t.Errorf("expected 1 infra result, got %d", len(groups["infra"]))
	}
	if len(groups[""]) != 1 {
		t.Errorf("expected 1 untagged result, got %d", len(groups[""]))
	}
}

func TestSortedTagKeys_Order(t *testing.T) {
	groups := map[string][]TaggedResult{
		"zebra": {},
		"alpha": {},
		"mango": {},
	}
	keys := SortedTagKeys(groups)
	expected := []string{"alpha", "mango", "zebra"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], k)
		}
	}
}
