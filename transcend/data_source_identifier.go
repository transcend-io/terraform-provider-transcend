package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/transcend-io/terraform-provider-transcend/transcend/types"
)

func dataSourceIdentifier() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentifierRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"text": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The text to lookup an identifier by",
			},
		},
	}
}

func dataSourceIdentifierRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var query struct {
		Identifiers struct {
			Nodes []types.Identifier
		} `graphql:"identifiers(filterBy: { text: $text })"`
	}

	vars := map[string]interface{}{
		"text": graphql.String(d.Get("text").(string)),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("Identifier"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding identifier with text " + d.Get("text").(string),
			Detail:   "Error when finding identifier: " + err.Error(),
		})
		return diags
	}
	if len(query.Identifiers.Nodes) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding identifier with text " + d.Get("text").(string),
			Detail:   "Found 0 identifiers for given text",
		})
		return diags
	}
	if len(query.Identifiers.Nodes) > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding identifier with text " + d.Get("text").(string),
			Detail:   "Found multiple identifiers for given text",
		})
		return diags
	}

	identifier := query.Identifiers.Nodes[0]
	d.Set("id", identifier.ID)
	d.Set("name", identifier.Name)
	d.SetId(string(identifier.ID))

	return diags
}
