package cloud

import "context"

// ResourceFetcher is the interface that wraps FetchResource.
// Implementations retrieve live attributes for a given resource type and ID
// from the underlying cloud provider.
type ResourceFetcher interface {
	FetchResource(ctx context.Context, resourceType, resourceID string) (ResourceAttributes, error)
}

// FetchAll retrieves live attributes for every resource in the provided map.
// The map key format is "<resource_type>/<resource_name>" and the value is the
// resource ID used to query the cloud provider.
// Errors are collected per-resource; a partial result is returned alongside any
// errors encountered.
func FetchAll(ctx context.Context, fetcher ResourceFetcher, resources map[string]string) (map[string]ResourceAttributes, map[string]error) {
	results := make(map[string]ResourceAttributes, len(resources))
	errs := make(map[string]error)

	for key, id := range resources {
		rType := resourceTypeFromKey(key)
		attrs, err := fetcher.FetchResource(ctx, rType, id)
		if err != nil {
			errs[key] = err
			continue
		}
		results[key] = attrs
	}
	return results, errs
}

// resourceTypeFromKey extracts the resource type from a "type/name" key.
// If the key has no slash the entire string is returned as the type.
func resourceTypeFromKey(key string) string {
	for i, ch := range key {
		if ch == '/' {
			return key[:i]
		}
	}
	return key
}
