package drift

import (
	"time"
)

// Suppression represents a rule that silences drift results for a given
// resource key until an optional expiry time.
type Suppression struct {
	ResourceKey string    `json:"resource_key"`
	Reason      string    `json:"reason"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Suppressor holds a set of suppressions and filters drift results.
type Suppressor struct {
	suppressions []Suppression
	now          func() time.Time
}

// NewSuppressor creates a Suppressor from the provided suppression list.
func NewSuppressor(suppressions []Suppression) *Suppressor {
	return &Suppressor{
		suppressions: suppressions,
		now:          time.Now,
	}
}

// IsSuppressed reports whether the given drift result should be suppressed.
func (s *Suppressor) IsSuppressed(result Result) bool {
	for _, sup := range s.suppressions {
		if sup.ResourceKey != result.ResourceKey {
			continue
		}
		if sup.ExpiresAt.IsZero() || s.now().Before(sup.ExpiresAt) {
			return true
		}
	}
	return false
}

// Apply filters out suppressed results and returns the remaining ones along
// with a count of how many were suppressed.
func (s *Suppressor) Apply(results []Result) (kept []Result, suppressedCount int) {
	for _, r := range results {
		if s.IsSuppressed(r) {
			suppressedCount++
		} else {
			kept = append(kept, r)
		}
	}
	return kept, suppressedCount
}
