package types

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
)

type Scope struct {
	Type graphql.String `json:"type"`
}

type APIKey struct {
	ID    graphql.String `json:"id"`
	Title graphql.String `json:"title"`
	// Scopes    []Scope        `json:"scopes"`
	// DataSilos []DataSilo     `json:"dataSilos"`
}

type ApiKeyInput struct {
	Title graphql.String `json:"title"`
	// Scopes    []ScopeName    `json:"scopes"`
	// DataSilos []graphql.ID   `json:"dataSilos"`
}

type UpdateApiKeyInput ApiKeyInput

func MakeApiKeyInput(d *schema.ResourceData) ApiKeyInput {
	return ApiKeyInput{
		Title: graphql.String(d.Get("title").(string)),
		// DataSilos: ToIDList(d.Get("data_silos").([]interface{})),
		// Scopes:    CreateScopeNames(d.Get("scopes").([]interface{})),
	}
}

func ReadApiKeyIntoState(d *schema.ResourceData, key APIKey) {
	d.Set("title", key.Title)
	// d.Set("scopes", types.FlattenScopes(key.Scopes))
	// d.Set("data_silos", types.FlattenDataSilos(key.DataSilos))
}

func CreateScopeNames(rawScopes []interface{}) []ScopeName {
	scopes := make([]ScopeName, len(rawScopes))
	for i, scope := range rawScopes {
		scopes[i] = ScopeName(scope.(string))
	}
	return scopes
}

func FlattenScopes(scopes []Scope) []interface{} {
	ret := make([]interface{}, len(scopes))
	for i, scope := range scopes {
		ret[i] = scope.Type
	}
	return ret
}

func FlattenDataSilos(silos []DataSilo) []interface{} {
	ret := make([]interface{}, len(silos))
	for i, silo := range silos {
		ret[i] = silo.ID
	}
	return ret
}
