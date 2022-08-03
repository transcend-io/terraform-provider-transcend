package transcend

import (
	"context"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
)

func createDataSiloUpdatableFields(d *schema.ResourceData) DataSiloUpdatableFields {
	return DataSiloUpdatableFields{
		Title:              graphql.String(d.Get("title").(string)),
		Description:        graphql.String(d.Get("description").(string)),
		URL:                graphql.String(d.Get("url").(string)),
		NotifyEmailAddress: graphql.String(d.Get("notify_email_address").(string)),
		IsLive:             graphql.Boolean(d.Get("is_live").(bool)),
		OwnerEmails:        toStringList(d.Get("owner_emails").([]interface{})),
		Headers:            toCustomHeaderInputList((d.Get("headers").([]interface{}))),

		// TODO: Add more fields
		// DataSubjectBlockListIds: toStringList(d.Get("data_subject_block_list_ids")),
		// Identifiers:             toStringList(d.Get("identifiers").([]interface{})),
		// "api_key_id":                   graphql.ID(d.Get("api_key_id").(string)),
		// "depended_on_data_silo_titles": toStringList(d.Get("depended_on_data_silo_titles").([]interface{})),
		// "team_names":                   toStringList(d.Get("team_names").([]interface{})),
	}
}

func createDataSiloInput(d *schema.ResourceData) DataSiloInput {
	return DataSiloInput{
		Name:                    graphql.String(d.Get("type").(string)),
		DataSiloUpdatableFields: createDataSiloUpdatableFields(d),
	}
}

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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The title of the data silo",
			},
			"link": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to the data silo",
			},
			"aws_external_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The external ID for the AWS IAM Role for AWS data silos",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of silo",
				ForceNew:    true,
			},
			"has_avc_functionality": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the data silo supports automated vendor coordination",
			},
			"headers": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Custom headers to include in outbound webhook",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the custom header",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
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
			},
			"outer_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The catalog name responsible for the cosmetics of the integration (name, description, logo, email fields)",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the data silo",
			},
			// "prompt_email_template_id": &schema.Schema{
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The id of template to use when prompting via email",
			// },
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the server to post to if a server silo",
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email address that should be notified whenever new requests are made",
			},
			"is_live": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the data silo should be live",
			},
			// "api_key_id": &schema.Schema{
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The id of the existing api key to attach to",
			// },
			// "identifiers": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// 	Description: "The names of the identifiers that the data silo should be connected to",
			// },
			// "depended_on_data_silo_ids": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// 	Description: "The IDs of the data silo that this data silo depends on during a deletion request.",
			// },
			// "data_subject_block_list_ids": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// 	Description: "The list of subject IDs to block list from this data silo",
			// },
			"owner_emails": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The emails of the users to assign as owners of this data silo. These emails must have matching users on Transcend.",
			},
			// "team_names": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// 	Description: "The names of the teams that should be responsible for this data silo",
			// },
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDataSilosCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	// Determine the type of the data silo. Most often, this is just the `type` field.
	// But for AVC silos, the `outer_type` actually contains the name to use, as the `type`
	// is always "promptAPerson"
	integrationName := d.Get("outer_type")
	if integrationName == "" {
		integrationName = d.Get("type")
	}

	// Create an empty data silo
	var createMutation struct {
		CreateDataSilos struct {
			DataSilos []DataSilo
		} `graphql:"createDataSilos(input: [{name: $name}])"`
	}
	createVars := map[string]interface{}{
		"name": graphql.String(integrationName.(string)),
	}
	err := client.graphql.Mutate(context.Background(), &createMutation, createVars)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error connecting to " + d.Get("type").(string),
			Detail:   "Error when connecting to data silo: " + err.Error(),
		})
		return diags
	}

	if len(createMutation.CreateDataSilos.DataSilos) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create data silo of type " + d.Get("type").(string),
			Detail:   "The request to create the silo completed, but no data was returned.",
		})
		return diags
	}
	d.SetId(string(createMutation.CreateDataSilos.DataSilos[0].ID))

	// Update the data silo with all fields
	resourceDataSilosUpdate(ctx, d, m)

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

	d.Set("id", query.DataSilo.ID)
	d.Set("link", query.DataSilo.Link)
	d.Set("aws_external_id", query.DataSilo.ExternalId)
	d.Set("has_avc_functionality", query.DataSilo.Catalog.HasAvcFunctionality)
	d.Set("type", query.DataSilo.Type)
	d.Set("title", query.DataSilo.Title)
	d.Set("description", query.DataSilo.Description)
	d.Set("url", query.DataSilo.URL)
	d.Set("outer_type", query.DataSilo.OuterType)
	d.Set("notify_email_address", query.DataSilo.NotifyEmailAddress)
	d.Set("is_live", query.DataSilo.IsLive)
	d.Set("owner_emails", flattenOwners(query.DataSilo))
	d.Set("headers", flattenHeaders(&query.DataSilo.Headers))

	// TODO: Support these fields being read in
	// d.Set("data_subject_block_list", flattenDataSiloBlockList(query.DataSilo))
	// d.Set("identifiers", query.DataSilo.Identifiers)
	// d.Set("prompt_email_template_id", query.DataSilo.PromptEmailTemplate.ID)
	// d.Set("team_names", ...)
	// d.Set("depended_on_data_silo_ids", ...)
	// d.Set("data_subject_block_list_ids", ...)
	// d.Set("headers", ...)
	// d.Set("api_key_id", ...)

	return nil
}

func resourceDataSilosUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var mutation struct {
		UpdateDataSilo struct {
			DataSilo DataSilo
		} `graphql:"updateDataSilo(input: $input)"`
	}

	vars := map[string]interface{}{
		"input": UpdateDataSiloInput{
			Id:                      graphql.ID(d.Get("id").(string)),
			DataSiloUpdatableFields: createDataSiloUpdatableFields(d),
		},
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
