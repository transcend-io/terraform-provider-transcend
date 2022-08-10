package transcend

import (
	"context"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

func resourceEnricher() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnricherCreate,
		ReadContext:   resourceEnricherRead,
		UpdateContext: resourceEnricherUpdate,
		DeleteContext: resourceEnricherDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The enricher's title",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The enricher's title",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The enricher's description",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The url that the enricher should post to",
			},
			"input_identifier": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the identifier that will be the input to the enricher",
			},
			"output_identifiers": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IDs of the identifiers that can possibly be output from the enricher",
			},
			"actions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The action types that the enricher should run for",
			},
			"headers": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the custom header",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The value of the custom header",
						},
						"is_secret": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When true, the value of this header will be considered sensitive",
						},
					},
				},
				Description: "Custom headers to include in outbound webhook",
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
			Enricher types.Enricher
		} `graphql:"createEnricher(input: $input)"`
	}

	vars := map[string]interface{}{
		"input": types.MakeEnricherInput(d),
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars, graphql.OperationName("CreateEnricher"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error creating " + d.Get("title").(string),
			Detail:   "Error when creating enricher: " + err.Error(),
		})
		return diags
	}
	d.SetId(string(mutation.CreateEnricher.Enricher.ID))

	return resourceEnricherRead(ctx, d, m)
}

func resourceEnricherRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var query struct {
		Enricher types.Enricher `graphql:"enricher(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.ID(d.Get("id").(string)),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("Enricher"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading " + d.Get("title").(string),
			Detail:   "Error when reading enricher: " + err.Error(),
		})
		return diags
	}

	types.ReadEnricherIntoState(d, query.Enricher)

	return nil
}

func resourceEnricherUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		UpdateEnricher struct {
			Enricher types.Enricher
		} `graphql:"updateEnricher(input: $input)"`
	}

	vars := map[string]interface{}{
		"input": types.MakeUpdateEnricherInput(d),
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars, graphql.OperationName("UpdateEnricher"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating " + d.Get("title").(string),
			Detail:   "Error when updating enricher: " + err.Error(),
		})
		return diags
	}

	return resourceEnricherRead(ctx, d, m)
}

func resourceEnricherDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		DeleteEnricher struct {
			Success graphql.Boolean
		} `graphql:"deleteEnricher(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.ID(d.Get("id").(string)),
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars, graphql.OperationName("DeleteEnricher"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error deleting enricher " + d.Get("title").(string),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return nil
}
