package transcend

import (
	"context"
	"testing"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/gruntwork-io/terratest/modules/terraform"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
)

func lookupEnricher(t *testing.T, id string) types.Enricher {
	client := getTestClient()

	var query struct {
		Enricher types.Enricher `graphql:"enricher(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": graphql.ID(id),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("ApiKey"))
	assert.Nil(t, err)

	return query.Enricher
}

func prepareEnricherOptions(t *testing.T, vars map[string]interface{}) *terraform.Options {
	defaultVars := map[string]interface{}{"title": t.Name()}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/enricher",
		Vars:         defaultVars,
	})
	return terraformOptions
}

func deployEnricher(t *testing.T, terraformOptions *terraform.Options) types.Enricher {
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "enricherId"))
	enricher := lookupEnricher(t, terraform.Output(t, terraformOptions, "enricherId"))
	return enricher
}

func TestCanCreateAndDestroyEnricher(t *testing.T) {
	options := prepareEnricherOptions(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	enricher := deployEnricher(t, options)
	assert.Equal(t, graphql.String(t.Name()), enricher.Title)
}
