package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/transcend-io/terraform-provider-transcend/transcend/types"
)

func dataSourceDataSilo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataSiloRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"discoveredby": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the parent data silo that discovered this data silo",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the data silo",
			},
			"title": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The title of the data silo",
			},
		},
	}
}

func dataSourceDataSiloRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var query struct {
		DataSilos []types.DataSilo `graphql:"dataSilos(filterBy: $filterByInput)"`
	}

	vars := map[string]interface{}{
		"filterByInput": types.DataSiloFilter{
			DiscoveredBy: types.WrapValueToList(d.Get("discoveredby")),
			Type:         types.WrapValueToList(d.Get("type")),
			Title:        types.WrapValueToList(d.Get("title")),
		},
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
	if len(query.DataSilos) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding data silo",
			Detail:   "Found 0 data silos for given query",
		})
		return diags
	}
	if len(query.DataSilos) > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding data silos",
			Detail:   "Found multiple data silos matching the given query",
		})
		return diags
	}

	dataSilo := query.DataSilos[0]
	d.Set("id", dataSilo.ID)
	d.SetId(string(dataSilo.ID))

	return diags
}
