package transcend

import (
	"context"
	"os"
	"strconv"
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

func TestCanChangeDataPointTitle(t *testing.T) {
	dataPoint, options := deployDataPoint(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String(t.Name()), dataPoint.Title.DefaultMessage)

	dataPoint, _ = deployDataPoint(t, map[string]interface{}{"title": t.Name() + "_2"})
	assert.Equal(t, graphql.String(t.Name()+"_2"), dataPoint.Title.DefaultMessage)
}

func TestCanChangeDataPointName(t *testing.T) {
	dataPoint, options := deployDataPoint(t, map[string]interface{}{"name": t.Name()})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String(t.Name()), dataPoint.Name)

	dataPoint, _ = deployDataPoint(t, map[string]interface{}{"name": t.Name() + "_2"})
	assert.Equal(t, graphql.String(t.Name()+"_2"), dataPoint.Name)
}

func TestCanChangeDataPointDescription(t *testing.T) {
	dataPoint, options := deployDataPoint(t, map[string]interface{}{"description": t.Name()})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String(t.Name()), dataPoint.Description.DefaultMessage)

	dataPoint, _ = deployDataPoint(t, map[string]interface{}{"description": t.Name() + "_2"})
	assert.Equal(t, graphql.String(t.Name()+"_2"), dataPoint.Description.DefaultMessage)
}

func TestCanChangeDataPointSilo(t *testing.T) {
	dataPoint, options := deployDataPoint(t, map[string]interface{}{"data_silo_type": "server"})
	defer terraform.Destroy(t, options)
	originalSiloId := terraform.Output(t, options, "dataSiloId")
	assert.Equal(t, graphql.String(originalSiloId), dataPoint.DataSilo.ID)

	dataPoint, options = deployDataPoint(t, map[string]interface{}{"data_silo_type": "promptAPerson"})
	newSiloId := terraform.Output(t, options, "dataSiloId")
	assert.Equal(t, graphql.String(newSiloId), dataPoint.DataSilo.ID)

	// Ensure that the data silo was recreated so that the API key would have to have been updated
	assert.NotEqual(t, originalSiloId, newSiloId)
}

func TestCanCreateDataPointWithSubDataPoints(t *testing.T) {
	_, options := deployDataPoint(t, map[string]interface{}{
		"properties": []map[string]interface{}{
			{"name": "subDataPoint1"},
			{"name": "subDataPoint2"},
			{"name": "subDataPoint3"},
			{"name": "subDataPoint4"},
		},
	})
	defer terraform.Destroy(t, options)
	properties := terraform.OutputListOfObjects(t, options, "properties")
	assert.Len(t, properties, 4)
}

func TestCanChangeSubDataPoints(t *testing.T) {
	_, options := deployDataPoint(t, map[string]interface{}{
		"properties": []map[string]interface{}{
			{"name": "subDataPoint1"},
			{"name": "subDataPoint2"},
			{"name": "subDataPoint3"},
			{"name": "subDataPoint4"},
		},
	})
	defer terraform.Destroy(t, options)
	properties := terraform.OutputListOfObjects(t, options, "properties")
	assert.Len(t, properties, 4)

	_, options = deployDataPoint(t, map[string]interface{}{
		"properties": []map[string]interface{}{
			{"name": "onlySubDataPoint"},
		},
	})
	properties = terraform.OutputListOfObjects(t, options, "properties")
	assert.Len(t, properties, 1)

	_, options = deployDataPoint(t, map[string]interface{}{})
	properties = terraform.OutputListOfObjects(t, options, "properties")
	assert.Len(t, properties, 0)
}

func TestCanPaginateSubDataPoints(t *testing.T) {
	properties := make([]map[string]interface{}, 251)
	for i := 0; i < 251; i++ {
		properties[i] = map[string]interface{}{
			"name": "subDataPoint" + strconv.Itoa(i),
		}
	}

	_, options := deployDataPoint(t, map[string]interface{}{
		"properties": properties,
	})
	defer terraform.Destroy(t, options)
	propertiesOutput := terraform.OutputListOfObjects(t, options, "properties")
	assert.Len(t, propertiesOutput, 251)
}
