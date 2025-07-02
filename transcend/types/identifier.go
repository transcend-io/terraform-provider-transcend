package types

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

type Identifier struct {
	   ID   graphql.String `json:"id"`
	   Name graphql.String `json:"name"`
	   Type graphql.String `json:"type"`
}

type IdentifierInput struct {
	   Name graphql.String `json:"name"`
	   Type graphql.String `json:"type"`
}

type UpdateIdentifierInput struct {
	   ID   graphql.String `json:"id"`
	   Name graphql.String `json:"name"`
	   Type graphql.String `json:"type"`
}

func MakeIdentifierInput(d *schema.ResourceData) IdentifierInput {
	return IdentifierInput{
		Name: graphql.String(d.Get("name").(string)),
		Type: graphql.String(d.Get("type").(string)),
	}
}

func MakeUpdateIdentifierInput(d *schema.ResourceData) UpdateIdentifierInput {
	return UpdateIdentifierInput{
		ID:   graphql.String(d.Get("id").(string)),
		Name: graphql.String(d.Get("name").(string)),
		Type: graphql.String(d.Get("type").(string)),
	}
}

func ReadIdentifierIntoState(d *schema.ResourceData, identifier Identifier) {
	d.Set("name", identifier.Name)
	d.Set("type", identifier.Type)
}
