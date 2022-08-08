package transcend

import (
	"context"
	"strings"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
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
			"plaintext_context": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "This is where you put non-secretive values that go in the form when connecting a data silo",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the plaintext input",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the plaintext input",
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
			"skip_connecting": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the data silo will be left unconnected. When false, the provided credentials will be tested against a live environment",
			},
			"connection_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current state of the integration",
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

	// Create an empty data silo
	var createMutation struct {
		CreateDataSilos struct {
			DataSilos []types.DataSilo
		} `graphql:"createDataSilos(input: [$dataSilo])"`
	}
	createVars := map[string]interface{}{
		"dataSilo": types.CreateDataSiloInput(d),
	}
	err := client.graphql.Mutate(context.Background(), &createMutation, createVars, graphql.OperationName("CreateDataSilos"))
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
	return resourceDataSilosUpdate(ctx, d, m)
}

func resourceDataSilosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var query struct {
		DataSilo types.DataSilo `graphql:"dataSilo(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.String(d.Get("id").(string)),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("DataSilo"))
	if err != nil {
		return diag.FromErr(err)
	}

	types.ReadDataSiloIntoState(d, query.DataSilo)

	return nil
}

func resourceDataSilosUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var updateMutation struct {
		UpdateDataSilo struct {
			DataSilo types.DataSilo
		} `graphql:"updateDataSilo(input: $input)"`
	}
	updateVars := map[string]interface{}{
		"input": types.UpdateDataSiloInput{
			Id:                      graphql.ID(d.Get("id").(string)),
			DataSiloUpdatableFields: types.CreateDataSiloUpdatableFields(d),
		},
	}
	err := client.graphql.Mutate(context.Background(), &updateMutation, updateVars, graphql.OperationName("UpdateDataSilo"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating data silos",
			Detail:   "Error when updating data silo: " + err.Error(),
		})
		return diags
	}

	shouldSkipConnecting := d.Get("skip_connecting").(bool)
	if !shouldSkipConnecting {
		var connectMutation struct {
			ReconnectDataSilo struct {
				DataSilo types.DataSilo
			} `graphql:"reconnectDataSilo(input: $input, dhEncrypted: $dhEncrypted)"`
		}
		connectVars := map[string]interface{}{
			"input":       types.CreateReconnectDataSiloFields(d),
			"dhEncrypted": graphql.String(""), // This is not needed when no encrypted saas contexts are provided
		}
		err = client.graphql.Mutate(context.Background(), &connectMutation, connectVars, graphql.OperationName("ReconnectDataSilo"))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error connecting data silos",
				Detail:   "Error when connecting data silo: " + err.Error(),
			})
			return diags
		}
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

	err := client.graphql.Mutate(context.Background(), &mutation, vars, graphql.OperationName("DeleteDataSilos"))
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
