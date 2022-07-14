package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
)

func resourceDataPoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataPointCreate,
		ReadContext:   resourceDataPointRead,
		UpdateContext: resourceDataPointUpdate,
		DeleteContext: resourceDataPointDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_silo_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_collection_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_suggestions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"suggested_query": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"request_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"enabled_actions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// TODO: schema for subdatapoint
			// "sub_data_points": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"name": &schema.Schema{
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 			"description": &schema.Schema{
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 			"categories": &schema.Schema{
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"name": &schema.Schema{
			// 							Type:     schema.TypeString,
			// 							Required: true,
			// 						},
			// 						"category": &schema.Schema{
			// 							Type:     schema.TypeString,
			// 							Required: true,
			// 						},
			// 					},
			// 				},
			// 			},
			// 			"purposes": &schema.Schema{
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"name": &schema.Schema{
			// 							Type:     schema.TypeString,
			// 							Required: true,
			// 						},
			// 						"purpose": &schema.Schema{
			// 							Type:     schema.TypeString,
			// 							Required: true,
			// 						},
			// 					},
			// 				},
			// 			},
			// 			"attributes": &schema.Schema{
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"key": &schema.Schema{
			// 							Type:     schema.TypeString,
			// 							Required: true,
			// 						},
			// 						"values": &schema.Schema{
			// 							Type:     schema.TypeString,
			// 							Required: true,
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// },
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDataPointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	//TODO: check whether it is a database integration (has different input fields)
	id, err := mutateDataPoint(client, d, diags)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error creating datapoint " + d.Get("name").(string),
			Detail:   err.Error(),
		})
		return diags
	}

	resourceDataPointRead(ctx, d, m)
	d.SetId(id)

	return nil
}

func resourceDataPointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var query struct {
		DataPoint struct {
			Title struct {
				DefaultMessage graphql.String
			}
			Description struct {
				DefaultMessage graphql.String
			}
		} `graphql:"dataPoint(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.ID(d.Get("id").(string)),
	}

	err := client.graphql.Query(context.Background(), &query, vars)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: sync up all fields
	d.Set("title", query.DataPoint.Title.DefaultMessage)
	d.Set("description", query.DataPoint.Description.DefaultMessage)

	return nil

	return nil
}

func resourceDataPointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	_, err := mutateDataPoint(client, d, diags)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating datapoint " + d.Get("name").(string),
			Detail:   err.Error(),
		})
		return diags
	}

	return resourceDataPointRead(ctx, d, m)
}

func resourceDataPointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		DeleteDataPoints struct {
			Success graphql.Boolean
		} `graphql:"deleteDataPoints(input: { ids: $ids })"`
	}

	ids := make([]graphql.ID, 1)
	ids[0] = graphql.ID(d.Get("id").(string))

	vars := map[string]interface{}{
		"ids": ids,
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error deleting datapoint " + d.Get("title").(string),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return nil
}

func mutateDataPoint(client *Client, d *schema.ResourceData, diags diag.Diagnostics) (string, error) {
	var mutation struct {
		UpdateOrCreateDataPoint struct {
			DataPoint struct {
				ID   graphql.String
				Name graphql.String
			}
		} `graphql:"updateOrCreateDataPoint(input: {dataSiloId: $dataSiloId, name: $name, title: $title, dataCollectionTag: $dataCollectionTag, description: $description, enabledActions: $enabledActions})"`
	}

	vars := map[string]interface{}{
		"dataSiloId":        graphql.ID(d.Get("data_silo_id").(string)),
		"name":              graphql.String(d.Get("name").(string)),
		"title":             graphql.String(d.Get("title").(string)),
		"dataCollectionTag": graphql.String(d.Get("data_collection_tag").(string)),
		"description":       graphql.String(d.Get("description").(string)),
		"enabledActions":    toRequestActionObjectResolverList(d.Get("enabled_actions").([]interface{})),
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		return "", err
	}
	return string(mutation.UpdateOrCreateDataPoint.DataPoint.ID), nil
}
