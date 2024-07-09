package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/transcend-io/terraform-provider-transcend/transcend/types"
)

func dataSourceDataSilos() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataSilosRead,
		Schema: map[string]*schema.Schema{
			"ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The attribute values used to label resources",
			},
			"discoveredby": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the parent data silo that discovered these data silos",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the data silos",
			},
		},
	}
}

func dataSourceDataSilosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var query struct {
		DataSilos types.DataSilosPayload `graphql:"dataSilos(filterBy: $filterByInput, first: 100)"`
	}

	filters := types.DataSiloFiltersInput{}
	discoveredByList := types.WrapValueToIDList(d.Get("discoveredby"))
	if len(discoveredByList) > 0 {
		filters.DiscoveredBy = discoveredByList
	}

	typeList := types.WrapValueToList(d.Get("type"))
	if len(typeList) > 0 {
		filters.Type = typeList
	}

	vars := map[string]interface{}{
		"filterByInput": filters,
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("DataSilos"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding data silo",
			Detail:   "Error when finding data silo: " + err.Error(),
		})
		return diags
	}
	if len(query.DataSilos.Nodes) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding data silo",
			Detail:   "Found 0 data silos for given query",
		})
		return diags
	}

	ids := make([]string, len(query.DataSilos.Nodes))
	for i, dataSilo := range query.DataSilos.Nodes {
		ids[i] = string(dataSilo.ID)
	}

	// Make the ID a csv of the found IDs
	id := ""
	for _, i := range ids {
		id += i + ","
	}

	d.SetId(id)
	d.Set("ids", ids)

	return diags
}
