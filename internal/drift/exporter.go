package drift

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ExportFormat defines supported export formats.
type ExportFormat string

const (
	ExportCSV  ExportFormat = "csv"
	ExportJSON ExportFormat = "json"
)

// ExportRecord is a flat representation of a drift result for export.
type ExportRecord struct {
	Timestamp    string `json:"timestamp"`
	ResourceKey  string `json:"resource_key"`
	Kind         string `json:"kind"`
	Attribute    string `json:"attribute"`
	Expected     string `json:"expected"`
	Actual       string `json:"actual"`
	Severity     string `json:"severity"`
	Score        int    `json:"score"`
}

// ExportResults writes drift results to w in the given format.
func ExportResults(w io.Writer, results []AnnotatedResult, format ExportFormat) error {
	records := flattenResults(results)
	switch format {
	case ExportCSV:
		return writeCSV(w, records)
	case ExportJSON:
		return writeJSONExport(w, records)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func flattenResults(results []AnnotatedResult) []ExportRecord {
	ts := time.Now().UTC().Format(time.RFC3339)
	var records []ExportRecord
	for _, r := range results {
		if len(r.Result.ChangedAttributes) == 0 {
			records = append(records, ExportRecord{
				Timestamp:   ts,
				ResourceKey: r.Result.ResourceKey,
				Kind:        string(r.Result.Kind),
				Attribute:   "",
				Expected:    "",
				Actual:      "",
				Severity:    string(r.Severity),
				Score:       r.Score,
			})
			continue
		}
		for _, ch := range r.Result.ChangedAttributes {
			records = append(records, ExportRecord{
				Timestamp:   ts,
				ResourceKey: r.Result.ResourceKey,
				Kind:        string(r.Result.Kind),
				Attribute:   ch.Attribute,
				Expected:    fmt.Sprintf("%v", ch.Expected),
				Actual:      fmt.Sprintf("%v", ch.Actual),
				Severity:    string(r.Severity),
				Score:       r.Score,
			})
		}
	}
	return records
}

func writeCSV(w io.Writer, records []ExportRecord) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "resource_key", "kind", "attribute", "expected", "actual", "severity", "score"}); err != nil {
		return err
	}
	for _, r := range records {
		row := []string{r.Timestamp, r.ResourceKey, r.Kind, r.Attribute, r.Expected, r.Actual, r.Severity, fmt.Sprintf("%d", r.Score)}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeJSONExport(w io.Writer, records []ExportRecord) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}
