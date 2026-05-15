package drift

import "sort"

// GroupedDrift holds drift results organized by resource type.
type GroupedDrift struct {
	ByType map[string][]DriftResult
	Ordered []string // sorted type keys
}

// GroupByType organizes a slice of DriftResults into a GroupedDrift,
// keyed by the resource type extracted from each result's ResourceKey.
func GroupByType(results []DriftResult) GroupedDrift {
	byType := make(map[string][]DriftResult)

	for _, r := range results {
		rType, _ := splitKey(r.ResourceKey)
		if rType == "" {
			rType = "unknown"
		}
		byType[rType] = append(byType[rType], r)
	}

	ordered := make([]string, 0, len(byType))
	for k := range byType {
		ordered = append(ordered, k)
	}
	sort.Strings(ordered)

	return GroupedDrift{
		ByType:  byType,
		Ordered: ordered,
	}
}

// GroupByKind organizes drift results by DriftKind (changed, missing, extra).
func GroupByKind(results []DriftResult) map[DriftKind][]DriftResult {
	byKind := make(map[DriftKind][]DriftResult)
	for _, r := range results {
		byKind[r.Kind] = append(byKind[r.Kind], r)
	}
	return byKind
}

// TopDriftedTypes returns up to n resource types with the most drift entries,
// ordered descending by count.
func TopDriftedTypes(g GroupedDrift, n int) []string {
	type entry struct {
		name  string
		count int
	}

	entries := make([]entry, 0, len(g.ByType))
	for _, name := range g.Ordered {
		entries = append(entries, entry{name: name, count: len(g.ByType[name])})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count != entries[j].count {
			return entries[i].count > entries[j].count
		}
		return entries[i].name < entries[j].name
	})

	if n > len(entries) {
		n = len(entries)
	}

	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = entries[i].name
	}
	return result
}
