package transcend

import (
	"context"
	"testing"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/gruntwork-io/terratest/modules/terraform"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
)

func lookupPlugin(t *testing.T, dataSiloId string, typ string) types.Plugin {
	client := getTestClient()

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

func TestCanChangeFrequency(t *testing.T) {
	options := preparePluginOptions(t, map[string]interface{}{"schedule_frequency_minutes": "2000"})
	defer terraform.Destroy(t, options)

	plugin := deployPlugin(t, options)
	assert.Equal(t, graphql.String("2000"), plugin.ScheduleFrequency)

	plugin = deployPlugin(t, preparePluginOptions(t, map[string]interface{}{"schedule_frequency_minutes": "3000"}))
	assert.Equal(t, graphql.String("3000"), plugin.ScheduleFrequency)
}

// func TestCanScheduleStartAt(t *testing.T) {
// 	options := preparePluginOptions(t, map[string]interface{}{"schedule_start_at": "2022-08-16T08:00:00.000Z"})
// 	defer terraform.Destroy(t, options)

// 	plugin := deployPlugin(t, options)
// 	assert.Equal(t, graphql.String("2022-08-16T08:00:00.000Z"), plugin.ScheduleStartAt)

// 	plugin = deployPlugin(t, preparePluginOptions(t, map[string]interface{}{"schedule_start_at": "2022-08-16T09:00:00.000Z"}))
// 	assert.Equal(t, graphql.String("2022-08-16T09:00:00.000Z"), plugin.ScheduleStartAt)
// }
