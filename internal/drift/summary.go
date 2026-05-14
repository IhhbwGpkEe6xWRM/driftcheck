package drift

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// DriftSummary holds aggregated statistics about detected drift.
type DriftSummary struct {
	Total    int
	Missing  int
	Changed  int
	Extra    int
	ByType   map[string]int
}

// Summarize computes a DriftSummary from a slice of DriftResult.
func Summarize(results []DriftResult) DriftSummary {
	s := DriftSummary{
		ByType: make(map[string]int),
	}
	for _, r := range results {
		s.Total++
		switch r.Kind {
		case KindMissing:
			s.Missing++
		case KindChanged:
			s.Changed++
		case KindExtra:
			s.Extra++
		}
		rtype, _, _ := splitKey(r.ResourceKey)
		s.ByType[rtype]++
	}
	return s
}

// PrintSummary writes a human-readable summary table to w.
func PrintSummary(w io.Writer, s DriftSummary) {
	fmt.Fprintf(w, "Drift Summary\n")
	fmt.Fprintf(w, "  Total:   %d\n", s.Total)
	fmt.Fprintf(w, "  Missing: %d\n", s.Missing)
	fmt.Fprintf(w, "  Changed: %d\n", s.Changed)
	fmt.Fprintf(w, "  Extra:   %d\n", s.Extra)

	if len(s.ByType) == 0 {
		return
	}

	fmt.Fprintf(w, "\nBy Resource Type:\n")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	types := make([]string, 0, len(s.ByType))
	for t := range s.ByType {
		types = append(types, t)
	}
	sort.Strings(types)

	for _, t := range types {
		fmt.Fprintf(tw, "  %s\t%d\n", t, s.ByType[t])
	}
	tw.Flush()
}
