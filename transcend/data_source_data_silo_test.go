package transcend

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestCanLookupDataSilo(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	vars := map[string]interface{}{
		"title": t.Name(),
	}
	options := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/data_silo_data_source",
		Vars:         vars,
	})
	defer terraform.Destroy(t, options)

	terraform.InitAndApplyAndIdempotent(t, options)
	assert.NotEmpty(t, terraform.Output(t, options, "dataSiloId"))
	assert.NotEmpty(t, terraform.Output(t, options, "dataSiloLink"))
	assert.NotEmpty(t, terraform.Output(t, options, "dataSiloDescription"))
	assert.Equal(t, terraform.Output(t, options, "dataSiloTitle"), t.Name())

	assert.NotEmpty(t, terraform.Output(t, options, "dataSiloOwners"))
	assert.Contains(t, terraform.Output(t, options, "dataSiloOwners"), "david@transcend.io")
}
