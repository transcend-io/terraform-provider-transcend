package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/machinebox/graphql"
)

func dataDataSilo() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataDataSilosRead,
		// UpdateContext: resourceDataSilosUpdate,
		// DeleteContext: resourceDataSilosDelete,
		Schema: map[string]*schema.Schema{
			// "last_updated": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// 	Computed: true,
			// },
			// "data_silos": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Required: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": &schema.Schema{
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 				Computed: true,
			// 			},
			// 			"title": &schema.Schema{
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 				Computed: true,
			// 			},
			// 			"type": &schema.Schema{
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 				Computed: true,
			// 			},
			// 			"link": &schema.Schema{
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 				Computed: true,
			// 			},
			// 			"catalog": &schema.Schema{
			// 				Type:     schema.TypeMap,
			// 				Optional: true,
			// 				Computed: true,
			// 				Elem: &schema.Schema{
			// 					Type: schema.TypeBool,
			// 				},
			// 			},
			// 		},
			// 	},
			// },
			"text": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"first": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"offset": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func DataDataSilosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// var query struct {
	// 	dataSilos DataSilo `graphql:"dataSilos(filterBy: { text: $text }, first: $first, offset: $offset)"`
	// }

	req := graphql.NewRequest(`
		dataSilos(
			filterBy: { text: $title }
			first: $first
			offset: $offset
		) {
			nodes {
			id
			title
			link
			type
			catalog {
				hasAvcFunctionality
			}
		}
	  }`)

	req.Var("text", d.Get("text").(string))
	req.Var("first", d.Get("first").(int))
	req.Var("offset", d.Get("offset").(int))
	req.Header.Set("Authorization", client.apiToken)

	err := client.graphql.Run(context.Background(), req, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
