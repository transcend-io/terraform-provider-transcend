package transcend

import (
	"context"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

func resourceDiscoClassScanConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDiscoClassScanConfigCreate,
		ReadContext:   resourceDiscoClassScanConfigRead,
		UpdateContext: resourceDiscoClassScanConfigUpdate,
		DeleteContext: resourceDiscoClassScanConfigDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_silo_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the data silo to connect",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of disco class scan config",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not scheduling is enabled",
			},
			"schedule_frequency_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"last_disco_class_scan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the last disco class scan",
			},
			"schedule_start_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDiscoClassScanConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceDiscoClassScanConfigUpdate(ctx, d, m)
}

func resourceDiscoClassScanConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	var discoClassScanConfigQuery struct {
		DiscoClassScanConfig types.DiscoClassScanConfig `graphql:"discoClassScanConfig(input: { dataSiloId: $dataSiloId })"`
	}
	discoClassScanConfigVars := map[string]interface{}{
		"dataSiloId": graphql.ID(d.Get("data_silo_id").(string)),
	}
	err := client.graphql.Query(context.Background(), &discoClassScanConfigQuery, discoClassScanConfigVars, graphql.OperationName("DiscoClassScanConfig"))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(discoClassScanConfigQuery.DiscoClassScanConfig.ID))

	return diags
}

func resourceDiscoClassScanConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	// DiscoClassScanConfig already exists when a data silo is created, it is just disabled.
	// So here we fetch the ID / settings on the existing discoClassScanConfig
	// Read the data silo discoClassScanConfig information
	var discoClassScanConfigQuery struct {
		DiscoClassScanConfig types.DiscoClassScanConfig `graphql:"discoClassScanConfig(input: { dataSiloId: $dataSiloId })"`
	}
	discoClassScanConfigVars := map[string]interface{}{
		"dataSiloId": graphql.ID(d.Get("data_silo_id").(string)),
	}
	err := client.graphql.Query(context.Background(), &discoClassScanConfigQuery, discoClassScanConfigVars, graphql.OperationName("DiscoClassScanConfig"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error finding discoClassScanConfig for data silo",
			Detail:   "Error when reading discoClassScanConfig: " + err.Error(),
		})
		return diags
	}

	d.SetId(string(discoClassScanConfigQuery.DiscoClassScanConfig.ID))

	var updateMutation struct {
		UpdateDiscoClassScanConfig struct {
			DiscoClassScanConfig types.DiscoClassScanConfig
		} `graphql:"updateDiscoClassScanConfig(input: $input)"`
	}
	updateVars := map[string]interface{}{
		"input": types.MakeStandaloneUpdateDiscoClassScanConfigInput(d),
	}
	err = client.graphql.Mutate(context.Background(), &updateMutation, updateVars, graphql.OperationName("UpdateDataSiloPlugin"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating disco class scan config",
			Detail:   "Error when updating disco class scan config: " + err.Error(),
		})
		return diags
	}

	return resourceDiscoClassScanConfigRead(ctx, d, m)
}

// Scans cannot be deleted, but they can be disabled, so we do that here
func resourceDiscoClassScanConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	resourceDiscoClassScanConfigRead(ctx, d, m)

	var updateMutation struct {
		UpdateDiscoClassScanConfig struct {
			DiscoClassScanConfig types.DiscoClassScanConfig
		} `graphql:"updateDiscoClassScanConfig(input: $input)"`
	}
	input := types.MakeStandaloneUpdateDiscoClassScanConfigInput(d)
	input.Enabled = false
	updateVars := map[string]interface{}{
		"input": input,
	}
	err := client.graphql.Mutate(context.Background(), &updateMutation, updateVars, graphql.OperationName("UpdateDataSiloPlugin"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating disco class scan config",
			Detail:   "Error when updating disco class scan config: " + err.Error(),
		})
		return diags
	}

	return nil
}
