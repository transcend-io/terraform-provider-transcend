package types

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
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
	// Description struct {
	// 	DefaultMessage graphql.String
	// }
	// DataCollection struct {
	// 	VisualID graphql.String
	// }
}

type DataPointUpdatableFields struct {
	DataSiloId graphql.String `json:"dataSiloId"`
	Name       graphql.String `json:"name"`
	Title      graphql.String `json:"title"`
	// Description graphql.String `json:"description"`
	// Categories []DataSubCategoryInput    `json:"categories"`
	// Purposes   []PurposeSubCategoryInput `json:"purposes"`
	// Attributes []AttributeInput          `json:"attributes"`
}

type UpdateOrCreateDataPointInput struct {
	ID graphql.String `json:"id,omitempty"`
	DataPointUpdatableFields
}

// type PurposeSubCategoryInput struct {
// 	Name    graphql.String    `json:"name"`
// 	Purpose ProcessingPurpose `json:"purpose"`
// }

// type DataSubCategoryInput struct {
// 	Name     graphql.String   `json:"name"`
// 	Category DataCategoryType `json:"category"`
// }

// type AttributeInput struct {
// 	Key    graphql.String   `json:"key"`
// 	Values []graphql.String `json:"values"`
// }

func MakeUpdateOrCreateDataPointInput(d *schema.ResourceData) UpdateOrCreateDataPointInput {
	return UpdateOrCreateDataPointInput{
		ID: graphql.String(d.Get("id").(string)),
		DataPointUpdatableFields: DataPointUpdatableFields{
			Name:       graphql.String(d.Get("name").(string)),
			DataSiloId: graphql.String(d.Get("data_silo_id").(string)),
			Title:      graphql.String(d.Get("title").(string)),
		},
	}
}

func ReadDataPointIntoState(d *schema.ResourceData, dataPoint DataPoint) {
	d.Set("name", dataPoint.Name)
	d.Set("data_silo_id", dataPoint.DataSilo.ID)
	d.Set("title", dataPoint.Title.DefaultMessage)
}

// func ToDataPointSubDataPointInputList(origs []interface{}) []DataPointSubDataPointInput {
// 	vals := make([]DataPointSubDataPointInput, len(origs))
// 	for i, orig := range origs {
// 		newVal := orig.(map[string]interface{})
// 		vals[i] = DataPointSubDataPointInput{
// 			Name: graphql.String(newVal["name"].(string)),
// 			Description: graphql.String(newVal["description"].(string)),
// 			// ToDataSubCategoryInputList(newVal["categories"].([]interface{})),
// 			// ToPurposeSubCategoryInputList(newVal["purposes"].([]interface{})),
// 			// ToAttributeInputList(newVal["attributes"].([]interface{})),
// 		}
// 	}

// 	return vals
// }

// func ToDataSubCategoryInputList(origs []interface{}) []DataSubCategoryInput {
// 	vals := make([]DataSubCategoryInput, len(origs))
// 	for i, orig := range origs {
// 		newVal := orig.(map[string]interface{})
// 		vals[i] = DataSubCategoryInput{
// 			graphql.String(newVal["name"].(string)),
// 			DataCategoryType(newVal["category"].(string)),
// 		}
// 	}

// 	return vals
// }

// func ToPurposeSubCategoryInputList(origs []interface{}) []PurposeSubCategoryInput {
// 	vals := make([]PurposeSubCategoryInput, len(origs))
// 	for i, orig := range origs {
// 		newVal := orig.(map[string]interface{})
// 		vals[i] = PurposeSubCategoryInput{
// 			graphql.String(newVal["name"].(string)),
// 			ProcessingPurpose(newVal["purpose"].(string)),
// 		}
// 	}

// 	return vals
// }

// func ToAttributeInputList(origs []interface{}) []AttributeInput {
// 	vals := make([]AttributeInput, len(origs))
// 	for i, orig := range origs {
// 		newVal := orig.(map[string]interface{})
// 		vals[i] = AttributeInput{
// 			graphql.String(newVal["key"].(string)),
// 			ToStringList(newVal["values"].([]interface{})),
// 		}
// 	}

// 	return vals
// }
