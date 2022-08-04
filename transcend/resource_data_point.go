package transcend

import (
	"context"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the data silo to create the datapoint for",
				ForceNew:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "he datapoint name (used to key by)",
			},
			"title": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the datapoint",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description for the datapoint",
			},
			// "data_collection_tag": &schema.Schema{
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The title of the data collection to assign to the datapoint. If the collection does not exist, one will be created.",
			// },
			// "query_suggestions": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"suggested_query": &schema.Schema{
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 			"request_type": &schema.Schema{
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 		},
			// 	},
			// 	Description: "The suggested SQL queries to run for a DSR",
			// },
			// "enabled_actions": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// 	Description: "The actions that the datapoint should connect to",
			// },
			// "sub_data_points": &schema.Schema{
			// 	Type:        schema.TypeList,
			// 	Optional:    true,
			// 	Description: "The subdatapoints associated with this datapoint",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"name": &schema.Schema{
			// 				Type:        schema.TypeString,
			// 				Required:    true,
			// 				Description: "The name of the subdatapoint",
			// 			},
			// 			"description": &schema.Schema{
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 				Description: "A description for the subdatapoint",
			// 			},
			// 			"categories": &schema.Schema{
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"name": &schema.Schema{
			// 							Type:        schema.TypeString,
			// 							Required:    true,
			// 							Description: "The name of the subcategory",
			// 						},
			// 						"category": &schema.Schema{
			// 							Type:        schema.TypeString,
			// 							Required:    true,
			// 							Description: "The category of personal data",
			// 						},
			// 					},
			// 				},
			// 				Description: "The category of personal data for this subdatapoint",
			// 			},
			// 			"purposes": &schema.Schema{
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"name": &schema.Schema{
			// 							Type:        schema.TypeString,
			// 							Required:    true,
			// 							Description: "The purpose of processing sub category",
			// 						},
			// 						"purpose": &schema.Schema{
			// 							Type:        schema.TypeString,
			// 							Required:    true,
			// 							Description: "The purpose of processing",
			// 						},
			// 					},
			// 				},
			// 				Description: "The processing purposes for this subdatapoint",
			// 			},
			// 			"attributes": &schema.Schema{
			// 				Type:     schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"key": &schema.Schema{
			// 							Type:        schema.TypeString,
			// 							Required:    true,
			// 							Description: "The attribute key that houses the attribute values",
			// 						},
			// 						"values": &schema.Schema{
			// 							Type:     schema.TypeList,
			// 							Required: true,
			// 							Elem: &schema.Schema{
			// 								Type: schema.TypeString,
			// 							},
			// 							Description: "The attribute values used to label resources",
			// 						},
			// 					},
			// 				},
			// 				Description: "The attribute values used to label this subdatapoint",
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

	var mutation struct {
		CreateApiKey struct {
			DataPoint types.DataPoint
		} `graphql:"updateOrCreateDataPoint(input: $input)"`
	}

	vars := map[string]interface{}{
		"input": types.MakeUpdateOrCreateDataPointInput(d),
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error creating Data Point",
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId(string(mutation.CreateApiKey.DataPoint.ID))

	resourceDataPointRead(ctx, d, m)

	return nil
}

func resourceDataPointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var query struct {
		DataPoints struct {
			Nodes []types.DataPoint
		} `graphql:"dataPoints(filterBy: { ids: [$id] })"`
	}

	vars := map[string]interface{}{
		"id": graphql.ID(d.Get("id").(string)),
	}

	err := client.graphql.Query(context.Background(), &query, vars)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading datapoint " + d.Get("name").(string),
			Detail:   err.Error(),
		})
		return diags
	}

	if len(query.DataPoints.Nodes) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading datapoint " + d.Get("name").(string),
			Detail:   "Cannot find datapoint.",
		})
		return diags
	}

	types.ReadDataPointIntoState(d, query.DataPoints.Nodes[0])

	return nil
}

func resourceDataPointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		UpdateApiKey struct {
			DataPoint types.DataPoint
		} `graphql:"updateOrCreateDataPoint(input: $input)"`
	}

	vars := map[string]interface{}{
		"input": types.MakeUpdateOrCreateDataPointInput(d),
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating Data Point",
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
