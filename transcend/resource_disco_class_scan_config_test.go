package transcend

import (
	"testing"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/gruntwork-io/terratest/modules/terraform"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
)

func prepareDiscoClassScanConfigOptions(t *testing.T, vars map[string]interface{}) *terraform.Options {
	defaultVars := map[string]interface{}{"title": t.Name()}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/disco_class_scan_config",
		Vars:         defaultVars,
	})
	return terraformOptions
}

func deployDiscoClassScanConfig(t *testing.T, terraformOptions *terraform.Options) (types.DataSilo, types.DiscoClassScanConfig, []types.Plugin) {
	// TODO: Use the Idempotent version eventually
	terraform.InitAndApply(t, terraformOptions)
	// terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	silo := lookupDataSilo(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	discoClassScanConfig := lookupDataSiloDiscoClassScanConfig(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	plugins := lookupDataSiloPlugin(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	return silo, discoClassScanConfig, plugins
}

func TestCanUseSeparateDiscoClassScanConfigResource(t *testing.T) {
	options := prepareDiscoClassScanConfigOptions(t, map[string]interface{}{
		"disco_class_scan_config_vars": map[string]interface{}{
			"enabled":                    true,
			"type":                       "FULL_SCAN",
			"schedule_frequency_minutes": 120,
			// Schedule far in the future so that the test works for a long time
			"schedule_start_at": "2122-09-06T17:51:13.000Z",
		},
	})
	defer terraform.Destroy(t, options)
	silo, discoClassScanConfig, plugins := deployDiscoClassScanConfig(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.NotEmpty(t, terraform.Output(t, options, "awsExternalId"))
	assert.Equal(t, types.DiscoClassScanType("FULL_SCAN"), discoClassScanConfig.Type)
	assert.Equal(t, graphql.Boolean(true), discoClassScanConfig.Enabled)
	assert.Equal(t, graphql.Int(120*1000*60), discoClassScanConfig.ScheduleFrequency) // API returns milliseconds
	assert.Equal(t, graphql.String("2122-09-06T17:51:13.000Z"), discoClassScanConfig.ScheduleStartAt)
	assert.Len(t, plugins, 2)
	hasSchemaDiscovery := false
	hasContentClassification := false
	for _, plugin := range plugins {
		assert.True(t, bool(plugin.Enabled))
		if plugin.Type == "SCHEMA_DISCOVERY" {
			hasSchemaDiscovery = true
		}
		if plugin.Type == "CONTENT_CLASSIFICATION" {
			hasContentClassification = true
		}
		assert.NotEmpty(t, plugin.ID)
	}
	assert.True(t, hasSchemaDiscovery)
	assert.True(t, hasContentClassification)
}
