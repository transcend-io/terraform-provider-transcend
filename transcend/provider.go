package transcend

import (
	"context"

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
				DefaultFunc: schema.EnvDefaultFunc("TRANSCEND_URL", nil),
			},
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TRANSCEND_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"transcend_data_silo": resourceDataSilos(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"transcend_data_silo": dataDataSilo(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	url := d.Get("url").(string)
	key := d.Get("key").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if url == "" || key == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to authenticate provider",
			Detail:   "Some fields are missing",
		})
		return nil, diags
	}

	url = url + "graphql"

	return NewClient(url, key), nil
}
