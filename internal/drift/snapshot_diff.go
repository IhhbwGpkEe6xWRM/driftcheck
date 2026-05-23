package drift

import (
	"fmt"
	"sort"
	"time"
)

// SnapshotDiff represents the difference between two trend snapshots.
type SnapshotDiff struct {
	From          time.Time
	To            time.Time
	AddedDrifts   []string
	RemovedDrifts []string
	Persisted     []string
	DeltaTotal    int
}

// DiffSnapshots compares two TrendSnapshots and returns a SnapshotDiff
// describing which drifts appeared, disappeared, or persisted between them.
func DiffSnapshots(older, newer TrendSnapshot) SnapshotDiff {
	oldKeys := snapshotResultKeys(older)
	newKeys := snapshotResultKeys(newer)

	oldSet := toSet(oldKeys)
	newSet := toSet(newKeys)

	var added, removed, persisted []string

	for k := range newSet {
		if oldSet[k] {
			persisted = append(persisted, k)
		} else {
			added = append(added, k)
		}
	}

	for k := range oldSet {
		if !newSet[k] {
			removed = append(removed, k)
		}
	}

	sort.Strings(added)
	sort.Strings(removed)
	sort.Strings(persisted)

	return SnapshotDiff{
		From:          older.Timestamp,
		To:            newer.Timestamp,
		AddedDrifts:   added,
		RemovedDrifts: removed,
		Persisted:     persisted,
		DeltaTotal:    len(newer.Results) - len(older.Results),
	}
}

// Summary returns a human-readable one-line summary of the diff.
func (d SnapshotDiff) Summary() string {
	return fmt.Sprintf(
		"from %s to %s: +%d added, -%d resolved, %d persisted (delta: %+d)",
		d.From.Format(time.RFC3339),
		d.To.Format(time.RFC3339),
		len(d.AddedDrifts),
		len(d.RemovedDrifts),
		len(d.Persisted),
		d.DeltaTotal,
	)
}

func snapshotResultKeys(s TrendSnapshot) []string {
	keys := make([]string, 0, len(s.Results))
	for _, r := range s.Results {
		keys = append(keys, driftResultKey(r))
	}
	return keys
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
