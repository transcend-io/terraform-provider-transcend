package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
)

// TODO: add support for scopes

func resourceAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIKeyCreate,
		ReadContext:   resourceAPIKeyRead,
		UpdateContext: resourceAPIKeyUpdate,
		DeleteContext: resourceAPIKeyDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title used to identify the API key",
			},
			"scopes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The names of the scopes to add",
			},
			"data_silos": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The ids of the data silos to assign to",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAPIKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		CreateApiKey struct {
			APIKey APIKey
		} `graphql:"createApiKey(input: {title: $title, dataSilos: $data_silos, scopes: $scopes})"`
	}

	sc := d.Get("scopes").([]interface{})
	scopes := make([]ScopeName, len(sc))
	for i, scope := range sc {
		scopes[i] = ScopeName(scope.(string))
	}

	vars := map[string]interface{}{
		"title":      graphql.String(d.Get("title").(string)),
		"data_silos": toIDList(d.Get("data_silos").([]interface{})),
		"scopes":     scopes,
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error creating API Key " + d.Get("title").(string),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(string(mutation.CreateApiKey.APIKey.ID))
	d.Set("title", mutation.CreateApiKey.APIKey.Title)
	// TODO: Set scopes/data_silos

	return nil
}
func resourceAPIKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var query struct {
		APIKey APIKey `graphql:"apiKey(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.ID(d.Get("id").(string)),
	}

	err := client.graphql.Query(context.Background(), &query, vars)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("title", query.APIKey.Title)

	return nil
}
func resourceAPIKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		UpdateApiKey struct {
			APIKey APIKey
		} `graphql:"updateApiKey(input: {id: $id, title: $title, dataSilos: $data_silos})"`
	}

	vars := map[string]interface{}{
		"id":         graphql.ID(d.Get("id").(string)),
		"title":      graphql.String(d.Get("title").(string)),
		"data_silos": toIDList(d.Get("data_silos").([]interface{})),
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating API Key " + d.Get("title").(string),
			Detail:   err.Error(),
		})
		return diags
	}
	return nil
}
func resourceAPIKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		DeleteApiKey struct {
			Success graphql.Boolean
		} `graphql:"deleteApiKey(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.ID(d.Get("id").(string)),
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error deleting API Key " + d.Get("title").(string),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return nil
}
