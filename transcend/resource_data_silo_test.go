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

// Helper to destroy any data silo with a given title before a test runs
func destroyDataSiloByTitle(t *testing.T, title string) {
	client := getTestClient()
	var query struct {
		DataSilos struct {
			Nodes []struct {
				ID    graphql.String `json:"id"`
				Title graphql.String `json:"title"`
			}
		} `graphql:"dataSilos(filterBy: { titles: [$title] })"`
	}
	vars := map[string]interface{}{
		"title": graphql.String(title),
	}
	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("DataSilos"))
	if err != nil {
		t.Logf("Could not query for silos to destroy: %v", err)
		return
	}
	for _, node := range query.DataSilos.Nodes {
		var mutation struct {
			DeleteDataSilos struct {
				Success graphql.Boolean
			} `graphql:"deleteDataSilos(input: { ids: $ids })"`
		}
		ids := []graphql.ID{graphql.ID(node.ID)}
		mvars := map[string]interface{}{
			"ids": ids,
		}
		_ = client.graphql.Mutate(context.Background(), &mutation, mvars, graphql.OperationName("DeleteDataSilos"))
	}
}

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
	if err != nil {
		return []types.Plugin{}
	}

	return query.Plugins.Plugins
}

func lookupDataSiloDiscoClassScanConfig(t *testing.T, id string) types.DiscoClassScanConfig {
	client := getTestClient()

	var query struct {
		DiscoClassScanConfig types.DiscoClassScanConfig `graphql:"discoClassScanConfig(input: { dataSiloId: $dataSiloId })"`
	}
	vars := map[string]interface{}{
		"dataSiloId": graphql.ID(id),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("DiscoClassScanConfig"))
	if err != nil {
		return types.DiscoClassScanConfig{}
	}

	return query.DiscoClassScanConfig
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

func deployDataSilo(t *testing.T, terraformOptions *terraform.Options) (types.DataSilo, []types.Plugin, types.DiscoClassScanConfig) {
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.NotEmpty(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	silo := lookupDataSilo(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	plugin := lookupDataSiloPlugin(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	discoClassScanConfig := lookupDataSiloDiscoClassScanConfig(t, terraform.Output(t, terraformOptions, "dataSiloId"))
	return silo, plugin, discoClassScanConfig
}

func TestCanCreateAndDestroyDataSilo(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.NotEmpty(t, terraform.Output(t, options, "awsExternalId"))
}

func TestCanConnectAwsDataSilo(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"skip_connecting": false})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.NotEmpty(t, terraform.Output(t, options, "awsExternalId"))
	assert.Equal(t, types.DataSiloConnectionState("CONNECTED"), silo.ConnectionState)
}

func TestCanConnectDatadogDataSilo(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
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
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.Equal(t, types.DataSiloConnectionState("CONNECTED"), silo.ConnectionState)
}

func TestCanConnectSchemaDiscoveryAndContentClassificationPlugin(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{
		"skip_connecting": false,
		"disco_class_scan_config_vars": []map[string]interface{}{
			{
				"enabled":                    true,
				"type":                       "FULL_SCAN",
				"schedule_frequency_minutes": 120,
				// Schedule far in the future so that the test works for a long time
				"schedule_start_at": "2122-09-06T17:51:13.000Z",
			},
		},
		"schema_discovery_plugin_config": []map[string]interface{}{
			{
				"enabled":                    true,
				"schedule_frequency_minutes": 120,
				// Schedule far in the future so that the test works for a long time
				"schedule_start_at": "2122-09-06T17:51:13.000Z",
			},
		},
		"content_classification_plugin_config": []map[string]interface{}{
			{
				"enabled":                    true,
				"schedule_frequency_minutes": 120,
				// Schedule far in the future so that the test works for a long time
				"schedule_start_at": "2122-09-06T17:51:13.000Z",
			},
		},
	})
	defer terraform.Destroy(t, options)
	silo, plugins, discoClassScanConfig := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.Equal(t, types.DataSiloConnectionState("CONNECTED"), silo.ConnectionState)
	assert.True(t, bool(discoClassScanConfig.Enabled))
	assert.NotEmpty(t, discoClassScanConfig.ID)
	assert.Equal(t, types.DiscoClassScanType("FULL_SCAN"), discoClassScanConfig.Type)
	assert.Len(t, plugins, 2)
	for _, plugin := range plugins {
		assert.True(t, bool(plugin.Enabled))
		assert.NotEmpty(t, plugin.ID)
	}
}

func TestCanChangeTitle(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"title": t.Name()})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)

	silo, _, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"title": t.Name() + "_2"}))
	assert.Equal(t, graphql.String(t.Name()+"_2"), silo.Title)
}

func TestCanChangeDescription(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"description": t.Name()})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)

	silo, _, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"description": t.Name() + "_2"}))
	assert.Equal(t, graphql.String(t.Name()+"_2"), silo.Description)
}

func TestCanChangeSaasContext(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
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
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.Equal(t, types.DataSiloConnectionState("CONNECTED"), silo.ConnectionState)

	silo, _, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{
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
				"value": "service:a-different-service AND @email:{{identifier}}",
			},
		},
	}))
	assert.Equal(t, types.DataSiloConnectionState("CONNECTED"), silo.ConnectionState)
}

func TestThatChangingSaasContextToInvalidValueDoesNotDestroySilo(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
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
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String(t.Name()), silo.Title)
	assert.Equal(t, types.DataSiloConnectionState("CONNECTED"), silo.ConnectionState)

	updatedOptions := prepareDataSiloOptions(t, map[string]interface{}{
		"skip_connecting": false,
		"type":            "datadog",
		"secret_context": []map[string]interface{}{
			{
				"name":  "apiKey",
				"value": "not-an-api-key-at-least-most-likely-unless-uuid-changes-formats-and-this-string-is-randomly-generated",
			},
			{
				"name":  "applicationKey",
				"value": "applicationKeyMcApplicationKeyFace",
			},
			{
				"name":  "queryTemplate",
				"value": "service:a-different-service AND @email:{{identifier}}",
			},
		},
	})
	terraform.ApplyE(t, updatedOptions)
	assert.NotEmpty(t, terraform.Output(t, updatedOptions, "dataSiloId"))
	updatedSilo := lookupDataSilo(t, terraform.Output(t, updatedOptions, "dataSiloId"))
	assert.Equal(t, types.DataSiloConnectionState("CONNECTED"), updatedSilo.ConnectionState)
}

func TestCanChangeUrl(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"url": "https://some.webhook", "type": "server"})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("https://some.webhook"), silo.URL)

	silo, _, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"url": "https://some.other.webhook", "type": "server"}))
	assert.Equal(t, graphql.String("https://some.other.webhook"), silo.URL)
}

func TestCanChangeNotifyEmailAddress(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"notify_email_address": "david@transcend.io"})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("david@transcend.io"), silo.NotifyEmailAddress)

	silo, _, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"notify_email_address": "mike@transcend.io"}))
	assert.Equal(t, graphql.String("mike@transcend.io"), silo.NotifyEmailAddress)
}

func TestCanChangeIsLive(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"is_live": false})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.Boolean(false), silo.IsLive)

	silo, _, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"is_live": true}))
	assert.Equal(t, graphql.Boolean(true), silo.IsLive)

	silo, _, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"is_live": false}))
	assert.Equal(t, graphql.Boolean(false), silo.IsLive)
}

func TestCanChangeOwners(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"owner_emails": []string{"david@transcend.io"}})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("david@transcend.io"), silo.Owners[0].Email)

	silo, _, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"owner_emails": []string{"mike@transcend.io"}}))
	assert.Equal(t, graphql.String("mike@transcend.io"), silo.Owners[0].Email)
}

func TestCanChangeOwnerTeams(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"owner_teams": []string{"Engineers"}})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("Engineers"), silo.Teams[0].Name)

	silo, _, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"owner_teams": []string{"Legal"}}))
	assert.Equal(t, graphql.String("Legal"), silo.Teams[0].Name)
}

func TestCanChangeHeaders(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"headers": []map[string]interface{}{
		{
			"name":      "someHeader",
			"value":     "someHeaderValue",
			"is_secret": "false",
		},
	}})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("someHeader"), silo.Headers[0].Name)
	assert.Equal(t, graphql.String("someHeaderValue"), silo.Headers[0].Value)

	silo, _, _ = deployDataSilo(t, prepareDataSiloOptions(t, map[string]interface{}{"headers": []map[string]interface{}{
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
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{
		"type":       "promptAPerson",
		"outer_type": "coupa",
	})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("coupa"), silo.OuterType)
	assert.Equal(t, graphql.String("promptAPerson"), silo.Type)
	assert.Equal(t, graphql.Boolean(true), silo.Catalog.HasAvcFunctionality)
	assert.Equal(t, graphql.String("dpo@coupa.com"), silo.NotifyEmailAddress)
}

func TestCanSetPromptAPersonNotifyEmailAddress(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{
		"type":                 "promptAPerson",
		"notify_email_address": "not.real.email@transcend.io",
	})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("promptAPerson"), silo.Type)
	assert.Equal(t, graphql.Boolean(true), silo.Catalog.HasAvcFunctionality)
	assert.Equal(t, graphql.String("not.real.email@transcend.io"), silo.NotifyEmailAddress)
	assert.Empty(t, silo.OuterType)
}

func TestCanAddSombraId(t *testing.T) {
	destroyDataSiloByTitle(t, t.Name())
	options := prepareDataSiloOptions(t, map[string]interface{}{"sombra_id": "c29ba149-b7b4-4ff1-a93f-b24641271ea7"})
	defer terraform.Destroy(t, options)
	silo, _, _ := deployDataSilo(t, options)
	assert.Equal(t, graphql.String("c29ba149-b7b4-4ff1-a93f-b24641271ea7"), silo.SombraId)
}

func TestOwnerEmailsOrderIndependence(t *testing.T) {
	// Clean up any silos with the same titles before running the test
	destroyDataSiloByTitle(t, t.Name()+"_A")
	destroyDataSiloByTitle(t, t.Name()+"_B")
	// Test with owner_emails in one order
	emailsA := []string{"b@example.com", "a@example.com"}
	emailsB := []string{"a@example.com", "b@example.com"}

	optionsA := prepareDataSiloOptions(t, map[string]interface{}{
		"owner_emails": emailsA,
		"title":        t.Name() + "_A",
	})
	optionsB := prepareDataSiloOptions(t, map[string]interface{}{
		"owner_emails": emailsB,
		"title":        t.Name() + "_B",
	})

	defer terraform.Destroy(t, optionsA)
	defer terraform.Destroy(t, optionsB)

	siloA, _, _ := deployDataSilo(t, optionsA)
	siloB, _, _ := deployDataSilo(t, optionsB)

	// Both silos should have the same owner emails in state, regardless of input order
	var outputA, outputB []string
	for _, o := range siloA.Owners {
		outputA = append(outputA, string(o.Email))
	}
	for _, o := range siloB.Owners {
		outputB = append(outputB, string(o.Email))
	}

	assert.ElementsMatch(t, outputA, outputB, "Owner emails should be order-independent in state and not change between applies")
}
