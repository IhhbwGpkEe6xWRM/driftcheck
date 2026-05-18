package drift

import "fmt"

// DeduplicateResults removes duplicate DriftResult entries from a slice,
// keeping the first occurrence of each unique (resource key, attribute) pair.
func DeduplicateResults(results []DriftResult) []DriftResult {
	seen := make(map[string]struct{})
	out := make([]DriftResult, 0, len(results))

	for _, r := range results {
		key := dedupeKey(r)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, r)
	}

	return out
}

// DeduplicateByResource collapses multiple drift results for the same resource
// into a single entry, merging changed attributes.
func DeduplicateByResource(results []DriftResult) []DriftResult {
	type entry struct {
		index  int
		result DriftResult
	}

	index := make(map[string]*entry)
	order := []string{}

	for _, r := range results {
		rk := fmt.Sprintf("%s::%s", r.ResourceType, r.ResourceName)
		if e, exists := index[rk]; exists {
			for k, v := range r.ChangedAttributes {
				if e.result.ChangedAttributes == nil {
					e.result.ChangedAttributes = make(map[string]AttributeDiff)
				}
				e.result.ChangedAttributes[k] = v
			}
			index[rk].result = e.result
		} else {
			copy := r
			if copy.ChangedAttributes != nil {
				merged := make(map[string]AttributeDiff, len(r.ChangedAttributes))
				for k, v := range r.ChangedAttributes {
					merged[k] = v
				}
				copy.ChangedAttributes = merged
			}
			index[rk] = &entry{index: len(order), result: copy}
			order = append(order, rk)
		}
	}

	out := make([]DriftResult, 0, len(order))
	for _, rk := range order {
		out = append(out, index[rk].result)
	}
	return out
}

func dedupeKey(r DriftResult) string {
	attrs := ""
	for k := range r.ChangedAttributes {
		attrs += k + ","
	}
	return fmt.Sprintf("%s::%s::%s::%s", r.ResourceType, r.ResourceName, r.Kind, attrs)
}
