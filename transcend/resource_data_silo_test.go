package transcend

import (
	"context"
	"os"
	"testing"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/gruntwork-io/terratest/modules/terraform"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
)

func lookupDataSiloPlugin(t *testing.T, id string) []types.Plugin {
	client := getTestClient()

	var query struct {
		Plugins struct {
			Plugins []types.Plugin
		} `graphql:"plugins(filterBy: { dataSiloId: $dataSiloId })"`
	}
	vars := map[string]interface{}{
		"dataSiloId": graphql.String(id),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("Plugins"))
	assert.Nil(t, err)

	return query.Plugins.Plugins
}

func lookupDataSilo(t *testing.T, id string) types.DataSilo {
	client := getTestClient()

	var query struct {
		DataSilo types.DataSilo `graphql:"dataSilo(id: $id)"`
	}
	vars := map[string]interface{}{
		"id": graphql.String(id),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("DataSilo"))
	assert.Nil(t, err)

	return query.DataSilo
}

func prepareDataSiloOptions(t *testing.T, vars map[string]interface{}) *terraform.Options {
	defaultVars := map[string]interface{}{"title": t.Name()}
	for k, v := range vars {
		defaultVars[k] = v
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/tests/data_silo",
		Vars:         defaultVars,
	})
	return terraformOptions
}

func deployDataSilo(t *testing.T, terraformOptions *terraform.Options) (types.DataSilo, []types.Plugin) {
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	silo := lookupDataSilo(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	plugin := lookupDataSiloPlugin(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	return silo, plugin
}

func TestCanCreateAndDestroyDataSilo(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.NotEmpty(t, terraform.Output(t, options, "awsExternalId"))
}

func TestCanConnectAwsDataSilo(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{"skip_connecting": false})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.NotEmpty(t, terraform.Output(t, options, "awsExternalId"))
	assert.Equal(t, types.DataSiloConnectionState("CONNECTED"), silo.ConnectionState)
}

func TestCanConnectDatadogDataSilo(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{
		"skip_connecting": false,
		"type":            "datadog",
		"secret_context": []map[string]interface{}{
			{
				"name":  "apiKey",
				"value": os.Getenv("DD_API_KEY"),
			},
			{
				"name":  "applicationKey",
				"value": os.Getenv("DD_APP_KEY"),
			},
			{
				"name":  "queryTemplate",
				"value": "service:programmatic-remote-seeding AND @email:{{identifier}}",
			},
		},
	})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.Equal(t, types.DataSiloConnectionState("CONNECTED"), silo.ConnectionState)
}

func TestCanConnectPlugin(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{
		"skip_connecting": false,
		"plugin_config": []map[string]interface{}{
			{
				"enabled":                    true,
				"type":                       "DATA_SILO_DISCOVERY",
				"schedule_frequency_minutes": 120,
				// Schedule far in the future so that the test works for a long time
				"schedule_start_at": "2122-09-06T17:51:13+0000",
				"schedule_now":      false,
			},
		},
	})
	defer terraform.Destroy(t, options)
	silo, plugins := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.Equal(t, types.DataSiloConnectionState("CONNECTED"), silo.ConnectionState)
	assert.Len(t, plugins, 1)
	assert.True(t, bool(plugins[0].Enabled))
	assert.NotEmpty(t, plugins[0].ID)
}

func TestCanChangeTitle(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)

	silo, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"title": t.Name() + "_2"}))
	assert.Equal(t, graphql.String(t.Name()+"_2"), silo.Title)
}

func TestCanChangeDescription(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{"description": t.Name()})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)

	silo, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"description": t.Name() + "_2"}))
	assert.Equal(t, graphql.String(t.Name()+"_2"), silo.Description)
}

func TestCanChangeUrl(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{"url": "https://some.webhook", "type": "server"})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("https://some.webhook"), silo.URL)

	silo, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"url": "https://some.other.webhook", "type": "server"}))
	assert.Equal(t, graphql.String("https://some.other.webhook"), silo.URL)
}

func TestCanChangeNotifyEmailAddress(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{"notify_email_address": "david@transcend.io"})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("david@transcend.io"), silo.NotifyEmailAddress)

	silo, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"notify_email_address": "mike@transcend.io"}))
	assert.Equal(t, graphql.String("mike@transcend.io"), silo.NotifyEmailAddress)
}

func TestCanChangeIsLive(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{"is_live": false})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.Boolean(false), silo.IsLive)

	silo, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"is_live": true}))
	assert.Equal(t, graphql.Boolean(true), silo.IsLive)

	silo, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"is_live": false}))
	assert.Equal(t, graphql.Boolean(false), silo.IsLive)
}

func TestCanChangeOwners(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{"owner_emails": []string{"david@transcend.io"}})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("david@transcend.io"), silo.Owners[0].Email)

	silo, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"owner_emails": []string{"mike@transcend.io"}}))
	assert.Equal(t, graphql.String("mike@transcend.io"), silo.Owners[0].Email)
}

func TestCanChangeHeaders(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{"headers": []map[string]interface{}{
		{
			"name":      "someHeader",
			"value":     "someHeaderValue",
			"is_secret": "false",
		},
	}})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("someHeader"), silo.Headers[0].Name)
	assert.Equal(t, graphql.String("someHeaderValue"), silo.Headers[0].Value)

	silo, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"headers": []map[string]interface{}{
		{
			"name":      "someOtherHeader",
			"value":     "someOtherHeaderValue",
			"is_secret": "false",
		},
	}}))
	assert.Equal(t, graphql.String("someOtherHeader"), silo.Headers[0].Name)
	assert.Equal(t, graphql.String("someOtherHeaderValue"), silo.Headers[0].Value)
}

func TestCanCreatePromptAPersonSilo(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{
		"type":       "promptAPerson",
		"outer_type": "coupa",
	})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("coupa"), silo.OuterType)
	assert.Equal(t, graphql.String("promptAPerson"), silo.Type)
	assert.Equal(t, graphql.Boolean(true), silo.Catalog.HasAvcFunctionality)
	assert.Equal(t, graphql.String("dpo@coupa.com"), silo.NotifyEmailAddress)
}

func TestCanSetPromptAPersonNotifyEmailAddress(t *testing.T) {
	options := prepareDataSiloOptions(t, map[string]interface{}{
		"type":                 "promptAPerson",
		"notify_email_address": "not.real.email@transcend.io",
	})
	defer terraform.Destroy(t, options)
	silo, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("promptAPerson"), silo.Type)
	assert.Equal(t, graphql.Boolean(true), silo.Catalog.HasAvcFunctionality)
	assert.Equal(t, graphql.String("not.real.email@transcend.io"), silo.NotifyEmailAddress)
	assert.Empty(t, silo.OuterType)
}
