package drift

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Summary holds aggregated counts from a drift detection run.
type Summary struct {
	Total   int
	Missing int
	Changed int
	Extra   int
}

// Summarize computes a Summary from a slice of DriftResults.
func Summarize(results []DriftResult) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.DriftType {
		case DriftTypeMissing:
			s.Missing++
		case DriftTypeChanged:
			s.Changed++
		case DriftTypeExtra:
			s.Extra++
		}
	}
	return s
}

// WriteReport writes a human-readable drift report to w.
func WriteReport(w io.Writer, results []DriftResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "✓ No drift detected. Infrastructure matches Terraform state.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TYPE\tRESOURCE\tATTRIBUTE\tSTATE VALUE\tLIVE VALUE")
	fmt.Fprintln(tw, "----\t--------\t---------\t-----------\t----------")

	for _, r := range results {
		resource := r.ResourceType + "." + r.ResourceName
		switch r.DriftType {
		case DriftTypeMissing:
			fmt.Fprintf(tw, "MISSING\t%s\t-\t-\t-\n", resource)
		case DriftTypeExtra:
			fmt.Fprintf(tw, "EXTRA\t%s\t-\t-\t-\n", resource)
		case DriftTypeChanged:
			fmt.Fprintf(tw, "CHANGED\t%s\t%s\t%v\t%v\n",
				resource, r.Attribute, r.StateValue, r.LiveValue)
		}
	}

	if err := tw.Flush(); err != nil {
		return err
	}

	s := Summarize(results)
	_, err := fmt.Fprintf(w, "\nSummary: %d drift(s) found — %d missing, %d changed, %d extra\n",
		s.Total, s.Missing, s.Changed, s.Extra)
	return err
}
