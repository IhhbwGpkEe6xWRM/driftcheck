package drift

import (
	"fmt"

	"github.com/driftcheck/internal/tfstate"
)

// DriftType describes the kind of drift detected.
type DriftType string

const (
	DriftTypeMissing  DriftType = "missing"   // resource exists in state but not in cloud
	DriftTypeChanged  DriftType = "changed"   // resource exists but attributes differ
	DriftTypeExtra    DriftType = "extra"     // resource exists in cloud but not in state
)

// DriftResult holds a single detected drift item.
type DriftResult struct {
	ResourceType string
	ResourceName string
	DriftType    DriftType
	Attribute    string
	StateValue   interface{}
	LiveValue    interface{}
}

func (d DriftResult) String() string {
	switch d.DriftType {
	case DriftTypeMissing:
		return fmt.Sprintf("[MISSING] %s.%s not found in live infrastructure", d.ResourceType, d.ResourceName)
	case DriftTypeExtra:
		return fmt.Sprintf("[EXTRA] %s.%s exists in cloud but not in Terraform state", d.ResourceType, d.ResourceName)
	case DriftTypeChanged:
		return fmt.Sprintf("[CHANGED] %s.%s .%s: state=%v live=%v", d.ResourceType, d.ResourceName, d.Attribute, d.StateValue, d.LiveValue)
	}
	return ""
}

// LiveResourceMap maps "resourceType.resourceName" to its live attribute map.
type LiveResourceMap map[string]map[string]interface{}

// Detect compares Terraform state resources against live cloud resources and
// returns all detected drift items.
func Detect(resources []tfstate.Resource, live LiveResourceMap) []DriftResult {
	var results []DriftResult

	seen := make(map[string]bool)

	for _, res := range resources {
		for _, inst := range res.Instances {
			key := resourceKey(res.Type, res.Name)
			seen[key] = true

			liveAttrs, exists := live[key]
			if !exists {
				results = append(results, DriftResult{
					ResourceType: res.Type,
					ResourceName: res.Name,
					DriftType:    DriftTypeMissing,
				})
				continue
			}

			for attr, stateVal := range inst.Attributes {
				liveVal, ok := liveAttrs[attr]
				if !ok || fmt.Sprintf("%v", liveVal) != fmt.Sprintf("%v", stateVal) {
					results = append(results, DriftResult{
						ResourceType: res.Type,
						ResourceName: res.Name,
						DriftType:    DriftTypeChanged,
						Attribute:    attr,
						StateValue:   stateVal,
						LiveValue:    liveVal,
					})
				}
			}
		}
	}

	for key := range live {
		if !seen[key] {
			rType, rName := splitKey(key)
			results = append(results, DriftResult{
				ResourceType: rType,
				ResourceName: rName,
				DriftType:    DriftTypeExtra,
			})
		}
	}

	return results
}

func resourceKey(rType, rName string) string {
	return rType + "." + rName
}

func splitKey(key string) (string, string) {
	for i, c := range key {
		if c == '.' {
			return key[:i], key[i+1:]
		}
	}
	return key, ""
}
