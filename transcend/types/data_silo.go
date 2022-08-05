package types

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
)

type DataSiloUpdatableFields struct {
	Title                   graphql.String      `json:"title,omitempty"`
	Description             graphql.String      `json:"description,omitempty"`
	URL                     graphql.String      `json:"url,omitempty"`
	NotifyEmailAddress      graphql.String      `json:"notifyEmailAddress,omitempty"`
	IsLive                  graphql.Boolean     `json:"isLive"`
	OwnerEmails             []graphql.String    `json:"ownerEmails"`
	DataSubjectBlockListIds []graphql.String    `json:"dataSubjectBlockListIds"`
	Headers                 []CustomHeaderInput `json:"headers"`

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

type DataSiloInput struct {
	Name graphql.String `json:"name,omitempty"`
	DataSiloUpdatableFields
}

type UpdateDataSiloInput struct {
	Id graphql.ID `json:"id"`
	DataSiloUpdatableFields
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
	// Identifiers        []struct {
	// 	Name graphql.String `json:"name"`
	// } `json:"identifiers"`
	Owners []struct {
		ID    graphql.String `json:"id"`
		Email graphql.String `json:"email"`
	} `json:"owners"`
	SubjectBlocklist []struct {
		ID graphql.String `json:"id"`
	} `json:"subjectBlocklist"`
	Headers   []Header       `json:"headers"`
	OuterType graphql.String `json:"outerType"`

	// TODO: Add support to DataSiloInput first
	// PromptEmailTemplate struct {
	// 	ID graphql.String `json:"id,omitempty"`
	// } `json:"promptEmailTemplate,omitempty"`

	// TODO: Look up the schema here
	// Teams   []struct{} `json:"teams"`
	// ApiKeys []struct{} `json:"apiKeys"`
	// DependentDataSilos []struct{} `json:"dependentDataSilos"`
}

func CreateDataSiloUpdatableFields(d *schema.ResourceData) DataSiloUpdatableFields {
	return DataSiloUpdatableFields{
		Title:              graphql.String(d.Get("title").(string)),
		Description:        graphql.String(d.Get("description").(string)),
		URL:                graphql.String(d.Get("url").(string)),
		NotifyEmailAddress: graphql.String(d.Get("notify_email_address").(string)),
		IsLive:             graphql.Boolean(d.Get("is_live").(bool)),
		OwnerEmails:        ToStringList(d.Get("owner_emails").([]interface{})),
		Headers:            ToCustomHeaderInputList((d.Get("headers").([]interface{}))),

		// TODO: Add more fields
		// DataSubjectBlockListIds: toStringList(d.Get("data_subject_block_list_ids")),
		// Identifiers:             toStringList(d.Get("identifiers").([]interface{})),
		// "api_key_id":                   graphql.ID(d.Get("api_key_id").(string)),
		// "depended_on_data_silo_titles": toStringList(d.Get("depended_on_data_silo_titles").([]interface{})),
		// "team_names":                   toStringList(d.Get("team_names").([]interface{})),
	}
}

func createDataSiloInput(d *schema.ResourceData) DataSiloInput {
	return DataSiloInput{
		Name:                    graphql.String(d.Get("type").(string)),
		DataSiloUpdatableFields: CreateDataSiloUpdatableFields(d),
	}
}

func ReadDataSiloIntoState(d *schema.ResourceData, silo DataSilo) {
	d.Set("id", silo.ID)
	d.Set("link", silo.Link)
	d.Set("aws_external_id", silo.ExternalId)
	d.Set("has_avc_functionality", silo.Catalog.HasAvcFunctionality)
	d.Set("type", silo.Type)
	d.Set("title", silo.Title)
	d.Set("description", silo.Description)
	d.Set("url", silo.URL)
	d.Set("outer_type", silo.OuterType)
	d.Set("notify_email_address", silo.NotifyEmailAddress)
	d.Set("is_live", silo.IsLive)
	d.Set("owner_emails", FlattenOwners(silo))
	d.Set("headers", FlattenHeaders(&silo.Headers))

	// TODO: Support these fields being read in
	// d.Set("data_subject_block_list", flattenDataSiloBlockList(silo))
	// d.Set("identifiers", silo.Identifiers)
	// d.Set("prompt_email_template_id", silo.PromptEmailTemplate.ID)
	// d.Set("team_names", ...)
	// d.Set("depended_on_data_silo_ids", ...)
	// d.Set("data_subject_block_list_ids", ...)
	// d.Set("headers", ...)
	// d.Set("api_key_id", ...)
}

func FlattenOwners(dataSilo DataSilo) []interface{} {
	owners := dataSilo.Owners
	ret := make([]interface{}, len(owners))
	for i, owner := range owners {
		ret[i] = owner.Email
	}
	return ret
}

func FlattenDataSiloBlockList(dataSilo DataSilo) []interface{} {
	owners := dataSilo.SubjectBlocklist
	ret := make([]interface{}, len(owners))
	for i, owner := range owners {
		ret[i] = owner.ID
	}
	return ret
}
