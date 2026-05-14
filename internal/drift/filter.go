package drift

// FilterDrifts removes DriftResult entries that are suppressed by the given RuleSet.
// If rs is nil, the original slice is returned unchanged.
func FilterDrifts(drifts []DriftResult, rs *RuleSet) []DriftResult {
	if rs == nil || len(drifts) == 0 {
		return drifts
	}

	filtered := make([]DriftResult, 0, len(drifts))
	for _, d := range drifts {
		if !rs.ShouldIgnore(d.ResourceKey, d.Attribute) {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

// FilterByType returns only DriftResult entries whose Kind matches one of the
// provided kinds. An empty kinds list returns all drifts unmodified.
func FilterByType(drifts []DriftResult, kinds ...DriftKind) []DriftResult {
	if len(kinds) == 0 {
		return drifts
	}
	allowed := make(map[DriftKind]struct{}, len(kinds))
	for _, k := range kinds {
		allowed[k] = struct{}{}
	}

	filtered := make([]DriftResult, 0, len(drifts))
	for _, d := range drifts {
		if _, ok := allowed[d.Kind]; ok {
			filtered = append(filtered, d)
		}
	}
	return filtered
}
