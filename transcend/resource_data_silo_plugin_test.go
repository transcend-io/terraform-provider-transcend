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

func deployPlugin(t *testing.T, vars map[string]interface{}) (types.Plugin, *terraform.Options) {
	defaultVars := map[string]interface{}{}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/data_silo_plugin",
		Vars:         defaultVars,
	})
	terraform.InitAndApply(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "gradlePluginDataSiloId"))
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "gradlePluginType"))
	plugin := lookupPlugin(t, terraform.Output(t, terraformOptions, "gradlePluginDataSiloId"), terraform.Output(t, terraformOptions, "gradlePluginType"))
	return plugin, terraformOptions
}

func TestCanCreateAndDestroyPlugin(t *testing.T) {
	plugin, options := deployPlugin(t, map[string]interface{}{"enabled": true})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.Boolean(true), plugin.Enabled)
}

// func TestCanChangeApiKeyTitle(t *testing.T) {
// 	key, options := deployApiKey(t, map[string]interface{}{"enabled": true})
// 	defer terraform.Destroy(t, options)
// 	assert.Equal(t, graphql.String(t.Name()), key.Title)
// 	originalKeyId := key.ID

// 	key, _ = deployApiKey(t, map[string]interface{}{"title": t.Name() + "_2"})
// 	assert.Equal(t, graphql.String(t.Name()+"_2"), key.Title)

// 	// Ensure that a new API key was created
// 	assert.NotEqual(t, originalKeyId, key.ID)
// }

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
