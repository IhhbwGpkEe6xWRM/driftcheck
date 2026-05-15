package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/driftcheck/internal/drift"
)

// WriteAnnotatedText writes annotated drift results in a human-readable format.
func WriteAnnotatedText(w io.Writer, items []drift.AnnotatedDrift) error {
	if len(items) == 0 {
		_, err := fmt.Fprintln(w, "No drift detected.")
		return err
	}
	for _, item := range items {
		sev := strings.ToUpper(string(item.Annotation.Severity))
		fmt.Fprintf(w, "[%s] %s\n", sev, item.Result.ResourceKey)
		fmt.Fprintf(w, "  Kind:    %s\n", item.Result.Kind)
		fmt.Fprintf(w, "  Message: %s\n", item.Annotation.Message)
		if item.Annotation.Hint != "" {
			fmt.Fprintf(w, "  Hint:    %s\n", item.Annotation.Hint)
		}
		for _, attr := range item.Result.Attributes {
			fmt.Fprintf(w, "    ~ %s: %q => %q\n", attr.Key, attr.StateVal, attr.CloudVal)
		}
		fmt.Fprintln(w)
	}
	return nil
}

// WriteAnnotatedJSON writes annotated drift results as a JSON array.
func WriteAnnotatedJSON(w io.Writer, items []drift.AnnotatedDrift) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(items)
}
