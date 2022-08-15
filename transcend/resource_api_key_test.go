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

func lookupApiKey(t *testing.T, id string) types.APIKey {
	client := NewClient("https://api.dev.trancsend.com/graphql", os.Getenv("TRANSCEND_KEY"))

	var query struct {
		APIKey types.APIKey `graphql:"apiKey(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": graphql.ID(id),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("ApiKey"))
	assert.Nil(t, err)

	return query.APIKey
}

func deployApiKey(t *testing.T, vars map[string]interface{}) (types.APIKey, *terraform.Options) {
	defaultVars := map[string]interface{}{"title": t.Name()}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/api_key",
		Vars:         defaultVars,
	})
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "apiKeyId"))
	key := lookupApiKey(t, terraform.Output(t, terraformOptions, "apiKeyId"))
	return key, terraformOptions
}

func TestCanCreateAndDestroyAPIKey(t *testing.T) {
	key, options := deployApiKey(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String(t.Name()), key.Title)
}

func TestCanChangeApiKeyTitle(t *testing.T) {
	key, options := deployApiKey(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String(t.Name()), key.Title)
	originalKeyId := key.ID

	key, _ = deployApiKey(t, map[string]interface{}{"title": t.Name() + "_2"})
	assert.Equal(t, graphql.String(t.Name()+"_2"), key.Title)

	// Ensure that a new API key was created
	assert.NotEqual(t, originalKeyId, key.ID)
}

func TestCanChangeScopes(t *testing.T) {
	key, options := deployApiKey(t, map[string]interface{}{"scopes": []string{"connectDataSilos"}})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String("connectDataSilos"), key.Scopes[0].Name)

	key, _ = deployApiKey(t, map[string]interface{}{"scopes": []string{"makeDataSubjectRequest"}})
	assert.Equal(t, graphql.String("makeDataSubjectRequest"), key.Scopes[0].Name)
}

func TestCanChangeDataSilos(t *testing.T) {
	key, options := deployApiKey(t, map[string]interface{}{"data_silo_type": "amazonS3"})
	defer terraform.Destroy(t, options)
	originalSiloId := terraform.Output(t, options, "dataSiloId")
	assert.Equal(t, graphql.String(originalSiloId), key.DataSilos[0].ID)

	key, options = deployApiKey(t, map[string]interface{}{"data_silo_type": "asana"})
	newSiloId := terraform.Output(t, options, "dataSiloId")
	assert.Equal(t, graphql.String(newSiloId), key.DataSilos[0].ID)

	// Ensure that the data silo was recreated so that the API key would have to have been updated
	assert.NotEqual(t, originalSiloId, newSiloId)
}
