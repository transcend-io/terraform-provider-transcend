package transcend

import (
	"context"
	"time"

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
			"plugin_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"data_silo_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the data silo to connect",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of plugin",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "State to toggle plugin to",
			},
			"schedule_frequency": &schema.Schema{
				Type:        schema.TypeString,
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
			"scheduled_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the next plugin run is scheduled",
			},
			"last_run_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date at which this data silo was last run",
			},
			"last_enabled_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date at which this data silo was last enabled",
			},
			"error": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current error message for the most recent run of the plugin",
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

	plugin, err := types.PluginsReadQuery(*client.graphql, d)
	if err != "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading plugin",
			Detail:   err,
		})
		return diags
	}

	types.ReadDataSiloPluginIntoState(d, plugin)
	d.SetId(string(plugin.ID))

	return nil
}

func resourceDataSiloPluginUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var updateMutation struct {
		UpdateDataSiloPluginPayload struct {
			Plugin types.Plugin
		} `graphql:"updateDataSiloPlugin(input: $input)"`
	}

	dateTime, err := time.Parse("2022-08-16T07:00:00.000Z", d.Get("schedule_start_at").(string))
	date := dateTime.String()
	if err != nil {
		plugin, err := types.PluginsReadQuery(*client.graphql, d)
		if err != "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error updating plugin",
				Detail:   err,
			})
			return diags
		}
		date = string(plugin.ScheduleStartAt)
	}

	updateVars := map[string]interface{}{
		"input": types.MakeUpdatePluginInput(d, date),
	}

	err = client.graphql.Mutate(context.Background(), &updateMutation, updateVars, graphql.OperationName("UpdateDataSiloPlugin"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating plugin",
			Detail:   "Error when updating plugin: " + err.Error(),
		})
		return diags
	}

	return resourceDataSiloPluginRead(ctx, d, m)
}

// Data silos cannot be disconnected, so just no-op
func resourceDataSiloPluginDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
