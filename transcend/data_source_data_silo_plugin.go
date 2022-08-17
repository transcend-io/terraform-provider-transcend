package transcend

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDataSiloPlugin() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDataSiloPluginRead,
		Schema: map[string]*schema.Schema{
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
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_frequency_minutes": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The updated frequency with which we should schedule this plugin, in milliseconds",
			},
			"schedule_start_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The updated start time when we should start scheduling this plugin, in ISO format",
			},
			"schedule_now": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
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
	}
}
