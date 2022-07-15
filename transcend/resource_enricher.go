package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
)

func resourceEnricher() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataSilosCreate,
		ReadContext:   resourceDataSilosRead,
		UpdateContext: resourceDataSilosUpdate,
		DeleteContext: resourceDataSilosDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"input_identifier": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"identifiers": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"actions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"headers": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_secret": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceEnricherCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		CreateEnricher struct {
			Enricher Enricher
		} `graphql:"createEnricher(input: {title: $title, type: SERVER, description: $description, url: $url, inputIdentifier: $inputIdentifier, headers: $headers, identifiers: $identifiers, actions: $actions})"`
	}

	vars := map[string]interface{}{
		"title":           graphql.String(d.Get("title").(string)),
		"description":     graphql.String(d.Get("description").(string)),
		"url":             graphql.String(d.Get("url").(string)),
		"inputIdentifier": graphql.ID(d.Get("input_identifier").(string)),
		"headers":         toHeadersList(d.Get("headers").([]interface{})),
		"identifiers":     toIDList(d.Get("identifiers").([]interface{})),
		"actions":         toRequestActionList(d.Get("actions").([]interface{})),
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error creating " + d.Get("title").(string),
			Detail:   "Error when creating enricher: " + err.Error(),
		})
		return diags
	}

	d.SetId(string(mutation.CreateEnricher.Enricher.ID))
	return nil
}

func resourceEnricherRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return nil
}

func resourceEnricherUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return nil
}

func resourceEnricherDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return nil
}
