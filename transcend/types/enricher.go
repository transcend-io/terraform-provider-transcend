package types

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

type IDObject struct {
	ID graphql.String `json:"id"`
}

type Enricher struct {
	ID              graphql.String  `json:"id"`
	Title           graphql.String  `json:"title"`
	Description     graphql.String  `json:"description"`
	URL             graphql.String  `json:"url"`
	InputIdentifier IDObject        `json:"inputIdentifier"`
	Identifiers     []IDObject      `json:"identifiers"`
	Headers         []Header        `json:"headers"`
	Actions         []RequestAction `json:"actions"`
	Type            EnricherType    `json:"type"`
}

type UpdateEnricherInput struct {
	ID graphql.String `json:"id"`
	EnricherUpdatableFields
}

type EnricherInput struct {
	EnricherUpdatableFields
}

type EnricherUpdatableFields struct {
	Title           graphql.String      `json:"title"`
	Description     graphql.String      `json:"description,omitempty"`
	URL             graphql.String      `json:"url,omitempty"`
	InputIdentifier graphql.String      `json:"inputIdentifier"`
	Identifiers     []graphql.String    `json:"identifiers"`
	Headers         []CustomHeaderInput `json:"headers,omitempty"`
	Actions         []RequestAction     `json:"actions,omitempty"`
	Type            EnricherType        `json:"type"`
	PhoneNumbers    []graphql.String    `json:"phoneNumbers"`

	// TODO: Add more fields
	// DataSiloId
	// signedIdentifierInputs
	// userId
}

func MakeUpdateEnricherInput(d *schema.ResourceData) UpdateEnricherInput {
	return UpdateEnricherInput{
		ID:                      graphql.String(d.Get("id").(string)),
		EnricherUpdatableFields: MakeEnricherUpdatableFields(d),
	}
}

func MakeEnricherInput(d *schema.ResourceData) EnricherInput {
	return EnricherInput{
		EnricherUpdatableFields: MakeEnricherUpdatableFields(d),
	}
}

func MakeEnricherUpdatableFields(d *schema.ResourceData) EnricherUpdatableFields {
	return EnricherUpdatableFields{
		Title:           graphql.String(d.Get("title").(string)),
		Description:     graphql.String(d.Get("description").(string)),
		URL:             graphql.String(d.Get("url").(string)),
		Headers:         ToCustomHeaderInputList(d.Get("headers").([]interface{})),
		Actions:         ToRequestActionList(d.Get("actions").([]interface{})),
		Identifiers:     ToStringList(d.Get("output_identifiers").([]interface{})),
		InputIdentifier: graphql.String(d.Get("input_identifier").(string)),
		Type:            EnricherType(d.Get("type").(string)),
		// This is not a fully supported type, but it must be present on the creation mutation
		PhoneNumbers: []graphql.String{},
	}
}

func ReadEnricherIntoState(d *schema.ResourceData, enricher Enricher) {
	d.Set("title", enricher.Title)
	d.Set("type", enricher.Type)
	d.Set("description", enricher.Description)
	d.Set("url", enricher.URL)
	d.Set("input_identifier", enricher.InputIdentifier.ID)
	d.Set("output_identifiers", FlattenIDObject(enricher.Identifiers))
	d.Set("actions", enricher.Actions)
	d.Set("headers", FlattenHeaders(&enricher.Headers))
}

func FlattenRequestAction(actions []RequestAction) []interface{} {
	ret := make([]interface{}, len(actions))
	for i, action := range actions {
		ret[i] = string(action)
	}
	return ret
}

func FlattenIDObject(objects []IDObject) []interface{} {
	ret := make([]interface{}, len(objects))
	for i, object := range objects {
		ret[i] = object.ID
	}
	return ret
}
