package transcend

import (
	"context"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
)

func resourceDataSilo() *schema.Resource {
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
				Optional: true,
				Computed: true,
			},
			"link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"has_avc_functionality": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
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
			"outer_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"prompt_email_template_id": &schema.Schema{
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
			"identifiers": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"depended_on_data_silo_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"depended_on_data_silo_titles": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"data_subject_block_list_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_emails": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"teams": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"team_names": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
		} `graphql:"connectDataSilo(input: {name: $type, headers: $headers, outerType: $outer_type, title: $title, description: $description, url: $url, notifyEmailAddress: $notify_email_address, isLive: $is_live, apiKeyId: $api_key_id, identifiers: $identifiers, dependedOnDataSiloTitles: $depended_on_data_silo_titles, ownerEmails: $owner_emails, teamNames: $team_names})"`
	}

	heads := d.Get("headers").([]interface{})

	headers := make([]CustomHeaderInput, len(heads))

	for i, head := range heads {

		newHead := head.(map[string]interface{})

		headers[i] = CustomHeaderInput{
			graphql.String(newHead["name"].(string)),
			graphql.String(newHead["value"].(string)),
			graphql.Boolean(newHead["is_secret"].(bool)),
		}
	}

	vars := map[string]interface{}{
		"type":                         graphql.String(d.Get("type").(string)),
		"headers":                      headers,
		"outer_type":                   graphql.String(d.Get("outer_type").(string)),
		"title":                        graphql.String(d.Get("title").(string)),
		"description":                  graphql.String(d.Get("description").(string)),
		"url":                          graphql.String(d.Get("url").(string)),
		"notify_email_address":         graphql.String(d.Get("notify_email_address").(string)),
		"is_live":                      graphql.Boolean(d.Get("is_live").(bool)),
		"api_key_id":                   graphql.ID(d.Get("api_key_id").(string)),
		"identifiers":                  toStringList(d.Get("identifiers").([]interface{})),
		"depended_on_data_silo_titles": toStringList(d.Get("depended_on_data_silo_titles").([]interface{})),
		"owner_emails":                 toStringList(d.Get("owner_emails").([]interface{})),
		"team_names":                   toStringList(d.Get("team_names").([]interface{})),
	}

	err := client.graphql.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error connecting to " + d.Get("type").(string),
			Detail:   "Error when connecting to data silo: " + err.Error(),
		})
		return diags
	}

	d.SetId(string(mutation.ConnectDataSilo.DataSilo.ID))

	resourceDataSilosRead(ctx, d, m)

	return nil
}

func resourceDataSilosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var query struct {
		DataSilo DataSilo `graphql:"dataSilo(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.String(d.Get("id").(string)),
	}

	err := client.graphql.Query(context.Background(), &query, vars)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("title", query.DataSilo.Title)
	d.Set("link", query.DataSilo.Link)
	d.Set("type", query.DataSilo.Type)
	d.Set("has_avc_functionality", query.DataSilo.Catalog.HasAvcFunctionality)

	return nil
}

func resourceDataSilosUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		UpdateDataSilo struct {
			DataSilo DataSilo
		} `graphql:"updateDataSilo(input: {id: $id, title: $title, description: $description, url: $url, notifyEmailAddress: $notify_email_address, isLive: $is_live, apiKeyId: $api_key_id identifiers: $identifiers, dependedOnDataSiloTitles: $depended_on_data_silo_titles, ownerEmails: $owner_emails, teamNames: $team_names})"`
	}

	vars := map[string]interface{}{
		"id":                           graphql.ID(d.Get("id").(string)),
		"title":                        graphql.String(d.Get("title").(string)),
		"description":                  graphql.String(d.Get("description").(string)),
		"url":                          graphql.String(d.Get("url").(string)),
		"notify_email_address":         graphql.String(d.Get("notify_email_address").(string)),
		"is_live":                      graphql.Boolean(d.Get("is_live").(bool)),
		"api_key_id":                   graphql.ID(d.Get("api_key_id").(string)),
		"identifiers":                  toStringList(d.Get("identifiers").([]interface{})),
		"depended_on_data_silo_titles": toStringList(d.Get("depended_on_data_silo_titles").([]interface{})),
		"owner_emails":                 toStringList(d.Get("owner_emails").([]interface{})),
		"team_names":                   toStringList(d.Get("team_names").([]interface{})),
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

	return resourceDataSilosRead(ctx, d, m)
}

func resourceDataSilosDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		DeleteDataSilos struct {
			Success graphql.Boolean
		} `graphql:"deleteDataSilos(input: { ids: $ids })"`
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
			Summary:  "Error deleting data silo " + d.Get("type").(string),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")
	return nil
}