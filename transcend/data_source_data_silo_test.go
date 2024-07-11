package transcend

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestCanLookupDataSilo(t *testing.T) {
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
	assert.Equal(t, terraform.Output(t, options, "dataSiloTitle"), t.Name())
}
