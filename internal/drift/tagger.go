package drift

import (
	"sort"
	"strings"
)

// Tag represents a label attached to a drift result for categorization.
type Tag struct {
	Key   string
	Value string
}

// TaggedResult wraps an AnnotatedResult with user-defined tags.
type TaggedResult struct {
	AnnotatedResult
	Tags []Tag
}

// Tagger applies tags to annotated drift results based on configurable rules.
type Tagger struct {
	rules []TagRule
}

// TagRule maps a resource key pattern and optional attribute to a tag.
type TagRule struct {
	ResourcePattern string
	Attribute       string
	TagKey          string
	TagValue        string
}

// NewTagger creates a Tagger with the given rules.
func NewTagger(rules []TagRule) *Tagger {
	return &Tagger{rules: rules}
}

// Apply tags all results according to the configured rules.
func (t *Tagger) Apply(results []AnnotatedResult) []TaggedResult {
	tagged := make([]TaggedResult, 0, len(results))
	for _, r := range results {
		tr := TaggedResult{AnnotatedResult: r}
		for _, rule := range t.rules {
			if t.matches(rule, r) {
				tr.Tags = appendUniqueTag(tr.Tags, Tag{Key: rule.TagKey, Value: rule.TagValue})
			}
		}
		tagged = append(tagged, tr)
	}
	return tagged
}

// GroupByTag groups TaggedResults by the value of a specific tag key.
// Results without the tag are placed under the "" (empty) key.
func GroupByTag(results []TaggedResult, tagKey string) map[string][]TaggedResult {
	groups := make(map[string][]TaggedResult)
	for _, r := range results {
		val := ""
		for _, tag := range r.Tags {
			if tag.Key == tagKey {
				val = tag.Value
				break
			}
		}
		groups[val] = append(groups[val], r)
	}
	return groups
}

// SortedTagKeys returns the keys of a GroupByTag result in sorted order.
func SortedTagKeys(groups map[string][]TaggedResult) []string {
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (t *Tagger) matches(rule TagRule, r AnnotatedResult) bool {
	if !matchesPattern(rule.ResourcePattern, r.Resource) {
		return false
	}
	if rule.Attribute == "" || rule.Attribute == "*" {
		return true
	}
	for _, ch := range r.Changes {
		if strings.EqualFold(ch.Attribute, rule.Attribute) {
			return true
		}
	}
	return false
}

func appendUniqueTag(tags []Tag, t Tag) []Tag {
	for _, existing := range tags {
		if existing.Key == t.Key && existing.Value == t.Value {
			return tags
		}
	}
	return append(tags, t)
}
