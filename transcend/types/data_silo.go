package types

import "github.com/shurcooL/graphql"

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
