package transcend

import (
	"context"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
	"github.com/transcend-io/terraform-provider-transcend/transcend/types"
)

func lookupIdentifier(t *testing.T, id string) types.Identifier {
	client := getTestClient()

	var query struct {
		Identifier types.Identifier `graphql:"identifier(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": graphql.ID(id),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("Identifier"))
	assert.Nil(t, err)

	return query.Identifier
}

func prepareIdentifierOptions(t *testing.T, vars map[string]interface{}) *terraform.Options {
	defaultVars := map[string]interface{}{"name": t.Name()}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/identifier",
		Vars:         defaultVars,
	})
	return terraformOptions
}

func deployIdentifier(t *testing.T, terraformOptions *terraform.Options) types.Identifier {
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "identifierId"))
	identifier := lookupIdentifier(t, terraform.Output(t, terraformOptions, "identifierId"))
	return identifier
}

func TestCanCreateAndDestroyIdentifier(t *testing.T) {
	options := prepareIdentifierOptions(t, map[string]interface{}{"name": t.Name()})
	defer terraform.Destroy(t, options)
	identifier := deployIdentifier(t, options)
	assert.Equal(t, graphql.String(t.Name()), identifier.Name)
}
