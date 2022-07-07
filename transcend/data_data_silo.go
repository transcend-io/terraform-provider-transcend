package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
)

func dataDataSilo() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataDataSilosRead,
		// UpdateContext: resourceDataSilosUpdate,
		// DeleteContext: resourceDataSilosDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"data_silos": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
							Elem: &schema.Schema{
								Type: schema.TypeBool,
							},
						},
					},
				},
			},
			"text": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"first": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"offset": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func DataDataSilosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var data struct {
		DataSilos struct {
			Nodes []DataSilo
		} `graphql:"dataSilos(filterBy: { text: $text } first: $first offset: $offset)"`
	}

	vars := map[string]interface{}{
		"text":   graphql.String(d.Get("text").(string)),
		"first":  graphql.Int(d.Get("first").(int)),
		"offset": graphql.Int(d.Get("offset").(int)),
	}

	err := client.graphql.Query(context.Background(), &data, vars)
	if err != nil {
		return diag.FromErr(err)
	}

	res := flattenItems(&data.DataSilos.Nodes)
	d.Set("data_silos", res)

	d.SetId("datasilo")
	return nil
}
