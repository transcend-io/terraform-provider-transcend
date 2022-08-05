package types

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

type Scope struct {
	Name graphql.String `json:"name"`
}

type Resource struct {
	ID graphql.String `json:"id"`
}

type APIKey struct {
	ID        graphql.String `json:"id"`
	Title     graphql.String `json:"title"`
	Scopes    []Scope        `json:"scopes"`
	DataSilos []Resource     `json:"dataSilos"`
}

type APIKeyUpdatableFields struct {
	Scopes    []ScopeName  `json:"scopes,omitempty"`
	DataSilos []graphql.ID `json:"dataSilos,omitempty"`
}

type ApiKeyInput struct {
	Title graphql.String `json:"title"`
	APIKeyUpdatableFields
}

type UpdateApiKeyInput struct {
	ID graphql.String `json:"id"`
	APIKeyUpdatableFields
}

func MakeApiKeyInput(d *schema.ResourceData) ApiKeyInput {
	return ApiKeyInput{
		Title:                 graphql.String(d.Get("title").(string)),
		APIKeyUpdatableFields: MakeAPIKeyUpdatableFields((d)),
	}
}

func MakeUpdateApiKeyInput(d *schema.ResourceData) UpdateApiKeyInput {
	return UpdateApiKeyInput{
		ID:                    graphql.String(d.Get("id").(string)),
		APIKeyUpdatableFields: MakeAPIKeyUpdatableFields((d)),
	}
}

func MakeAPIKeyUpdatableFields(d *schema.ResourceData) APIKeyUpdatableFields {
	return APIKeyUpdatableFields{
		Scopes:    CreateScopeNames(d.Get("scopes").([]interface{})),
		DataSilos: ToIDList(d.Get("data_silos").([]interface{})),
	}
}

func ReadApiKeyIntoState(d *schema.ResourceData, key APIKey) {
	d.Set("title", key.Title)
	d.Set("scopes", FlattenScopes(key.Scopes))
	d.Set("data_silos", FlattenDataSilos(key.DataSilos))
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
		ret[i] = scope.Name
	}
	return ret
}

func FlattenDataSilos(silos []Resource) []interface{} {
	ret := make([]interface{}, len(silos))
	for i, silo := range silos {
		ret[i] = silo.ID
	}
	return ret
}
