package transcend

import "github.com/shurcooL/graphql"

// Enums
type DbIntegrationQuerySuggestionInput string
type RequestActionObjectResolver string
type DataCategoryType string
type ProcessingPurpose string
type RequestAction string
type DataSiloConnectionState string

type DataPoint struct {
	Name  graphql.String
	Title struct {
		DefaultMessage graphql.String
	}
	Description struct {
		DefaultMessage graphql.String
	}
	DataCollection struct {
		VisualID graphql.String
	}
}

type DataPointSubDataPointInput struct {
	Name       graphql.String            `json:"name"`
	Desciption graphql.String            `json:"description"`
	Categories []DataSubCategoryInput    `json:"categories"`
	Purposes   []PurposeSubCategoryInput `json:"purposes"`
	Attributes []AttributeInput          `json:"attributes"`
}

type DataSubCategoryInput struct {
	Name     graphql.String   `json:"name"`
	Category DataCategoryType `json:"category"`
}

type PurposeSubCategoryInput struct {
	Name    graphql.String    `json:"name"`
	Purpose ProcessingPurpose `json:"purpose"`
}

type AttributeInput struct {
	Key    graphql.String   `json:"key"`
	Values []graphql.String `json:"values"`
}

type Enricher struct {
	ID              graphql.String
	Title           graphql.String
	Description     graphql.String
	Url             graphql.String
	InputIdentifier struct {
		ID graphql.String
	}
	Identifiers []struct {
		ID graphql.String
	}
	Actions []RequestAction
	Headers []Header
}

type Header struct {
	Name     graphql.String  `json:"name"`
	Value    graphql.String  `json:"value"`
	IsSecret graphql.Boolean `json:"isSecret"`
}

type CustomHeaderInput Header

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

// TODO: Add plaintextContext

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

type APIKey struct {
	ID     graphql.String `json:"id"`
	Title  graphql.String `json:"title"`
	Scopes []struct {
		Type graphql.String `json:"type"`
	} `json:"scopes"`
}

type ScopeName string
