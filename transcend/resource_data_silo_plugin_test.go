package transcend

import (
	"context"
	"os"
	"testing"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/gruntwork-io/terratest/modules/terraform"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
)

func lookupPlugin(t *testing.T, dataSiloId string, typ string) types.Plugin {
	client := NewClient("https://api.dev.trancsend.com/graphql", os.Getenv("TRANSCEND_KEY"))

	var query struct {
		Plugins struct {
			Plugins []types.Plugin
		} `graphql:"plugins(filterBy: $filterBy)"`
	}

	vars := map[string]interface{}{
		"filterBy": types.PluginsFiltersInput{
			DataSiloID: graphql.ID(dataSiloId),
			Type:       types.PluginType(typ),
		},
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("Plugins"))
	assert.Nil(t, err)

	return query.Plugins.Plugins[0]
}

func preparePluginOptions(t *testing.T, vars map[string]interface{}) *terraform.Options {
	defaultVars := map[string]interface{}{}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/data_silo_plugin",
		Vars:         defaultVars,
	})
	return terraformOptions
}

func deployPlugin(t *testing.T, options *terraform.Options) types.Plugin {
	terraform.InitAndApply(t, options)
	assert.NotEmpty(t, terraform.Output(t, options, "gradlePluginDataSiloId"))
	assert.NotEmpty(t, terraform.Output(t, options, "gradlePluginType"))
	plugin := lookupPlugin(t, terraform.Output(t, options, "gradlePluginDataSiloId"), terraform.Output(t, options, "gradlePluginType"))
	return plugin
}

func TestCanCreateAndDestroyPlugin(t *testing.T) {
	options := preparePluginOptions(t, map[string]interface{}{"enabled": true})
	plugin := deployPlugin(t, options)
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.Boolean(true), plugin.Enabled)
}

func TestCanChangeEnabled(t *testing.T) {
	options := preparePluginOptions(t, map[string]interface{}{"enabled": true})
	defer terraform.Destroy(t, options)

	plugin := deployPlugin(t, options)
	assert.Equal(t, graphql.Boolean(true), plugin.Enabled)

	plugin = deployPlugin(t, preparePluginOptions(t, map[string]interface{}{"enabled": false}))
	assert.Equal(t, graphql.Boolean(false), plugin.Enabled)
}

func TestCanChangeFrequency(t *testing.T) {
	options := preparePluginOptions(t, map[string]interface{}{"schedule_frequency": "2000"})
	defer terraform.Destroy(t, options)

	plugin := deployPlugin(t, options)
	assert.Equal(t, graphql.String("2000"), plugin.ScheduleFrequency)

	plugin = deployPlugin(t, preparePluginOptions(t, map[string]interface{}{"schedule_frequency": "3000"}))
	assert.Equal(t, graphql.String("3000"), plugin.ScheduleFrequency)
}

// func TestCanChangeScopes(t *testing.T) {
// 	key, options := deployApiKey(t, map[string]interface{}{"scopes": []string{"connectDataSilos"}})
// 	defer terraform.Destroy(t, options)
// 	assert.Equal(t, graphql.String("connectDataSilos"), key.Scopes[0].Name)

// 	key, _ = deployApiKey(t, map[string]interface{}{"scopes": []string{"makeDataSubjectRequest"}})
// 	assert.Equal(t, graphql.String("makeDataSubjectRequest"), key.Scopes[0].Name)
// }

// func TestCanChangeDataSilos(t *testing.T) {
// 	key, options := deployApiKey(t, map[string]interface{}{"data_silo_type": "amazonS3"})
// 	defer terraform.Destroy(t, options)
// 	originalSiloId := terraform.Output(t, options, "dataSiloId")
// 	assert.Equal(t, graphql.String(originalSiloId), key.DataSilos[0].ID)

// 	key, options = deployApiKey(t, map[string]interface{}{"data_silo_type": "asana"})
// 	newSiloId := terraform.Output(t, options, "dataSiloId")
// 	assert.Equal(t, graphql.String(newSiloId), key.DataSilos[0].ID)

// 	// Ensure that the data silo was recreated so that the API key would have to have been updated
// 	assert.NotEqual(t, originalSiloId, newSiloId)
// }
