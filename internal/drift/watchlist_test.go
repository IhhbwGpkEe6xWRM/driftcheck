package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeWatchEntry(key string, attrs ...string) WatchEntry {
	return WatchEntry{
		ResourceKey: key,
		AddedAt:     time.Now().UTC(),
		AlertOnAny:  len(attrs) == 0,
		Attributes:  attrs,
	}
}

func TestWatchlist_AddAndContains(t *testing.T) {
	wl := &Watchlist{}
	wl.Add(makeWatchEntry("aws_instance.web"))
	if !wl.Contains("aws_instance.web") {
		t.Fatal("expected watchlist to contain aws_instance.web")
	}
	if wl.Contains("aws_s3_bucket.data") {
		t.Fatal("expected watchlist not to contain aws_s3_bucket.data")
	}
}

func TestWatchlist_Add_Deduplicates(t *testing.T) {
	wl := &Watchlist{}
	wl.Add(makeWatchEntry("aws_instance.web"))
	wl.Add(makeWatchEntry("aws_instance.web"))
	if len(wl.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(wl.Entries))
	}
}

func TestWatchlist_Remove(t *testing.T) {
	wl := &Watchlist{}
	wl.Add(makeWatchEntry("aws_instance.web"))
	removed := wl.Remove("aws_instance.web")
	if !removed {
		t.Fatal("expected Remove to return true")
	}
	if wl.Contains("aws_instance.web") {
		t.Fatal("expected entry to be removed")
	}
	if wl.Remove("nonexistent") {
		t.Fatal("expected Remove to return false for missing key")
	}
}

func TestSaveAndLoadWatchlist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "watchlist.json")

	wl := &Watchlist{}
	wl.Add(makeWatchEntry("aws_instance.api", "instance_type"))
	wl.Add(makeWatchEntry("aws_s3_bucket.logs"))

	if err := SaveWatchlist(path, wl); err != nil {
		t.Fatalf("SaveWatchlist: %v", err)
	}
	loaded, err := LoadWatchlist(path)
	if err != nil {
		t.Fatalf("LoadWatchlist: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded.Entries))
	}
}

func TestLoadWatchlist_MissingFile(t *testing.T) {
	wl, err := LoadWatchlist("/tmp/nonexistent_watchlist_xyz.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(wl.Entries) != 0 {
		t.Fatal("expected empty watchlist for missing file")
	}
}

func TestLoadWatchlist_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0644)
	_, err := LoadWatchlist(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFilterByWatchlist_AlertOnAny(t *testing.T) {
	wl := &Watchlist{}
	wl.Add(makeWatchEntry("aws_instance.web"))

	results := []DriftResult{
		{ResourceKey: "aws_instance.web", Kind: KindChanged},
		{ResourceKey: "aws_s3_bucket.data", Kind: KindMissing},
	}
	out := FilterByWatchlist(results, wl)
	if len(out) != 1 || out[0].ResourceKey != "aws_instance.web" {
		t.Fatalf("unexpected filter result: %+v", out)
	}
}

func TestFilterByWatchlist_SpecificAttributes(t *testing.T) {
	wl := &Watchlist{}
	wl.Add(makeWatchEntry("aws_instance.web", "instance_type"))

	results := []DriftResult{
		{
			ResourceKey: "aws_instance.web",
			Kind:        KindChanged,
			ChangedAttributes: []ChangedAttribute{
				{Attribute: "instance_type", TFValue: "t2.micro", LiveValue: "t3.micro"},
				{Attribute: "tags", TFValue: "a", LiveValue: "b"},
			},
		},
	}
	out := FilterByWatchlist(results, wl)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if len(out[0].ChangedAttributes) != 1 {
		t.Fatalf("expected 1 changed attribute, got %d", len(out[0].ChangedAttributes))
	}
	if out[0].ChangedAttributes[0].Attribute != "instance_type" {
		t.Fatalf("unexpected attribute: %s", out[0].ChangedAttributes[0].Attribute)
	}
	_ = json.Marshal // satisfy import
}
