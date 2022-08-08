package types

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

type ReconnectDataSiloInput struct {
	DataSiloId       graphql.ID              `json:"dataSiloId"`
	PlaintextContext []PlaintextContextInput `json:"plaintextContext,omitempty"`
}

func CreateReconnectDataSiloFields(d *schema.ResourceData) ReconnectDataSiloInput {
	return ReconnectDataSiloInput{
		DataSiloId:       graphql.String(d.Get("id").(string)),
		PlaintextContext: ToPlaintextContextList(d.Get("plaintext_context").([]interface{})),
	}
}

func ReadDataSiloConnectionIntoState(d *schema.ResourceData, silo DataSilo) {
	d.Set("id", silo.ID)
	d.Set("data_silo_id", silo.ID)
	d.Set("connectionState", silo.ConnectionState)
	d.Set("plaintext_context", FromPlaintextContextList(silo.PlaintextContext))
}

func ToPlaintextContextList(plaintextContexts []interface{}) []PlaintextContextInput {
	vals := make([]PlaintextContextInput, len(plaintextContexts))
	for i, rawContext := range plaintextContexts {
		context := rawContext.(map[string]interface{})
		vals[i] = PlaintextContextInput{
			Name:  graphql.String(context["name"].(string)),
			Value: graphql.String(context["value"].(string)),
		}
	}
	return vals
}

func FromPlaintextContextList(plaintextContexts []PlaintextContextInput) []map[string]interface{} {
	vals := make([]map[string]interface{}, len(plaintextContexts))
	for i, context := range plaintextContexts {
		vals[i] = map[string]interface{}{
			"name":  context.Name,
			"value": context.Value,
		}
	}
	return vals
}
