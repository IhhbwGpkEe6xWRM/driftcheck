package drift

import "strings"

// IgnoreRule defines a rule for suppressing specific drift findings.
type IgnoreRule struct {
	// ResourceKey matches against "type.name" (supports "*" wildcard prefix/suffix)
	ResourceKey string
	// Attribute is the attribute name to ignore; empty string means ignore all attributes
	Attribute string
}

// RuleSet holds a collection of ignore rules.
type RuleSet struct {
	rules []IgnoreRule
}

// NewRuleSet creates a RuleSet from a slice of IgnoreRules.
func NewRuleSet(rules []IgnoreRule) *RuleSet {
	return &RuleSet{rules: rules}
}

// ShouldIgnore returns true if the given resource key + attribute combination
// is suppressed by any rule in the set.
func (rs *RuleSet) ShouldIgnore(resourceKey, attribute string) bool {
	for _, r := range rs.rules {
		if matchesPattern(r.ResourceKey, resourceKey) {
			if r.Attribute == "" || r.Attribute == attribute {
				return true
			}
		}
	}
	return false
}

// matchesPattern checks whether value matches a simple glob pattern.
// Supports leading and trailing "*" wildcards.
func matchesPattern(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
		mid := pattern[1 : len(pattern)-1]
		return strings.Contains(value, mid)
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(value, pattern[1:])
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(value, pattern[:len(pattern)-1])
	}
	return pattern == value
}
