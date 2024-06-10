package types

import graphql "github.com/hasura/go-graphql-client"

func ToStringList(raw interface{}) []graphql.String {
	if raw == nil {
		return []graphql.String{}
	}
	origs := raw.([]interface{})
	vals := make([]graphql.String, len(origs))
	for i, orig := range origs {
		vals[i] = graphql.String(orig.(string))
	}
	return vals
}

func FromStringList(strings []graphql.String) []interface{} {
	vals := make([]interface{}, len(strings))
	for i, graphQlString := range strings {
		vals[i] = graphQlString
	}
	return vals
}

func ToString(raw interface{}) graphql.String {
	if raw == nil {
		return ""
	}
	return graphql.String(raw.(string))
}

func ToIDList(origs []interface{}) []graphql.ID {
	vals := make([]graphql.ID, len(origs))
	for i, orig := range origs {
		vals[i] = graphql.ID(orig.(string))
	}

	return vals
}

func WrapValueToList(orig interface{}) []graphql.String {
	if orig == nil {
		return []graphql.String{}
	}
	vals := make([]graphql.String, 1)
	vals[0] = graphql.String(orig.(string))
	return vals

}

func ToRequestActionList(origs []interface{}) []RequestAction {
	vals := make([]RequestAction, len(origs))
	for i, orig := range origs {
		vals[i] = RequestAction(orig.(string))
	}

	return vals
}
