package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// WatchEntry represents a single resource being actively watched for drift.
type WatchEntry struct {
	ResourceKey string    `json:"resource_key"`
	AddedAt     time.Time `json:"added_at"`
	Note        string    `json:"note,omitempty"`
	AlertOnAny  bool      `json:"alert_on_any"`
	Attributes  []string  `json:"attributes,omitempty"`
}

// Watchlist holds a set of resources to monitor closely.
type Watchlist struct {
	Entries []WatchEntry `json:"entries"`
}

// Add inserts a new entry into the watchlist, deduplicating by resource key.
func (w *Watchlist) Add(entry WatchEntry) {
	for _, e := range w.Entries {
		if e.ResourceKey == entry.ResourceKey {
			return
		}
	}
	if entry.AddedAt.IsZero() {
		entry.AddedAt = time.Now().UTC()
	}
	w.Entries = append(w.Entries, entry)
}

// Remove deletes an entry by resource key.
func (w *Watchlist) Remove(resourceKey string) bool {
	for i, e := range w.Entries {
		if e.ResourceKey == resourceKey {
			w.Entries = append(w.Entries[:i], w.Entries[i+1:]...)
			return true
		}
	}
	return false
}

// Contains reports whether the given resource key is in the watchlist.
func (w *Watchlist) Contains(resourceKey string) bool {
	for _, e := range w.Entries {
		if e.ResourceKey == resourceKey {
			return true
		}
	}
	return false
}

// SaveWatchlist persists the watchlist to a JSON file.
func SaveWatchlist(path string, wl *Watchlist) error {
	data, err := json.MarshalIndent(wl, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal watchlist: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadWatchlist reads a watchlist from a JSON file.
func LoadWatchlist(path string) (*Watchlist, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Watchlist{}, nil
		}
		return nil, fmt.Errorf("read watchlist: %w", err)
	}
	var wl Watchlist
	if err := json.Unmarshal(data, &wl); err != nil {
		return nil, fmt.Errorf("parse watchlist: %w", err)
	}
	return &wl, nil
}

// FilterByWatchlist returns only the DriftResults whose resource key appears
// in the watchlist. If the entry specifies attributes, only those attribute
// changes are retained.
func FilterByWatchlist(results []DriftResult, wl *Watchlist) []DriftResult {
	var out []DriftResult
	for _, r := range results {
		for _, e := range wl.Entries {
			if e.ResourceKey != r.ResourceKey {
				continue
			}
			if e.AlertOnAny || len(e.Attributes) == 0 {
				out = append(out, r)
				break
			}
			// Filter to only watched attributes.
			filtered := filterAttrs(r, e.Attributes)
			if len(filtered.ChangedAttributes) > 0 || filtered.Kind != KindChanged {
				out = append(out, filtered)
			}
			break
		}
	}
	return out
}

func filterAttrs(r DriftResult, attrs []string) DriftResult {
	set := make(map[string]struct{}, len(attrs))
	for _, a := range attrs {
		set[a] = struct{}{}
	}
	copy := r
	copy.ChangedAttributes = nil
	for _, ca := range r.ChangedAttributes {
		if _, ok := set[ca.Attribute]; ok {
			copy.ChangedAttributes = append(copy.ChangedAttributes, ca)
		}
	}
	return copy
}
