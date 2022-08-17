package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

type Plugin struct {
	ID                graphql.String  `json:"id"`
	Type              PluginType      `json:"type"`
	Enabled           graphql.Boolean `json:"enabled"`
	ScheduledAt       graphql.String  `json:"scheduledAt"`
	LastRunAt         graphql.String  `json:"lastRunAt"`
	LastEnabledAt     graphql.String  `json:"lastEnabledAt"`
	ScheduleStartAt   graphql.String  `json:"scheduleStartAt"`
	ScheduleFrequency graphql.String  `json:"scheduleFrequency"`
	Error             graphql.String  `json:"error"`
	DataSilo          struct {
		ID graphql.String `json:"id"`
	} `json:"DataSilo"`
}

type DataSiloPlugin struct {
	Plugin Plugin `json:"plugin"`
}

type PluginsFiltersInput struct {
	DataSiloID graphql.ID      `json:"dataSiloId"`
	Enabled    graphql.Boolean `json:"enabled,omitempty"`
	Type       PluginType      `json:"type"`
}

type UpdatePluginInput struct {
	DataSiloID        graphql.ID      `json:"dataSiloId"`
	PluginID          graphql.ID      `json:"pluginId"`
	Enabled           graphql.Boolean `json:"enabled,omitempty"`
	ScheduleFrequency graphql.String  `json:"scheduleFrequency"`
	ScheduleStartAt   graphql.String  `json:"scheduleStartAt"`
	ScheduleNow       graphql.Boolean `json:"scheduleNow"`
}

func ReadDataSiloPluginIntoState(d *schema.ResourceData, plugin Plugin) {
	d.Set("data_silo_id", plugin.DataSilo.ID)
	d.Set("enabled", plugin.Enabled)
	d.Set("scheduled_at", plugin.ScheduledAt)
	d.Set("last_run_at", plugin.LastRunAt)
	d.Set("last_enabled_at", plugin.LastEnabledAt)
	d.Set("schedule_start_at", plugin.ScheduleStartAt)
	d.Set("schedule_frequency", plugin.ScheduleFrequency)
	d.Set("error", plugin.Error)
}

func MakePluginsFiltersInput(d *schema.ResourceData) PluginsFiltersInput {
	return PluginsFiltersInput{
		DataSiloID: graphql.ID(d.Get("data_silo_id").(string)),
		Type:       PluginType(d.Get("type").(string)),
	}
}

func MakeUpdatePluginInput(d *schema.ResourceData, date string) UpdatePluginInput {
	return UpdatePluginInput{
		DataSiloID:        graphql.ID(d.Get("data_silo_id").(string)),
		PluginID:          graphql.ID(d.Get("id").(string)),
		Enabled:           graphql.Boolean(d.Get("enabled").(bool)),
		ScheduleFrequency: graphql.String(d.Get("schedule_frequency").(string)),
		ScheduleStartAt:   graphql.String(date),
		ScheduleNow:       graphql.Boolean(d.Get("schedule_now").(bool)),
	}
}

func PluginsReadQuery(client graphql.Client, d *schema.ResourceData) (Plugin, string) {
	var query struct {
		Plugins struct {
			Plugins []Plugin
		} `graphql:"plugins(filterBy: $filterBy)"`
	}

	vars := map[string]interface{}{
		"filterBy": MakePluginsFiltersInput(d),
	}

	err := client.Query(context.Background(), &query, vars, graphql.OperationName("Plugins"))
	if err != nil {
		return Plugin{}, "Error when reading plugin: " + err.Error()
	}

	if len(query.Plugins.Plugins) == 0 {
		return Plugin{}, "Did not able to find plugin"
	}

	if len(query.Plugins.Plugins) > 1 {
		return Plugin{}, "Found multiple plugins"
	}

	return query.Plugins.Plugins[0], ""
}

func PluginsUpdateQuery(client graphql.Client, d *schema.ResourceData) string {
	var updateMutation struct {
		UpdateDataSiloPlugin struct {
			Plugin Plugin
		} `graphql:"updateDataSiloPlugin(input: $input)"`
	}

	date := d.Get("schedule_start_at").(string)

	if date == "" {
		plugin, err := PluginsReadQuery(client, d)
		if err != "" {
			return err
		}
		date = string(plugin.ScheduleStartAt)
	}
	updateVars := map[string]interface{}{
		"input": MakeUpdatePluginInput(d, date),
	}

	err := client.Mutate(context.Background(), &updateMutation, updateVars, graphql.OperationName("UpdateDataSiloPlugin"))
	if err != nil {
		return "Error when updating plugin: " + err.Error()
	}

	return ""
}
