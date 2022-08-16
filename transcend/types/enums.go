package types

// Enums
type DbIntegrationQuerySuggestionInput string
type RequestActionObjectResolver string
type DataCategoryType string
type ProcessingPurpose string
type RequestAction string
type DataSiloConnectionState string
type ScopeName string
type EnricherType string
type PluginType string

func ToRequestActionObjectResolverList(origs []interface{}) []RequestActionObjectResolver {
	vals := make([]RequestActionObjectResolver, len(origs))
	for i, orig := range origs {
		vals[i] = RequestActionObjectResolver(orig.(string))
	}

	return vals
}
