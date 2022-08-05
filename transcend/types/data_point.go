package types

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

type DataPoint struct {
	ID       graphql.String `json:"id"`
	Name     graphql.String `json:"name"`
	DataSilo struct {
		ID graphql.String `json:"id"`
	} `json:"dataSilo"`
	Title struct {
		DefaultMessage graphql.String `json:"defaultMessage"`
	} `json:"title"`
	Description struct {
		DefaultMessage graphql.String `json:"defaultMessage"`
	} `json:"description"`
	// DataCollection struct {
	// 	DataCetegoryId graphql.String `json:"dataCetegoryId"`
	// } `json:"dataCollection"`
	// DataCollection struct {
	// 	VisualID graphql.String
	// }
}

type AttributeValues struct {
	Name         graphql.String `json:"name"`
	AttributeKey struct {
		Name graphql.String `json:"name"`
	} `json:"attributeKey"`
}

type SubDataPoint struct {
	Name      graphql.String `json:"name"`
	DataPoint struct {
		ID graphql.String `json:"id"`
	} `json:"dataPoint"`
	Description     graphql.String            `json:"description"`
	Categories      []DataSubCategoryInput    `json:"categories"`
	Purposes        []PurposeSubCategoryInput `json:"purposes"`
	AttributeValues []AttributeValues         `json:"attributeValues"`
}

type DataPointSubDataPointInput struct {
	Name        graphql.String            `json:"name"`
	Description graphql.String            `json:"description"`
	Categories  []DataSubCategoryInput    `json:"categories"`
	Purposes    []PurposeSubCategoryInput `json:"purposes"`
	Attributes  []AttributeInput          `json:"attributes"`
}

type DataPointUpdatableFields struct {
	DataSiloId    graphql.String               `json:"dataSiloId"`
	Name          graphql.String               `json:"name"`
	Title         graphql.String               `json:"title"`
	Description   graphql.String               `json:"description"`
	SubDataPoints []DataPointSubDataPointInput `json:"subDataPoints,omitempty"`

	// TODO: Add more fields
	// enabledActions
	// dataCollectionId
	// dataCollectionTag
	// erasureRedactionMethod
	// querySuggestions
}

type UpdateOrCreateDataPointInput struct {
	ID graphql.String `json:"id,omitempty"`
	DataPointUpdatableFields
}

type PurposeSubCategoryInput struct {
	Name    graphql.String    `json:"name"`
	Purpose ProcessingPurpose `json:"purpose"`
}

type DataSubCategoryInput struct {
	Name     graphql.String   `json:"name"`
	Category DataCategoryType `json:"category"`
}

type AttributeInput struct {
	Key    graphql.String   `json:"key"`
	Values []graphql.String `json:"values"`
}

func MakeUpdateOrCreateDataPointInput(d *schema.ResourceData) UpdateOrCreateDataPointInput {
	return UpdateOrCreateDataPointInput{
		ID: graphql.String(d.Get("id").(string)),
		DataPointUpdatableFields: DataPointUpdatableFields{
			Name:          graphql.String(d.Get("name").(string)),
			DataSiloId:    graphql.String(d.Get("data_silo_id").(string)),
			Title:         graphql.String(d.Get("title").(string)),
			Description:   graphql.String(d.Get("description").(string)),
			SubDataPoints: ToDataPointSubDataPointInputList(d.Get("properties").([]interface{})),
		},
	}
}

func ReadDataPointIntoState(d *schema.ResourceData, dataPoint DataPoint, properties []SubDataPoint) {
	d.Set("name", dataPoint.Name)
	d.Set("data_silo_id", dataPoint.DataSilo.ID)
	d.Set("title", dataPoint.Title.DefaultMessage)
	d.Set("description", dataPoint.Description.DefaultMessage)
	d.Set("properties", FromDataPointSubDataPointInputList(properties))
}

func ToDataPointSubDataPointInputList(properties []interface{}) []DataPointSubDataPointInput {
	vals := make([]DataPointSubDataPointInput, len(properties))
	for i, rawProperty := range properties {
		property := rawProperty.(map[string]interface{})
		vals[i] = DataPointSubDataPointInput{
			Name:        graphql.String(property["name"].(string)),
			Description: graphql.String(property["description"].(string)),
			Categories:  ToDataSubCategoryInputList(property["categories"].([]interface{})),
			Purposes:    ToPurposeSubCategoryInputList(property["purposes"].([]interface{})),
			Attributes:  ToAttributeInputList(property["attributes"].([]interface{})),
		}
	}
	return vals
}

func FromDataPointSubDataPointInputList(properties []SubDataPoint) []interface{} {
	vals := make([]interface{}, len(properties))
	for i, property := range properties {
		vals[i] = map[string]interface{}{
			"name":        property.Name,
			"description": property.Description,
			"categories":  FromDataSubCategoryInputList(property.Categories),
			"purposes":    FromPurposeSubCategoryInputList(property.Purposes),
			"attributes":  FromAttributeInputList(property.AttributeValues),
		}
	}
	return vals
}

func ToDataSubCategoryInputList(properties []interface{}) []DataSubCategoryInput {
	vals := make([]DataSubCategoryInput, len(properties))
	for i, rawProperty := range properties {
		property := rawProperty.(map[string]interface{})
		vals[i] = DataSubCategoryInput{
			Name:     graphql.String(property["name"].(string)),
			Category: DataCategoryType(property["category"].(string)),
		}
	}
	return vals
}

func FromDataSubCategoryInputList(categories []DataSubCategoryInput) []map[string]interface{} {
	vals := make([]map[string]interface{}, len(categories))
	for i, category := range categories {
		vals[i] = map[string]interface{}{
			"name":     category.Name,
			"category": category.Category,
		}
	}
	return vals
}

func ToPurposeSubCategoryInputList(categories []interface{}) []PurposeSubCategoryInput {
	vals := make([]PurposeSubCategoryInput, len(categories))
	for i, rawCategory := range categories {
		category := rawCategory.(map[string]interface{})
		vals[i] = PurposeSubCategoryInput{
			Name:    graphql.String(category["name"].(string)),
			Purpose: ProcessingPurpose(category["purpose"].(string)),
		}
	}
	return vals
}

func FromPurposeSubCategoryInputList(categories []PurposeSubCategoryInput) []map[string]interface{} {
	vals := make([]map[string]interface{}, len(categories))
	for i, category := range categories {
		vals[i] = map[string]interface{}{
			"name":    category.Name,
			"purpose": category.Purpose,
		}
	}
	return vals
}

func ToAttributeInputList(attributes []interface{}) []AttributeInput {
	vals := make([]AttributeInput, len(attributes))
	for i, rawAttribute := range attributes {
		attribute := rawAttribute.(map[string]interface{})
		vals[i] = AttributeInput{
			Key:    graphql.String(attribute["key"].(string)),
			Values: ToStringList(attribute["values"].([]interface{})),
		}
	}
	return vals
}

func FromAttributeInputList(attributes []AttributeValues) []map[string]interface{} {
	vals := make([]map[string]interface{}, len(attributes))
	for i, attribute := range attributes {
		vals[i] = map[string]interface{}{
			"key":    attribute.AttributeKey.Name,
			"values": attribute.Name,
		}
	}
	return vals
}
