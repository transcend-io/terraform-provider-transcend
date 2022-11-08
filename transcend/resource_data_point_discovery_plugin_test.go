package transcend

import (
	"testing"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/gruntwork-io/terratest/modules/terraform"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
)

func prepareDataPointDiscoveryPluginOptions(t *testing.T, vars map[string]interface{}) *terraform.Options {
	defaultVars := map[string]interface{}{"title": t.Name()}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/data_point_discovery_plugin",
		Vars:         defaultVars,
	})
	return terraformOptions
}

func deployDataPointDiscoveryPlugin(t *testing.T, terraformOptions *terraform.Options) (types.DataSilo, []types.Plugin) {
	// TODO: Use the Idempotent version eventually
	terraform.InitAndApply(t, terraformOptions)
	// terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	silo := lookupDataSilo(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	plugin := lookupDataSiloPlugin(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	return silo, plugin
}

func TestCanUseSeparateDataPointDiscoveryPluginResource(t *testing.T) {
	options := prepareDataPointDiscoveryPluginOptions(t, map[string]interface{}{
		"title": t.Name(),
		"plugin_config": map[string]interface{}{
			"enabled":                    true,
			"schedule_frequency_minutes": 120,
			// Schedule far in the future so that the test works for a long time
			"schedule_start_at": "2122-09-06T17:51:13.000Z",
			"schedule_now":      false,
		},
	})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataPointDiscoveryPlugin(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.NotEmpty(t, terraform.Output(t, options, "awsExternalId"))
}
