package transcend

import (
	"bytes"
	"context"
	"io"
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
			// TODO: What if we just had the API here be formItems as a schema.TypeMap and the provider
			// queried the catalog for if the values should be secret or not? In the statefile, all values would be secretive,
			// but the provider could separate out plaintext from secret context values and give better error messages if there are
			// missing fields or if invalid field names are provided.
			"plaintext_context": &schema.Schema{
				Type:        schema.TypeSet,
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
			"secret_context": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "This is where you put values that go in the form when connecting a data silo. In general, most form values are secret context.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the input",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the input in plaintext",
							Sensitive:   true,
						},
					},
				},
			},
			"plugin_configuration": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "This is where you configure how often you'd like data silo and data point plugins to run, if enabled.",
				MinItems:    0,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "State to toggle plugin to",
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of plugin",
						},
						"schedule_frequency_minutes": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The updated frequency with which we should schedule this plugin, in milliseconds",
						},
						"schedule_start_at": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The updated start time when we should start scheduling this plugin, in ISO format",
						},
						"schedule_now": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether we should schedule a run immediately after this request",
						},
						"last_enabled_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date at which this data silo was last enabled",
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

	// Read the data silo information
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

	// Read the data silo plugin information
	var pluginQuery struct {
		Plugins struct {
			Plugins []types.Plugin
		} `graphql:"plugins(filterBy: { dataSiloId: $dataSiloId })"`
	}
	pluginVars := map[string]interface{}{
		"dataSiloId": graphql.String(d.Get("id").(string)),
	}
	err = client.graphql.Query(context.Background(), &pluginQuery, pluginVars, graphql.OperationName("Plugins"))
	if err != nil {
		return diag.FromErr(err)
	}
	if len(pluginQuery.Plugins.Plugins) == 1 {
		types.ReadDataSiloPluginIntoState(d, pluginQuery.Plugins.Plugins[0])
	}

	return nil
}

func resourceDataSilosUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	// Perform updates to most fields on the data silo
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
		deletionDiags := resourceDataSilosDelete(ctx, d, m)
		if deletionDiags.HasError() {
			diags = append(diags, deletionDiags...)
		}
		return diags
	}

	// Presign the SaaS context if the integration has secrets
	// For Internal Transcend Folks, see: https://docs.google.com/document/d/1PURNdW7VI9r9kwDM4fud9Hx_58vZbhMhB8OPEYxl8O4/view#
	var saasContext []byte
	if d.Get("secret_context") != nil {
		// Lookup the sombra URL to talk to
		var sombraUrlQuery struct {
			Organization struct {
				Sombra struct {
					CustomerUrl  graphql.String `graphql:"customerUrl"`
					HostedMethod graphql.String `graphql:"hostedMethod"`
				} `graphql:"sombra"`
			} `graphql:"organization"`
		}
		err = client.graphql.Query(context.Background(), &sombraUrlQuery, map[string]interface{}{}, graphql.OperationName("SombraUrlQuery"))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error Finding sombra URL",
				Detail:   "Error when updating data silo: " + err.Error(),
			})
			deletionDiags := resourceDataSilosDelete(ctx, d, m)
			if deletionDiags.HasError() {
				diags = append(diags, deletionDiags...)
			}
			return diags
		}
		// Lookup the saas context metadata
		var catalogQuery struct {
			Catalog struct {
				Catalog types.Catalog `json:"catalog"`
			} `graphql:"catalog(input: { integrationName: $integrationName })"`
		}
		err = client.graphql.Query(context.Background(), &catalogQuery, map[string]interface{}{
			"integrationName": graphql.String(types.GetIntegrationName(d)),
		}, graphql.OperationName("catalog"))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error Finding saas context metadata",
				Detail:   "Error when updating data silo: " + err.Error(),
			})
			deletionDiags := resourceDataSilosDelete(ctx, d, m)
			if deletionDiags.HasError() {
				diags = append(diags, deletionDiags...)
			}
			return diags
		}
		// Have sombra encrypt the secret map and parse the resulting saas context
		allowedBaseHosts := catalogQuery.Catalog.Catalog.IntegrationConfig.ConfiguredBaseHosts.PROD
		jsonBody, err := types.ConstructSecretMapString(d, allowedBaseHosts, catalogQuery.Catalog.Catalog.PlaintextInformation)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error encoding secret map to create saas context",
				Detail:   "Error when updating data silo: " + err.Error(),
			})
			deletionDiags := resourceDataSilosDelete(ctx, d, m)
			if deletionDiags.HasError() {
				diags = append(diags, deletionDiags...)
			}
			return diags
		}
		sombraResponse, err := client.sombraClient.Post(string(sombraUrlQuery.Organization.Sombra.CustomerUrl)+"v1/register-saas", "application/json", bytes.NewReader(jsonBody))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error creating SaaS context for the secret values",
				Detail:   "Error when updating data silo: " + err.Error(),
			})
			deletionDiags := resourceDataSilosDelete(ctx, d, m)
			if deletionDiags.HasError() {
				diags = append(diags, deletionDiags...)
			}
			return diags
		}
		saasContext, err = io.ReadAll(sombraResponse.Body)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading response of the new SaaS context",
				Detail:   "Error when updating data silo: " + err.Error(),
			})
			deletionDiags := resourceDataSilosDelete(ctx, d, m)
			if deletionDiags.HasError() {
				diags = append(diags, deletionDiags...)
			}
			return diags
		}
	}

	// Optionally attempt to connect the data silo, setting the form fields on success
	shouldSkipConnecting := d.Get("skip_connecting").(bool)
	if !shouldSkipConnecting {
		var connectMutation struct {
			ReconnectDataSilo struct {
				DataSilo types.DataSilo
			} `graphql:"reconnectDataSilo(input: $input, dhEncrypted: $dhEncrypted)"`
		}
		connectVars := map[string]interface{}{
			"input":       types.CreateReconnectDataSiloFields(d, saasContext),
			"dhEncrypted": graphql.String(""), // This is not needed when no encrypted saas contexts are provided
		}
		err = client.graphql.Mutate(context.Background(), &connectMutation, connectVars, graphql.OperationName("ReconnectDataSilo"))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error connecting data silos",
				Detail:   "Error when connecting data silo: " + err.Error(),
			})
			deletionDiags := resourceDataSilosDelete(ctx, d, m)
			if deletionDiags.HasError() {
				diags = append(diags, deletionDiags...)
			}
			return diags
		}
	}

	// Handle the plugin settings if defined
	if d.Get("plugin_configuration") != nil && len(d.Get("plugin_configuration").([]interface{})) == 1 {
		// Read the data silo plugin information
		var pluginQuery struct {
			Plugins struct {
				Plugins []types.Plugin
			} `graphql:"plugins(filterBy: { dataSiloId: $dataSiloId })"`
		}
		pluginVars := map[string]interface{}{
			"dataSiloId": graphql.String(d.Get("id").(string)),
		}
		err = client.graphql.Query(context.Background(), &pluginQuery, pluginVars, graphql.OperationName("Plugins"))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error finding data silo plugin for data silo",
				Detail:   "Error when reading data silo plugin: " + err.Error(),
			})
			return diags
		}
		if len(pluginQuery.Plugins.Plugins) != 1 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error finding exactly one data silo plugin for data silo",
				Detail:   "Error when reading data silo plugin",
			})
			return diags
		}

		var updateMutation struct {
			UpdateDataSiloPlugin struct {
				Plugin types.Plugin
			} `graphql:"updateDataSiloPlugin(input: $input)"`
		}
		updateVars := map[string]interface{}{
			"input": types.MakeUpdatePluginInput(d, pluginQuery.Plugins.Plugins[0].ID),
		}

		err := client.graphql.Mutate(context.Background(), &updateMutation, updateVars, graphql.OperationName("UpdateDataSiloPlugin"))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error updating data silo plugin",
				Detail:   "Error when updating data silo plugin: " + err.Error(),
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
