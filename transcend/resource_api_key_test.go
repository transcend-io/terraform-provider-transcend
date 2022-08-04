package transcend

import (
	"context"
	"os"
	"testing"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/shurcooL/graphql"
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

	err := client.graphql.Query(context.Background(), &query, vars)
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
	terraform.InitAndApply(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "apiKeyId"))
	key := lookupApiKey(t, terraform.Output(t, terraformOptions, "apiKeyId"))
	return key, terraformOptions
}

func TestCanCreateAndDestroyAPIKey(t *testing.T) {
	key, options := deployApiKey(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String(t.Name()), key.Title)
}
