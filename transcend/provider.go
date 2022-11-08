package transcend

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TRANSCEND_URL", "https://api.transcend.io/"),
				Description: "The custom Transcend backend URL to talk to. Typically can be left to the default production URL.",
			},
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TRANSCEND_KEY", nil),
				Description: "The API Key to use to talk to Transcend. Ensure it has the scopes to perform whatever actions you need. Can be set using the TRANSCEND_KEY environment variable.",
			},
			"internal_sombra_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TRANSCEND_INTERNAL_SOMBRA_KEY", nil),
				Description: "The API Key to use to talk to a self-hosted sombra. Only used for enterprises with the self-hosted option",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"transcend_api_key":                       resourceAPIKey(),
			"transcend_data_point":                    resourceDataPoint(),
			"transcend_enricher":                      resourceEnricher(),
			"transcend_data_silo":                     resourceDataSilo(),
			"transcend_data_silo_connection":          resourceDataSiloConnection(),
			"transcend_schema_discovery_plugin":       resourceSchemaDiscoveryPlugin(),
			"transcend_content_classification_plugin": resourceContentClassificationPlugin(),
			"transcend_data_silo_discovery_plugin":    resourceDataSiloDiscoveryPlugin(),
			"transcend_data_point_discovery_plugin":   resourceDataPointDiscoveryPlugin()
		},
		DataSourcesMap: map[string]*schema.Resource{
			"transcend_identifier": dataSourceIdentifier(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	backendUrl := d.Get("url").(string)
	backendApiKey := d.Get("key").(string)
	sombraInternalKey := d.Get("internal_sombra_key").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if backendUrl == "" || backendApiKey == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to authenticate provider",
			Detail:   "Some fields are missing",
		})
		return nil, diags
	}

	graphQlUrl, err := url.JoinPath(backendUrl, "/graphql")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error generating graphql url",
			Detail:   "Details: " + err.Error(),
		})
		return nil, diags
	}

	return NewClient(graphQlUrl, backendApiKey, sombraInternalKey), nil
}
