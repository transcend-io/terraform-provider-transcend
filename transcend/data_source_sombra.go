package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/transcend-io/terraform-provider-transcend/transcend/types"
)

func dataSourceSombra() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSombraRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The url to lookup a sombra instance by",
			},
		},
	}
}

func dataSourceSombraRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var query struct {
		Sombras []types.Sombra `graphql:"sombras(filterBy: { urls: [$url] })"`
	}

	vars := map[string]interface{}{
		"url": graphql.String(d.Get("url").(string)),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("Sombras"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding sombra with url " + d.Get("url").(string),
			Detail:   "Error when finding sombra: " + err.Error(),
		})
		return diags
	}
	if len(query.Sombras) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding sombra with text " + d.Get("url").(string),
			Detail:   "Found 0 sombras for given text",
		})
		return diags
	}
	if len(query.Sombras) > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding sombra with text " + d.Get("url").(string),
			Detail:   "Found multiple sombras for given url",
		})
		return diags
	}

	sombra := query.Sombras[0]
	d.Set("id", sombra.ID)
	d.SetId(string(sombra.ID))

	return diags
}
