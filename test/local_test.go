package test

import (
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformHelloWorldExample(t *testing.T) {
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/test/",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "amazon_title")
	assert.Equal(t, "Amazon", output)

	output = terraform.Output(t, terraformOptions, "amazon_description")
	assert.Equal(t, "This is a test", output)

	output = terraform.Output(t, terraformOptions, "amazon_link")
	http_helper.HttpGetWithRetry(t, output, nil, 200, "", 30, 5*time.Second)
}
