package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDataSilos() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataSilosCreate,
		ReadContext:   resourceDataSilosRead,
		// UpdateContext: resourceDataSilosUpdate,
		// DeleteContext: resourceDataSilosDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"data_silos": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"title": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"link": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"catalog": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hasAvcFunctionality": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"identifiers": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"identifier": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"isLive": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
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

func resourceDataSilosCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// client = m.(*Client)
	return nil
}

func resourceDataSilosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// client := m.(*Client)

	// var query struct {
	// 	dataSilos DataSilo `graphql:"dataSilos(filterBy: { text: $title, ids: $ids }, first: $first, offset: $offset)"`
	// }

	// // vars := map[string]interface{}{
	// // 	"id":     graphql.String(d.Get("id").(string)),
	// // 	"title":  graphql.String(d.Get("title").(string)),
	// // 	"first":  graphql.Int(d.Get("first").(int)),
	// // 	"offset": graphql.Int(d.Get("offset").(int)),
	// // }

	// // err := client.graphql.Query(context.Background(), &query, vars)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	return nil
}
