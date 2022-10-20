package transcend

import (
	"context"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

func resourceDataSiloPlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataSiloPluginCreate,
		ReadContext:   resourceDataSiloPluginRead,
		UpdateContext: resourceDataSiloPluginUpdate,
		DeleteContext: resourceDataSiloPluginDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_silo_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the data silo to connect",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "State to toggle plugin to",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of plugin",
				ForceNew:    true,
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDataSiloPluginCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceDataSiloPluginUpdate(ctx, d, m)
}

func resourceDataSiloPluginRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	var pluginQuery struct {
		Plugins struct {
			Plugins []types.Plugin
		} `graphql:"plugins(filterBy: { dataSiloId: $dataSiloId })"`
	}
	pluginVars := map[string]interface{}{
		"dataSiloId": graphql.String(d.Get("data_silo_id").(string)),
	}
	err := client.graphql.Query(context.Background(), &pluginQuery, pluginVars, graphql.OperationName("Plugins"))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(pluginQuery.Plugins.Plugins) == 1 {
		types.ReadStandaloneDataSiloPluginIntoState(d, pluginQuery.Plugins.Plugins[0])
		d.SetId(string(pluginQuery.Plugins.Plugins[0].ID))
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error querying plugin",
			Detail:   "Error when querying for data silo plugin: Found undexpected number of plugins",
		})
	}

	return diags
}

func resourceDataSiloPluginUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	// Plugins already exist when a data silo is created, they are just disabled.
	// So here we fetch the ID / settings on the existing plugin
	// Read the data silo plugin information
	var pluginQuery struct {
		Plugins struct {
			Plugins []types.Plugin
		} `graphql:"plugins(filterBy: { dataSiloId: $dataSiloId })"`
	}
	pluginVars := map[string]interface{}{
		"dataSiloId": graphql.String(d.Get("data_silo_id").(string)),
	}
	err := client.graphql.Query(context.Background(), &pluginQuery, pluginVars, graphql.OperationName("Plugins"))
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
	d.Set("id", pluginQuery.Plugins.Plugins[0].ID)

	var updateMutation struct {
		UpdateDataSiloPlugin struct {
			Plugin types.Plugin
		} `graphql:"updateDataSiloPlugin(input: $input)"`
	}
	updateVars := map[string]interface{}{
		"input": types.MakeStandaloneUpdatePluginInput(d),
	}
	err = client.graphql.Mutate(context.Background(), &updateMutation, updateVars, graphql.OperationName("UpdateDataSiloPlugin"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating data silo plugin",
			Detail:   "Error when updating data silo plugin: " + err.Error(),
		})
		return diags
	}

	return resourceDataSiloPluginRead(ctx, d, m)
}

// Plugins cannot be deleted, but they can be disabled, so we do that here
func resourceDataSiloPluginDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	resourceDataSiloPluginRead(ctx, d, m)

	var updateMutation struct {
		UpdateDataSiloPlugin struct {
			Plugin types.Plugin
		} `graphql:"updateDataSiloPlugin(input: $input)"`
	}
	input := types.MakeStandaloneUpdatePluginInput(d)
	input.Enabled = false
	updateVars := map[string]interface{}{
		"input": input,
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

	return nil
}
