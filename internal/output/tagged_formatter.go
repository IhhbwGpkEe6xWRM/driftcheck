package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/driftcheck/internal/drift"
)

// WriteTaggedText writes tagged drift results grouped by a tag key as human-readable text.
func WriteTaggedText(w io.Writer, results []drift.TaggedResult, groupByTag string) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "No drift detected.")
		return err
	}

	groups := drift.GroupByTag(results, groupByTag)
	keys := drift.SortedTagKeys(groups)

	for _, key := range keys {
		label := key
		if label == "" {
			label = "(untagged)"
		}
		fmt.Fprintf(w, "[%s: %s]\n", groupByTag, label)
		for _, r := range groups[key] {
			tags := formatTags(r.Tags)
			fmt.Fprintf(w, "  %-40s  kind=%-10s  severity=%-8s  tags=%s\n",
				r.Resource, r.Kind, r.Severity, tags)
			for _, ch := range r.Changes {
				fmt.Fprintf(w, "    attribute=%-30s  want=%v  got=%v\n",
					ch.Attribute, ch.Expected, ch.Actual)
			}
		}
		fmt.Fprintln(w)
	}
	return nil
}

// WriteTaggedJSON writes tagged drift results as JSON, preserving tag metadata.
func WriteTaggedJSON(w io.Writer, results []drift.TaggedResult) error {
	type jsonTag struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	type jsonChange struct {
		Attribute string `json:"attribute"`
		Expected  any    `json:"expected"`
		Actual    any    `json:"actual"`
		Severity  string `json:"severity"`
	}
	type jsonResult struct {
		Resource string       `json:"resource"`
		Kind     string       `json:"kind"`
		Severity string       `json:"severity"`
		Tags     []jsonTag    `json:"tags"`
		Changes  []jsonChange `json:"changes"`
	}

	output := make([]jsonResult, 0, len(results))
	for _, r := range results {
		jr := jsonResult{
			Resource: r.Resource,
			Kind:     r.Kind,
			Severity: string(r.Severity),
			Tags:     make([]jsonTag, 0, len(r.Tags)),
			Changes:  make([]jsonChange, 0, len(r.Changes)),
		}
		for _, t := range r.Tags {
			jr.Tags = append(jr.Tags, jsonTag{Key: t.Key, Value: t.Value})
		}
		for _, ch := range r.Changes {
			jr.Changes = append(jr.Changes, jsonChange{
				Attribute: ch.Attribute,
				Expected:  ch.Expected,
				Actual:    ch.Actual,
				Severity:  string(ch.Severity),
			})
		}
		output = append(output, jr)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

func formatTags(tags []drift.Tag) string {
	if len(tags) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(tags))
	for _, t := range tags {
		parts = append(parts, t.Key+"="+t.Value)
	}
	return strings.Join(parts, ",")
}
