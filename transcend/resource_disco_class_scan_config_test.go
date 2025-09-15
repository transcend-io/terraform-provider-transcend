package transcend

import (
	"testing"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func prepareDiscoClassScanConfigOptions(t *testing.T, vars map[string]interface{}) *terraform.Options {
	defaultVars := map[string]interface{}{"title": t.Name()}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/disco_class_scan_config",
		Vars:         defaultVars,
	})
	return terraformOptions
}

func deployDiscoClassScanConfig(t *testing.T, terraformOptions *terraform.Options) (types.DataSilo, types.DiscoClassScanConfig) {
	// TODO: Use the Idempotent version eventually
	terraform.InitAndApply(t, terraformOptions)
	// terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	silo := lookupDataSilo(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	discoClassScanConfig := lookupDataSiloDiscoClassScanConfig(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	return silo, discoClassScanConfig
}
