package tfstate

import (
	"encoding/json"
	"fmt"
	"os"
)

// Resource represents a single resource entry from Terraform state.
type Resource struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Provider   string                 `json:"provider"`
	Attributes map[string]interface{} `json:"attributes"`
}

// State represents a parsed Terraform state file.
type State struct {
	Version   int        `json:"version"`
	Resources []Resource `json:"resources"`
}

// tfStateFile mirrors the top-level structure of a .tfstate JSON file.
type tfStateFile struct {
	Version   int              `json:"version"`
	Resources []tfStateResource `json:"resources"`
}

type tfStateResource struct {
	Type      string             `json:"type"`
	Name      string             `json:"name"`
	Provider  string             `json:"provider"`
	Instances []tfStateInstance  `json:"instances"`
}

type tfStateInstance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// ParseFile reads and parses a Terraform state file from the given path.
func ParseFile(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading state file %q: %w", path, err)
	}
	return Parse(data)
}

// Parse parses raw Terraform state JSON bytes into a State struct.
func Parse(data []byte) (*State, error) {
	var raw tfStateFile
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("unmarshalling state JSON: %w", err)
	}

	if raw.Version < 4 {
		return nil, fmt.Errorf("unsupported state version %d (minimum supported: 4)", raw.Version)
	}

	state := &State{
		Version:   raw.Version,
		Resources: make([]Resource, 0, len(raw.Resources)),
	}

	for _, r := range raw.Resources {
		if len(r.Instances) == 0 {
			continue
		}
		// Use the first instance's attributes; multi-count resources can be
		// extended in a future iteration.
		state.Resources = append(state.Resources, Resource{
			Type:       r.Type,
			Name:       r.Name,
			Provider:   r.Provider,
			Attributes: r.Instances[0].Attributes,
		})
	}

	return state, nil
}
