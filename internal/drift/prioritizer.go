package drift

import "sort"

// Priority levels for drift results.
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// PrioritizedResult wraps an AnnotatedResult with a computed priority label.
type PrioritizedResult struct {
	AnnotatedResult
	Priority string
}

// Prioritize assigns a priority level to each AnnotatedResult based on its
// severity and drift kind, then returns results sorted highest-priority first.
func Prioritize(results []AnnotatedResult) []PrioritizedResult {
	prioritized := make([]PrioritizedResult, 0, len(results))
	for _, r := range results {
		prioritized = append(prioritized, PrioritizedResult{
			AnnotatedResult: r,
			Priority:        computePriority(r),
		})
	}
	sort.SliceStable(prioritized, func(i, j int) bool {
		return priorityRank(prioritized[i].Priority) > priorityRank(prioritized[j].Priority)
	})
	return prioritized
}

// FilterByPriority returns only results matching the given priority level.
func FilterByPriority(results []PrioritizedResult, priority string) []PrioritizedResult {
	var out []PrioritizedResult
	for _, r := range results {
		if r.Priority == priority {
			out = append(out, r)
		}
	}
	return out
}

func computePriority(r AnnotatedResult) string {
	if r.Kind == DriftMissing || r.Severity == SeverityCritical {
		return PriorityHigh
	}
	if r.Severity == SeverityHigh {
		return PriorityHigh
	}
	if r.Severity == SeverityMedium || r.Kind == DriftChanged {
		return PriorityMedium
	}
	return PriorityLow
}

func priorityRank(p string) int {
	switch p {
	case PriorityHigh:
		return 3
	case PriorityMedium:
		return 2
	case PriorityLow:
		return 1
	}
	return 0
}
