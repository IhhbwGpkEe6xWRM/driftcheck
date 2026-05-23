package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/driftcheck/internal/drift"
)

func makeWatchlist(entries ...drift.WatchEntry) *drift.Watchlist {
	wl := &drift.Watchlist{}
	for _, e := range entries {
		wl.Add(e)
	}
	return wl
}

func TestWriteWatchlistText_Empty(t *testing.T) {
	var buf bytes.Buffer
	wl := &drift.Watchlist{}
	if err := WriteWatchlistText(&buf, wl); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "empty") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteWatchlistText_WithEntries(t *testing.T) {
	wl := makeWatchlist(
		drift.WatchEntry{
			ResourceKey: "aws_instance.web",
			AddedAt:     time.Now().UTC(),
			AlertOnAny:  true,
			Note:        "critical host",
		},
		drift.WatchEntry{
			ResourceKey: "aws_s3_bucket.data",
			AddedAt:     time.Now().UTC(),
			AlertOnAny:  false,
			Attributes:  []string{"versioning", "acl"},
		},
	)
	var buf bytes.Buffer
	if err := WriteWatchlistText(&buf, wl); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_instance.web") {
		t.Errorf("expected aws_instance.web in output")
	}
	if !strings.Contains(out, "critical host") {
		t.Errorf("expected note 'critical host' in output")
	}
	if !strings.Contains(out, "versioning,acl") {
		t.Errorf("expected attributes 'versioning,acl' in output")
	}
	if !strings.Contains(out, "RESOURCE KEY") {
		t.Errorf("expected header row in output")
	}
}

func TestWriteWatchlistJSON_ValidOutput(t *testing.T) {
	wl := makeWatchlist(
		drift.WatchEntry{
			ResourceKey: "aws_instance.api",
			AddedAt:     time.Now().UTC(),
			AlertOnAny:  true,
		},
	)
	var buf bytes.Buffer
	if err := WriteWatchlistJSON(&buf, wl); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed drift.Watchlist
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(parsed.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(parsed.Entries))
	}
	if parsed.Entries[0].ResourceKey != "aws_instance.api" {
		t.Errorf("unexpected resource key: %s", parsed.Entries[0].ResourceKey)
	}
}

func TestJoinStrings(t *testing.T) {
	if joinStrings(nil) != "" {
		t.Error("expected empty string for nil slice")
	}
	if joinStrings([]string{"a"}) != "a" {
		t.Error("expected 'a'")
	}
	if joinStrings([]string{"a", "b", "c"}) != "a,b,c" {
		t.Error("expected 'a,b,c'")
	}
}
