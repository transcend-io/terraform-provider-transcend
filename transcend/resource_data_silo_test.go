package transcend

import (
	"context"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

func lookupDataSilo(t *testing.T, id string) DataSilo {
	client := NewClient("https://api.dev.trancsend.com/graphql", os.Getenv("TRANSCEND_KEY"))

	var query struct {
		DataSilo DataSilo `graphql:"dataSilo(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": graphql.String(id),
	}

	err := client.graphql.Query(context.Background(), &query, vars)
	assert.Nil(t, err)

	return query.DataSilo
}

func deployDataSilo(t *testing.T, vars map[string]interface{}) (DataSilo, *terraform.Options) {
	defaultVars := map[string]interface{}{"title": t.Name()}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/data_silo",
		Vars:         defaultVars,
	})
	terraform.InitAndApply(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	silo := lookupDataSilo(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	return silo, terraformOptions
}

func TestCanCreateAndDestroyDataSilo(t *testing.T) {
	silo, options := deployDataSilo(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.NotEmpty(t, terraform.Output(t, options, "awsExternalId"))
}

func TestCanChangeTitle(t *testing.T) {
	silo, options := deployDataSilo(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)

	silo, _ = deployDataSilo(t, map[string]interface{}{"title": t.Name() + "_2"})
	assert.Equal(t, graphql.String(t.Name()+"_2"), silo.Title)
}

func TestCanChangeDescription(t *testing.T) {
	silo, options := deployDataSilo(t, map[string]interface{}{"description": t.Name()})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)

	silo, _ = deployDataSilo(t, map[string]interface{}{"description": t.Name() + "_2"})
	assert.Equal(t, graphql.String(t.Name()+"_2"), silo.Description)
}

func TestCanChangeUrl(t *testing.T) {
	silo, options := deployDataSilo(t, map[string]interface{}{"url": "https://some.webhook", "type": "server"})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String("https://some.webhook"), silo.URL)

	silo, _ = deployDataSilo(t, map[string]interface{}{"url": "https://some.other.webhook", "type": "server"})
	assert.Equal(t, graphql.String("https://some.other.webhook"), silo.URL)
}

func TestCanChangeNotifyEmailAddress(t *testing.T) {
	silo, options := deployDataSilo(t, map[string]interface{}{"notify_email_address": "david@transcend.io"})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String("david@transcend.io"), silo.NotifyEmailAddress)

	silo, _ = deployDataSilo(t, map[string]interface{}{"notify_email_address": "mike@transcend.io"})
	assert.Equal(t, graphql.String("mike@transcend.io"), silo.NotifyEmailAddress)
}

func TestCanChangeIsLive(t *testing.T) {
	silo, options := deployDataSilo(t, map[string]interface{}{"is_live": false})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.Boolean(false), silo.IsLive)

	silo, _ = deployDataSilo(t, map[string]interface{}{"is_live": true})
	assert.Equal(t, graphql.Boolean(true), silo.IsLive)

	silo, _ = deployDataSilo(t, map[string]interface{}{"is_live": false})
	assert.Equal(t, graphql.Boolean(false), silo.IsLive)
}

func TestCanChangeOwners(t *testing.T) {
	silo, options := deployDataSilo(t, map[string]interface{}{"owner_emails": []string{"david@transcend.io"}})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String("david@transcend.io"), silo.Owners[0].Email)

	silo, _ = deployDataSilo(t, map[string]interface{}{"owner_emails": []string{"mike@transcend.io"}})
	assert.Equal(t, graphql.String("mike@transcend.io"), silo.Owners[0].Email)
}

func TestCanChangeHeaders(t *testing.T) {
	silo, options := deployDataSilo(t, map[string]interface{}{"headers": []map[string]interface{}{
		{
			"name":      "someHeader",
			"value":     "someHeaderValue",
			"is_secret": "false",
		},
	}})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String("someHeader"), silo.Headers[0].Name)
	assert.Equal(t, graphql.String("someHeaderValue"), silo.Headers[0].Value)

	silo, _ = deployDataSilo(t, map[string]interface{}{"headers": []map[string]interface{}{
		{
			"name":      "someOtherHeader",
			"value":     "someOtherHeaderValue",
			"is_secret": "false",
		},
	}})
	assert.Equal(t, graphql.String("someOtherHeader"), silo.Headers[0].Name)
	assert.Equal(t, graphql.String("someOtherHeaderValue"), silo.Headers[0].Value)
}

func TestCanCreatePromptAPersonSilo(t *testing.T) {
	silo, options := deployDataSilo(t, map[string]interface{}{
		"type":       "promptAPerson",
		"outer_type": "coupa",
	})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String("coupa"), silo.OuterType)
	assert.Equal(t, graphql.String("promptAPerson"), silo.Type)
	assert.Equal(t, graphql.Boolean(true), silo.Catalog.HasAvcFunctionality)
	assert.Equal(t, graphql.String("dpo@coupa.com"), silo.NotifyEmailAddress)
}

func TestCanSetPromptAPersonNotifyEmailAddress(t *testing.T) {
	silo, options := deployDataSilo(t, map[string]interface{}{
		"type":                 "promptAPerson",
		"notify_email_address": "not.real.email@transcend.io",
	})
	defer terraform.Destroy(t, options)
	assert.Equal(t, graphql.String("promptAPerson"), silo.Type)
	assert.Equal(t, graphql.Boolean(true), silo.Catalog.HasAvcFunctionality)
	assert.Equal(t, graphql.String("not.real.email@transcend.io"), silo.NotifyEmailAddress)
	assert.Empty(t, silo.OuterType)
}
