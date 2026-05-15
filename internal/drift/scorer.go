package drift

import "fmt"

// SeverityLevel represents the severity of a drift finding.
type SeverityLevel string

const (
	SeverityLow      SeverityLevel = "low"
	SeverityMedium   SeverityLevel = "medium"
	SeverityHigh     SeverityLevel = "high"
	SeverityCritical SeverityLevel = "critical"
)

// ScoredDrift wraps a DriftResult with an assigned severity.
type ScoredDrift struct {
	Result   DriftResult
	Severity SeverityLevel
	Reason   string
}

// ScoringRule maps attribute name patterns to severity levels.
type ScoringRule struct {
	AttributePattern string
	Severity         SeverityLevel
}

var defaultScoringRules = []ScoringRule{
	{AttributePattern: "password", Severity: SeverityCritical},
	{AttributePattern: "secret", Severity: SeverityCritical},
	{AttributePattern: "token", Severity: SeverityHigh},
	{AttributePattern: "iam", Severity: SeverityHigh},
	{AttributePattern: "policy", Severity: SeverityMedium},
	{AttributePattern: "port", Severity: SeverityMedium},
}

// ScoreDrifts assigns severity levels to a slice of DriftResults.
func ScoreDrifts(results []DriftResult, rules []ScoringRule) []ScoredDrift {
	if rules == nil {
		rules = defaultScoringRules
	}
	scored := make([]ScoredDrift, 0, len(results))
	for _, r := range results {
		scored = append(scored, scoreSingle(r, rules))
	}
	return scored
}

func scoreSingle(r DriftResult, rules []ScoringRule) ScoredDrift {
	if r.Kind == DriftMissing {
		return ScoredDrift{
			Result:   r,
			Severity: SeverityHigh,
			Reason:   "resource missing from cloud",
		}
	}
	if r.Kind == DriftExtra {
		return ScoredDrift{
			Result:   r,
			Severity: SeverityLow,
			Reason:   "extra resource not tracked in state",
		}
	}
	// DriftChanged — inspect attribute names
	for _, attr := range r.ChangedAttributes {
		for _, rule := range rules {
			if matchesPattern(attr, rule.AttributePattern) {
				return ScoredDrift{
					Result:   r,
					Severity: rule.Severity,
					Reason:   fmt.Sprintf("attribute %q matches rule %q", attr, rule.AttributePattern),
				}
			}
		}
	}
	return ScoredDrift{
		Result:   r,
		Severity: SeverityLow,
		Reason:   "non-sensitive attribute changed",
	}
}
