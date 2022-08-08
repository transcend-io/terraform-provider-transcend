package transcend

import (
	"context"

	"github.com/transcend-io/terraform-provider-transcend/transcend/types"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphql "github.com/hasura/go-graphql-client"
)

func resourceDataSiloConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataSiloConnectionsCreate,
		ReadContext:   resourceDataSiloConnectionsRead,
		UpdateContext: resourceDataSiloConnectionsUpdate,
		DeleteContext: resourceDataSiloConnectionsDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_silo_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the data silo to connect",
			},
			"plaintext_context": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "This is where you put non-secretive values that go in the form when connecting a data silo",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the plaintext input",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the plaintext input",
						},
					},
				},
			},
			"connection_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current state of the integration",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDataSiloConnectionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceDataSiloConnectionsUpdate(ctx, d, m)
}

func resourceDataSiloConnectionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var query struct {
		DataSilo types.DataSilo `graphql:"dataSilo(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.String(d.Get("data_silo_id").(string)),
	}

	err := client.graphql.Query(context.Background(), &query, vars, graphql.OperationName("DataSilo"))
	if err != nil {
		return diag.FromErr(err)
	}

	types.ReadDataSiloConnectionIntoState(d, query.DataSilo)
	d.SetId(string(query.DataSilo.ID))

	return nil
}

func resourceDataSiloConnectionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	var connectMutation struct {
		ReconnectDataSilo struct {
			DataSilo types.DataSilo
		} `graphql:"reconnectDataSilo(input: $input, dhEncrypted: $dhEncrypted)"`
	}
	connectVars := map[string]interface{}{
		"input": types.ReconnectDataSiloInput{
			DataSiloId:       graphql.String(d.Get("data_silo_id").(string)),
			PlaintextContext: types.ToPlaintextContextList(d.Get("plaintext_context").([]interface{})),
		},
		"dhEncrypted": graphql.String(""), // This is not needed when no encrypted saas contexts are provided
	}
	err := client.graphql.Mutate(context.Background(), &connectMutation, connectVars, graphql.OperationName("ReconnectDataSilo"))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error connecting data silos",
			Detail:   "Error when connecting data silo: " + err.Error(),
		})
		return diags
	}

	return resourceDataSiloConnectionsRead(ctx, d, m)
}

// Data silos cannot be disconnected, so just no-op
func resourceDataSiloConnectionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
