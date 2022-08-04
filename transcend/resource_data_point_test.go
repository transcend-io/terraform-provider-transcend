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

func lookupDataPoint(t *testing.T, id string) types.DataPoint {
	client := NewClient("https://api.dev.trancsend.com/graphql", os.Getenv("TRANSCEND_KEY"))

	var query struct {
		DataPoints struct {
			Nodes []types.DataPoint
		} `graphql:"dataPoints(filterBy: { ids: [$id] })"`
	}
	vars := map[string]interface{}{
		"id": graphql.ID(id),
	}

	err := client.graphql.Query(context.Background(), &query, vars)
	assert.Nil(t, err)
	assert.NotEmpty(t, query.DataPoints)

	return query.DataPoints.Nodes[0]
}

func deployDataPoint(t *testing.T, vars map[string]interface{}) (types.DataPoint, *terraform.Options) {
	defaultVars := map[string]interface{}{"name": t.Name(), "title": t.Name()}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/data_point",
		Vars:         defaultVars,
	})
	terraform.InitAndApply(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "dataPointId"))
	dataPoint := lookupDataPoint(t, terraform.Output(t, terraformOptions, "dataPointId"))
	return dataPoint, terraformOptions
}

func TestCanCreateAndDestroyDataPoint(t *testing.T) {
	dataPoint, options := deployDataPoint(t, map[string]interface{}{})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String(t.Name()), dataPoint.Name)
	assert.Equal(t, graphql.String(t.Name()), dataPoint.Title.DefaultMessage)
}
