package transcend

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
)

func resourceDataSilos() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataSilosCreate,
		ReadContext:   resourceDataSilosRead,
		UpdateContext: resourceDataSilosUpdate,
		DeleteContext: resourceDataSilosDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"silos": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"title": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"catalog": &schema.Schema{
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeBool,
										},
									},
									"description": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"url": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
											value := v.(string)

											var diags diag.Diagnostics
											if !strings.HasPrefix(value, "https://") && !strings.HasPrefix(value, "http://") {
												diags = append(diags, diag.Diagnostic{
													Severity: diag.Error,
													Summary:  "Invalid URL",
													Detail:   "URL did not start with 'https://' or 'https://'",
												})
											}
											return diags
										},
									},
									"notify_email_address": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"is_live": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"api_key_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"link": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDataSilosCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		ConnectDataSilo struct {
			DataSilo DataSilo
		} `graphql:"connectDataSilo(input: {name: $type, title: $title, description: $description, url: $url, notifyEmailAddress: $notify_email_address, isLive: $is_live, apiKeyId: $api_key_id})"`
	}

	silos := d.Get("silos").([]interface{})

	for _, item := range silos {
		i := item.(map[string]interface{})

		in := i["input"].([]interface{})[0]
		input := in.(map[string]interface{})

		vars := map[string]interface{}{
			"type":                 graphql.String(input["type"].(string)),
			"title":                graphql.String(input["title"].(string)),
			"description":          graphql.String(input["description"].(string)),
			"url":                  graphql.String(input["url"].(string)),
			"notify_email_address": graphql.String(input["notify_email_address"].(string)),
			"is_live":              graphql.Boolean(input["is_live"].(bool)),
			"api_key_id":           graphql.ID(input["api_key_id"].(string)),
		}

		err := client.graphql.Mutate(context.Background(), &mutation, vars)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error connecting data silos",
				Detail:   "Error when connecting to data silos: " + err.Error(),
			})
			return diags
		}
	}

	resourceDataSilosRead(ctx, d, m)

	d.SetId(string(mutation.ConnectDataSilo.DataSilo.ID))

	return nil
}

func resourceDataSilosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var data struct {
		DataSilos struct {
			Nodes []DataSilo
		}
	}

	err := client.graphql.Query(context.Background(), &data, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	res := flattenItems(&data.DataSilos.Nodes)
	d.Set("silos", res)

	return nil
}

func resourceDataSilosUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		UpdateDataSilo struct {
			DataSilo DataSilo
		} `graphql:"updateDataSilo(input: {id: $id, title: $title, description: $description, url: $url, notifyEmailAddress: $notify_email_address, isLive: $is_live, apiKeyId: $api_key_id})"`
	}

	if d.HasChange("silos") {
		silos := d.Get("silos").([]interface{})

		for _, item := range silos {
			i := item.(map[string]interface{})

			in := i["input"].([]interface{})[0]
			input := in.(map[string]interface{})

			vars := map[string]interface{}{
				"id":                   graphql.ID(input["id"].(string)),
				"title":                graphql.String(input["title"].(string)),
				"description":          graphql.String(input["description"].(string)),
				"url":                  graphql.String(input["url"].(string)),
				"notify_email_address": graphql.String(input["notify_email_address"].(string)),
				"is_live":              graphql.Boolean(input["is_live"].(bool)),
				"api_key_id":           graphql.ID(input["api_key_id"].(string)),
			}

			err := client.graphql.Mutate(context.Background(), &mutation, vars)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error updating data silos",
					Detail:   "Error when updating data silo: " + err.Error(),
				})
				return diags
			}
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}
	return resourceDataSilosRead(ctx, d, m)
}

func resourceDataSilosDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return nil
}
