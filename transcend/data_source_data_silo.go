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
		DataSilos types.DataSilosPayload `graphql:"dataSilos(filterBy: $filterByInput)"`
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

	titleList := types.WrapValueToList(d.Get("title"))
	if len(titleList) > 0 {
		filters.Title = titleList
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
	if len(query.DataSilos.Nodes) > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding data silos",
			Detail:   "Found multiple data silos matching the given query",
		})
		return diags
	}

	dataSilo := query.DataSilos.Nodes[0]
	d.Set("id", dataSilo.ID)
	d.SetId(string(dataSilo.ID))

	return diags
}
