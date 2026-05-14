package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/yourorg/driftcheck/internal/drift"
)

// Format controls how drift results are rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes drift results to a writer in a given format.
type Formatter struct {
	Format Format
	Out    io.Writer
}

// NewFormatter creates a Formatter with the specified format and writer.
func NewFormatter(format Format, out io.Writer) *Formatter {
	return &Formatter{Format: format, Out: out}
}

// Write renders the drift results according to the configured format.
func (f *Formatter) Write(results []drift.DriftResult) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(results)
	case FormatText:
		return f.writeText(results)
	default:
		return fmt.Errorf("unsupported format: %q", f.Format)
	}
}

func (f *Formatter) writeJSON(results []drift.DriftResult) error {
	enc := json.NewEncoder(f.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func (f *Formatter) writeText(results []drift.DriftResult) error {
	if len(results) == 0 {
		fmt.Fprintln(f.Out, "No drift detected.")
		return nil
	}

	w := tabwriter.NewWriter(f.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "RESOURCE\tKIND\tATTRIBUTE\tEXPECTED\tACTUAL")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%v\n",
			r.ResourceKey,
			r.Kind,
			r.Attribute,
			r.Expected,
			r.Actual,
		)
	}
	return w.Flush()
}
