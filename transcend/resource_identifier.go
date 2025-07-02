package transcend

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/transcend-io/terraform-provider-transcend/transcend/types"
)

func resourceIdentifier() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentifierCreate,
		ReadContext:   resourceIdentifierRead,
		UpdateContext: resourceIdentifierUpdate,
		DeleteContext: resourceIdentifierDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the identifier.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the identifier (e.g. EMAIL, PHONE, etc). Must match backend enum.",
			},
		},
	}
}

func resourceIdentifierCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	var mutation struct {
		CreateIdentifier struct {
			Identifier types.Identifier
		} `graphql:"createIdentifier(input: $input)"`
	}

	vars := map[string]interface{}{
		"input": types.MakeIdentifierInput(d),
	}

	err := client.graphql.Mutate(ctx, &mutation, vars, graphql.OperationName("CreateIdentifier"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error creating identifier",
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId(string(mutation.CreateIdentifier.Identifier.ID))
	types.ReadIdentifierIntoState(d, mutation.CreateIdentifier.Identifier)
	return nil
}

func resourceIdentifierRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	var query struct {
		Identifier types.Identifier `graphql:"identifier(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.ID(d.Get("id").(string)),
	}

	err := client.graphql.Query(ctx, &query, vars, graphql.OperationName("Identifier"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading identifier",
			Detail:   err.Error(),
		})
		return diags
	}
	types.ReadIdentifierIntoState(d, query.Identifier)
	return nil
}

func resourceIdentifierUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	var mutation struct {
		UpdateIdentifier struct {
			Identifier types.Identifier
		} `graphql:"updateIdentifier(input: $input)"`
	}

	vars := map[string]interface{}{
		"input": types.MakeUpdateIdentifierInput(d),
	}

	err := client.graphql.Mutate(ctx, &mutation, vars, graphql.OperationName("UpdateIdentifier"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error updating identifier",
			Detail:   err.Error(),
		})
		return diags
	}
	types.ReadIdentifierIntoState(d, mutation.UpdateIdentifier.Identifier)
	return nil
}

func resourceIdentifierDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	var mutation struct {
		DeleteIdentifiers struct {
			Success graphql.Boolean
		} `graphql:"deleteIdentifiers(ids: $ids)"`
	}

	vars := map[string]interface{}{
		"ids": []graphql.ID{graphql.ID(d.Get("id").(string))},
	}

	err := client.graphql.Mutate(ctx, &mutation, vars, graphql.OperationName("DeleteIdentifiers"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error deleting identifier",
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId("")
	return nil
}
