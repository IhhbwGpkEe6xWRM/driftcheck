package drift

import (
	"testing"
)

func TestShouldIgnore_ExactMatch(t *testing.T) {
	rs := NewRuleSet([]IgnoreRule{
		{ResourceKey: "aws_instance.web", Attribute: "tags"},
	})
	if !rs.ShouldIgnore("aws_instance.web", "tags") {
		t.Error("expected tags on aws_instance.web to be ignored")
	}
	if rs.ShouldIgnore("aws_instance.web", "ami") {
		t.Error("expected ami on aws_instance.web NOT to be ignored")
	}
}

func TestShouldIgnore_AllAttributes(t *testing.T) {
	rs := NewRuleSet([]IgnoreRule{
		{ResourceKey: "aws_s3_bucket.logs", Attribute: ""},
	})
	if !rs.ShouldIgnore("aws_s3_bucket.logs", "acl") {
		t.Error("expected all attributes on aws_s3_bucket.logs to be ignored")
	}
	if !rs.ShouldIgnore("aws_s3_bucket.logs", "versioning") {
		t.Error("expected all attributes on aws_s3_bucket.logs to be ignored")
	}
}

func TestShouldIgnore_WildcardResourceKey(t *testing.T) {
	rs := NewRuleSet([]IgnoreRule{
		{ResourceKey: "aws_instance.*", Attribute: "tags"},
	})
	if !rs.ShouldIgnore("aws_instance.web", "tags") {
		t.Error("expected wildcard to match aws_instance.web")
	}
	if !rs.ShouldIgnore("aws_instance.api", "tags") {
		t.Error("expected wildcard to match aws_instance.api")
	}
	if rs.ShouldIgnore("aws_s3_bucket.data", "tags") {
		t.Error("expected wildcard NOT to match aws_s3_bucket.data")
	}
}

func TestShouldIgnore_GlobalWildcard(t *testing.T) {
	rs := NewRuleSet([]IgnoreRule{
		{ResourceKey: "*", Attribute: "last_modified"},
	})
	if !rs.ShouldIgnore("aws_instance.web", "last_modified") {
		t.Error("expected global wildcard to match any resource")
	}
	if rs.ShouldIgnore("aws_instance.web", "ami") {
		t.Error("expected global wildcard NOT to suppress unrelated attribute")
	}
}

func TestShouldIgnore_NoRules(t *testing.T) {
	rs := NewRuleSet(nil)
	if rs.ShouldIgnore("aws_instance.web", "ami") {
		t.Error("expected empty ruleset to suppress nothing")
	}
}

func TestMatchesPattern_ContainsWildcard(t *testing.T) {
	if !matchesPattern("*instance*", "aws_instance.web") {
		t.Error("expected *instance* to match aws_instance.web")
	}
	if matchesPattern("*instance*", "aws_s3_bucket.data") {
		t.Error("expected *instance* NOT to match aws_s3_bucket.data")
	}
}
