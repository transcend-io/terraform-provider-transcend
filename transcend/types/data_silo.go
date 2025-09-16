package types

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

type DataSiloUpdatableFields struct {
	Title                   graphql.String      `json:"title,omitempty"`
	Description             graphql.String      `json:"description,omitempty"`
	URL                     graphql.String      `json:"url,omitempty"`
	NotifyEmailAddress      graphql.String      `json:"notifyEmailAddress,omitempty"`
	IsLive                  graphql.Boolean     `json:"isLive"`
	OwnerEmails             []graphql.String    `json:"ownerEmails"`
	OwnerTeams              []graphql.String    `json:"teamNames"`
	DataSubjectBlockListIds []graphql.String    `json:"dataSubjectBlockListIds"`
	Headers                 []CustomHeaderInput `json:"headers"`
	SombraId                graphql.String      `json:"sombraId,omitempty"`

	// TODO: Support more fields
	// Identifiers             []graphql.String    `json:"identifiers"`
	// dependedOnDataSiloIds
	// dependedOnDataSiloTitles
	// ownerIds
	// apiKeyId
	// teams
	// teamNames
	// notes
	// dataRetentionNote
	// dataProcessingAgreementLink
	// contactName
	// contactEmail
	// dataProcessingAgreementStatus
	// recommendedForConsent
	// recommendedForPrivacy
	// hasPersonalData
	// deprecationState
}

type CreateDataSilosInput struct {
	Name  graphql.String `json:"name"`
	Title graphql.String `json:"title"`
}

type UpdateDataSiloInput struct {
	Id graphql.ID `json:"id"`
	DataSiloUpdatableFields
}

type PlaintextContextInput struct {
	Name  graphql.String `json:"name"`
	Value graphql.String `json:"value"`
}

type PlaintextInformation struct {
	Path graphql.String `json:"path"`
}

type Catalog struct {
	PlaintextInformation []PlaintextInformation `json:"plaintextInformation"`
	IntegrationConfig    struct {
		ConfiguredBaseHosts struct {
			PROD []graphql.String `graphql:"PROD"`
		} `json:"configuredBaseHosts"`
	} `json:"integrationConfig"`
}

type DataSiloFiltersInput struct {
	Ids          []graphql.ID     `json:"ids,omitempty"`
	DiscoveredBy []graphql.ID     `json:"discoveredBy,omitempty"`
	Type         []graphql.String `json:"type,omitempty"`
	Title        []graphql.String `json:"titles,omitempty"`
}

type DataSilo struct {
	ID         graphql.String `json:"id"`
	Link       graphql.String `json:"link,omitempty"`
	ExternalId graphql.String `json:"externalId,omitempty"`
	Catalog    struct {
		HasAvcFunctionality graphql.Boolean `json:"hasAvcFunctionality"`
	} `json:"catalog"`

	Type               graphql.String  `json:"type"`
	Title              graphql.String  `json:"title"`
	Description        graphql.String  `json:"description,omitempty"`
	URL                graphql.String  `json:"url,omitempty"`
	NotifyEmailAddress graphql.String  `json:"notifyEmailAddress,omitempty"`
	IsLive             graphql.Boolean `json:"isLive"`
	Owners             []struct {
		ID    graphql.String `json:"id"`
		Email graphql.String `json:"email"`
	} `json:"owners"`
	Teams []struct {
		ID   graphql.String `json:"id"`
		Name graphql.String `json:"name"`
	} `json:"teams"`
	SubjectBlocklist []struct {
		ID graphql.String `json:"id"`
	} `json:"subjectBlocklist"`
	Headers          []Header                `json:"headers"`
	OuterType        graphql.String          `json:"outerType"`
	PlaintextContext []PlaintextContextInput `json:"plaintextContext"`
	ConnectionState  DataSiloConnectionState `json:"connectionState"`
	SombraId         graphql.String          `json:"sombraId,omitempty"`

	// Sombra enricher identifier mappings
	EnricherIdentifierMappings []EnricherIdentifierMapping `json:"enricherIdentifierMappings"`
	// TODO: Add support to DataSiloInput first
	// Identifiers        []struct {
	//  Name graphql.String `json:"name"`
	// } `json:"identifiers"`
	// PromptEmailTemplate struct {
	//  ID graphql.String `json:"id,omitempty"`
	// } `json:"promptEmailTemplate,omitempty"`

	// TODO: Look up the schema here
	// Teams   []struct{} `json:"teams"`
	// ApiKeys []struct{} `json:"apiKeys"`
	// DependentDataSilos []struct{} `json:"dependentDataSilos"`
}

type DataSiloBulkPreview struct {
	ID          graphql.String `json:"id"`
	Title       graphql.String `json:"title"`
	Type        graphql.String `json:"type"`
	Link        graphql.String `json:"link"`
	Description graphql.String `json:"description"`
	Owners      []struct {
		ID    graphql.String `json:"id"`
		Email graphql.String `json:"email"`
	} `json:"owners"`
}

type DataSilosPayload struct {
	Nodes []DataSiloBulkPreview `json:"nodes"`
}

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

type UpdatePluginInput struct {
	DataSiloID               graphql.ID      `json:"dataSiloId"`
	PluginID                 graphql.ID      `json:"pluginId"`
	Enabled                  graphql.Boolean `json:"enabled"`
	ScheduleFrequencyMinutes graphql.String  `json:"scheduleFrequency"`
	ScheduleStartAt          graphql.String  `json:"scheduleStartAt"`
	ScheduleNow              graphql.Boolean `json:"scheduleNow"`
}

type DiscoClassScanConfig struct {
	ID                       graphql.String     `json:"id"`
	DataSiloID               graphql.String     `json:"dataSiloId"`
	Type                     DiscoClassScanType `json:"type"`
	Enabled                  graphql.Boolean    `json:"enabled"`
	ScheduleFrequency graphql.Int        `json:"scheduleFrequency"`
	LastDiscoClassScanId     graphql.String     `json:"lastDiscoClassScanId"`
	ScheduleStartAt          graphql.String     `json:"scheduleStartAt"`
	LastDiscoClassScan       struct {
		Type        DiscoClassScanType   `json:"type"`
		Status      DiscoClassScanStatus `json:"status"`
		StartedAt   graphql.String       `json:"startedAt"`
		CompletedAt graphql.String       `json:"completedAt"`
		AvgDuration graphql.Float        `json:"avgDuration"`
	} `json:"lastDiscoClassScan"`
}

type UpdateDiscoClassScanConfigInput struct {
	ID                			graphql.ID         	`json:"id"`
	Enabled           			graphql.Boolean    	`json:"enabled"`
	Type              			*DiscoClassScanType	`json:"type"`
	ScheduleFrequencyMinutes	graphql.Int       	`json:"scheduleFrequency"`
	ScheduleStartAt   			graphql.String    	`json:"scheduleStartAt"`
	// Omitting scanPluginConfigs for now
}

type SombraOutput struct {
	CustomerUrl  graphql.String `graphql:"customerUrl"`
	HostedMethod graphql.String `graphql:"hostedMethod"`
}

func MakeStandaloneUpdatePluginInput(d *schema.ResourceData) UpdatePluginInput {
	return UpdatePluginInput{
		PluginID:                 graphql.String(d.Get("id").(string)),
		DataSiloID:               graphql.String(d.Get("data_silo_id").(string)),
		Enabled:                  graphql.Boolean(d.Get("enabled").(bool)),
		ScheduleFrequencyMinutes: graphql.String(strconv.Itoa(d.Get("schedule_frequency_minutes").(int) * 1000 * 60)),
		ScheduleStartAt:          graphql.String(d.Get("schedule_start_at").(string)),
		ScheduleNow:              graphql.Boolean(false),
	}
}

func MakeUpdatePluginInput(d *schema.ResourceData, configuration map[string]interface{}, pluginId graphql.String) UpdatePluginInput {
	return UpdatePluginInput{
		DataSiloID:               graphql.String(d.Get("id").(string)),
		PluginID:                 pluginId,
		Enabled:                  graphql.Boolean(configuration["enabled"].(bool)),
		ScheduleFrequencyMinutes: graphql.String(strconv.Itoa(configuration["schedule_frequency_minutes"].(int) * 1000 * 60)),
		ScheduleStartAt:          graphql.String(configuration["schedule_start_at"].(string)),
	}
}

func MakeStandaloneUpdateDiscoClassScanConfigInput(d *schema.ResourceData) UpdateDiscoClassScanConfigInput {
	input := UpdateDiscoClassScanConfigInput{
		ID:                graphql.String(d.Get("id").(string)),
		Enabled:           graphql.Boolean(d.Get("enabled").(bool)),
		ScheduleFrequencyMinutes: graphql.Int(d.Get("schedule_frequency_minutes").(int) * 1000 * 60),
	}
	
	// Only set type if it's provided and not empty
	if typeVal, ok := d.Get("type").(string); ok && typeVal != "" {
		discoType := DiscoClassScanType(typeVal)
		input.Type = &discoType
	}
	
	// Only set scheduleStartAt if it's provided and not empty
	if startAtVal, ok := d.Get("schedule_start_at").(string); ok && startAtVal != "" {
		input.ScheduleStartAt = graphql.String(startAtVal)
	}
	
	return input
}

func MakeUpdateDiscoClassScanConfigInput(d *schema.ResourceData, configuration map[string]interface{}, discoClassScanConfigId graphql.String) UpdateDiscoClassScanConfigInput {
	input := UpdateDiscoClassScanConfigInput{
		ID:                discoClassScanConfigId,
		Enabled:           graphql.Boolean(configuration["enabled"].(bool)),
		ScheduleFrequencyMinutes: graphql.Int(configuration["schedule_frequency_minutes"].(int) * 1000 * 60),
	}
	
	// Only set type if it's provided and not empty
	if typeVal, ok := configuration["type"].(string); ok && typeVal != "" {
		discoType := DiscoClassScanType(typeVal)
		input.Type = &discoType
	}
	
	// Only set scheduleStartAt if it's provided and not empty
	if startAtVal, ok := configuration["schedule_start_at"].(string); ok && startAtVal != "" {
		input.ScheduleStartAt = graphql.String(startAtVal)
	}
	
	return input
}

func ReadDiscoClassScanConfigIntoState(d *schema.ResourceData, config DiscoClassScanConfig) {
	// Convert the disco class scan config to the nested block format
	// Even though there's only one config per data silo, Terraform requires it as a list
	
	// Convert frequency from milliseconds back to minutes
	frequencyMinutes := int(config.ScheduleFrequency) / (1000 * 60)
	
	// Handle optional fields safely
	typeStr := ""
	if config.Type != "" {
		typeStr = string(config.Type)
	}
	
	scheduleStartAt := ""
	if config.ScheduleStartAt != "" {
		scheduleStartAt = string(config.ScheduleStartAt)
	}
	
	lastDiscoClassScanId := ""
	if config.LastDiscoClassScanId != "" {
		lastDiscoClassScanId = string(config.LastDiscoClassScanId)
	}
	
	discoClassScanConfig := []interface{}{
		map[string]interface{}{
			"id":                       string(config.ID),
			"enabled":                  bool(config.Enabled),
			"type":                     typeStr,
			"schedule_frequency_minutes": frequencyMinutes,
			"schedule_start_at":        scheduleStartAt,
			"last_disco_class_scan_id": lastDiscoClassScanId,
		},
	}
	d.Set("disco_class_scan_config", discoClassScanConfig)
}

func ReadStandaloneDataSiloPluginIntoState(d *schema.ResourceData, plugin Plugin) {
	frequency, err := strconv.Atoi(string(plugin.ScheduleFrequency))
	if err != nil {
		return
	}

	d.Set("enabled", plugin.Enabled)
	d.Set("id", plugin.ID)
	d.Set("data_silo_id", plugin.DataSilo.ID)
	d.Set("schedule_frequency_minutes", frequency/60/1000)
	d.Set("schedule_start_at", plugin.ScheduleStartAt)
	d.Set("last_enabled_at", plugin.LastEnabledAt)
}

func ReadDataSiloPluginsIntoState(d *schema.ResourceData, plugins []Plugin) {
	for _, plugin := range plugins {
		frequency, err := strconv.Atoi(string(plugin.ScheduleFrequency))
		if err == nil {
			configuration := map[string]interface{}{
				"enabled":                    plugin.Enabled,
				"id":                         plugin.ID,
				"schedule_frequency_minutes": frequency / 60 / 1000,
				"schedule_start_at":          plugin.ScheduleStartAt,
				"last_enabled_at":            plugin.LastRunAt,
			}

			switch plugin.Type {
			case "SCHEMA_DISCOVERY":
				d.Set("schema_discovery_plugin", []interface{}{configuration})
			case "CONTENT_CLASSIFICATION":
				d.Set("content_classification_plugin", []interface{}{configuration})
			case "DATA_SILO_DISCOVERY":
				d.Set("data_silo_discovery_plugin", []interface{}{configuration})
			}
		}
	}
}

func CreateDataSiloUpdatableFields(d *schema.ResourceData) DataSiloUpdatableFields {

	emailsSet := d.Get("owner_emails").(*schema.Set)
	emailsList := emailsSet.List()
	emails := make([]string, len(emailsList))
	for i, v := range emailsList {
		emails[i] = v.(string)
	}
	sort.Strings(emails)
	emailsGraphql := make([]graphql.String, len(emails))
	for i, v := range emails {
		emailsGraphql[i] = graphql.String(v)
	}

	teamsSet := d.Get("owner_teams").(*schema.Set)
	teamsList := teamsSet.List()
	teams := make([]string, len(teamsList))
	for i, v := range teamsList {
		teams[i] = v.(string)
	}
	sort.Strings(teams)
	teamsGraphql := make([]graphql.String, len(teams))
	for i, v := range teams {
		teamsGraphql[i] = graphql.String(v)
	}

	return DataSiloUpdatableFields{
		Title:              graphql.String(d.Get("title").(string)),
		Description:        graphql.String(d.Get("description").(string)),
		URL:                graphql.String(d.Get("url").(string)),
		NotifyEmailAddress: graphql.String(d.Get("notify_email_address").(string)),
		IsLive:             graphql.Boolean(d.Get("is_live").(bool)),
		OwnerEmails:        emailsGraphql,
		OwnerTeams:         teamsGraphql,
		Headers:            ToCustomHeaderInputList((d.Get("headers").([]interface{}))),
		SombraId:           graphql.String(d.Get("sombra_id").(string)),

		// TODO: Add more fields
		// DataSubjectBlockListIds: toStringList(d.Get("data_subject_block_list_ids")),
		// Identifiers:             toStringList(d.Get("identifiers").([]interface{})),
		// "api_key_id":                   graphql.ID(d.Get("api_key_id").(string)),
		// "depended_on_data_silo_titles": toStringList(d.Get("depended_on_data_silo_titles").([]interface{})),
		// "team_names":                   toStringList(d.Get("team_names").([]interface{})),
	}
}

func GetIntegrationName(d *schema.ResourceData) string {
	// Determine the type of the data silo. Most often, this is just the `type` field.
	// But for AVC silos, the `outer_type` actually contains the name to use, as the `type`
	// is always "promptAPerson"
	integrationName := d.Get("outer_type")
	if integrationName == "" {
		integrationName = d.Get("type")
	}

	return integrationName.(string)
}

func CreateDataSiloInput(d *schema.ResourceData) CreateDataSilosInput {
	return CreateDataSilosInput{
		Name:  graphql.String(GetIntegrationName(d)),
		Title: graphql.String(d.Get("title").(string)),
	}
}

type ContextJson struct {
	SecretMap              map[string]string       `json:"secretMap"`
	AllowedHosts           []string                `json:"allowedHosts"`
	AllowedPlaintextPaths  []string                `json:"allowedPlaintextPaths"`
	AllowedIdentifierPaths []AllowedIdentifierPath `json:"allowedIdentifierPaths,omitempty"`
}

// Represents a single allowed identifier path for Sombra integration
type AllowedIdentifierPath struct {
	IdentifierName string `json:"identifierName"`
	Path           string `json:"path"`
}

// Represents the enricherIdentifierMappings structure from the backend
type EnricherIdentifierMapping struct {
	IdentifierName string   `json:"identifierName"`
	Paths          []string `json:"paths"`
}

// Flattens enricherIdentifierMappings into a list of AllowedIdentifierPath
func BuildAllowedIdentifierPaths(mappings []EnricherIdentifierMapping) []AllowedIdentifierPath {
	var allowed []AllowedIdentifierPath
	for _, mapping := range mappings {
		for _, path := range mapping.Paths {
			allowed = append(allowed, AllowedIdentifierPath{
				IdentifierName: mapping.IdentifierName,
				Path:           path,
			})
		}
	}
	return allowed
}

func toStringList(l []graphql.String) []string {
	ret := make([]string, len(l))
	for i, s := range l {
		ret[i] = string(s)
	}
	return ret
}

func ConstructSecretMapString(d *schema.ResourceData, allowedHosts []graphql.String, allowedPlaintextPathObjs []PlaintextInformation, allowedIdentifierPaths []AllowedIdentifierPath) ([]byte, error) {
	// Construct secret map
	contextSet := d.Get("secret_context").(*schema.Set)
	contextMap := map[string]string{}
	for _, rawContext := range contextSet.List() {
		context := rawContext.(map[string]interface{})
		contextMap[context["name"].(string)] = context["value"].(string)
	}

	// Construct plaintext paths
	allowedPlaintextPaths := make([]string, len(allowedPlaintextPathObjs))
	for i, obj := range allowedPlaintextPathObjs {
		allowedPlaintextPaths[i] = string(obj.Path)
	}

	return json.Marshal(ContextJson{
		SecretMap:              contextMap,
		AllowedHosts:           toStringList(allowedHosts),
		AllowedPlaintextPaths:  allowedPlaintextPaths,
		AllowedIdentifierPaths: allowedIdentifierPaths,
	})
}

func ReadDataSiloIntoState(d *schema.ResourceData, silo DataSilo) {
	d.Set("id", silo.ID)
	d.Set("link", silo.Link)
	d.Set("aws_external_id", silo.ExternalId)
	d.Set("has_avc_functionality", silo.Catalog.HasAvcFunctionality)
	d.Set("type", silo.Type)
	d.Set("title", silo.Title)
	if d.Get("description") != nil {
		d.Set("description", silo.Description)
	}
	if d.Get("url") != nil {
		d.Set("url", silo.URL)
	}
	d.Set("outer_type", silo.OuterType)
	if d.Get("notify_email") != nil {
		d.Set("notify_email_address", silo.NotifyEmailAddress)
	}
	if d.Get("is_live") != nil {
		d.Set("is_live", silo.IsLive)
	}
	d.Set("connection_state", silo.ConnectionState)
	d.Set("owner_emails", FlattenOwners(silo))
	d.Set("owner_teams", FlattenOwnerTeams(silo))
	d.Set("headers", FlattenHeaders(&silo.Headers))

	// TODO: Support these fields being read in
	// d.Set("data_subject_block_list", flattenDataSiloBlockList(silo))
	// d.Set("identifiers", silo.Identifiers)
	// d.Set("prompt_email_template_id", silo.PromptEmailTemplate.ID)
	// d.Set("team_names", ...)
	// d.Set("depended_on_data_silo_ids", ...)
	// d.Set("data_subject_block_list_ids", ...)
	// d.Set("api_key_id", ...)
}

func FlattenOwners(dataSilo DataSilo) []interface{} {
	owners := dataSilo.Owners
	ret := make([]string, len(owners))
	for i, owner := range owners {
		ret[i] = string(owner.Email)
	}
	sort.Strings(ret)
	result := make([]interface{}, len(ret))
	for i, v := range ret {
		result[i] = v
	}
	return result
}

func FlattenOwnerTeams(dataSilo DataSilo) []interface{} {
	teams := dataSilo.Teams
	ret := make([]string, len(teams))
	for i, team := range teams {
		ret[i] = string(team.Name)
	}
	sort.Strings(ret)
	result := make([]interface{}, len(ret))
	for i, v := range ret {
		result[i] = v
	}
	return result
}

func FlattenDataSiloBlockList(dataSilo DataSilo) []interface{} {
	owners := dataSilo.SubjectBlocklist
	ret := make([]interface{}, len(owners))
	for i, owner := range owners {
		ret[i] = owner.ID
	}
	return ret
}
