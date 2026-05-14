package output

import (
	"fmt"
	"strings"
)

// FormatFlag implements flag.Value so Format can be used as a CLI flag.
type FormatFlag struct {
	Value Format
}

// String returns the current format value as a string.
func (f *FormatFlag) String() string {
	if f.Value == "" {
		return string(FormatText)
	}
	return string(f.Value)
}

// Set parses and validates a format string from the command line.
func (f *FormatFlag) Set(s string) error {
	norm := Format(strings.ToLower(strings.TrimSpace(s)))
	switch norm {
	case FormatText, FormatJSON:
		f.Value = norm
		return nil
	default:
		return fmt.Errorf("invalid format %q: must be one of [text, json]", s)
	}
}

// Type returns the flag type name, used by pflag-compatible libraries.
func (f *FormatFlag) Type() string {
	return "format"
}

// DefaultFormat returns a FormatFlag initialised to the text format.
func DefaultFormat() FormatFlag {
	return FormatFlag{Value: FormatText}
}
