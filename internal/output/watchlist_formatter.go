package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/user/driftcheck/internal/drift"
)

// WriteWatchlistText writes the watchlist entries as a human-readable table.
func WriteWatchlistText(w io.Writer, wl *drift.Watchlist) error {
	if len(wl.Entries) == 0 {
		_, err := fmt.Fprintln(w, "Watchlist is empty.")
		return err
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE KEY\tADDED AT\tALERT ON ANY\tATTRIBUTES\tNOTE")
	fmt.Fprintln(tw, "------------\t--------\t------------\t----------\t----")
	for _, e := range wl.Entries {
		attrs := "*"
		if len(e.Attributes) > 0 {
			attrs = joinStrings(e.Attributes)
		}
		alertAny := "no"
		if e.AlertOnAny {
			alertAny = "yes"
		}
		note := e.Note
		if note == "" {
			note = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			e.ResourceKey,
			e.AddedAt.Format(time.RFC3339),
			alertAny,
			attrs,
			note,
		)
	}
	return tw.Flush()
}

// WriteWatchlistJSON writes the watchlist as a JSON document.
func WriteWatchlistJSON(w io.Writer, wl *drift.Watchlist) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(wl)
}

func joinStrings(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	out := ss[0]
	for _, s := range ss[1:] {
		out += "," + s
	}
	return out
}
