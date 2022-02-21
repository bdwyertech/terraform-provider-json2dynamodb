package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

func data() *schema.Resource {
	return &schema.Resource{
		Description: "JSON into DynamoDB JSON format",

		ReadContext: dataRead,

		Schema: map[string]*schema.Schema{
			"json": {
				Description: "JSON String",
				Type:        schema.TypeString,
				Required:    true,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newJson, _ := structure.NormalizeJsonString(new)
					oldJson, _ := structure.NormalizeJsonString(old)
					return newJson == oldJson
				},
			},

			"spec": {
				Description: "OpenAPI Schema specification in JSON format to validate the JSON against.",
				Type:        schema.TypeString,
				Optional:    true,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newJson, _ := structure.NormalizeJsonString(new)
					oldJson, _ := structure.NormalizeJsonString(old)
					return newJson == oldJson
				},
			},

			"result": {
				Description: "JSON rendered as DynamoDB JSON",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var jInt interface{}
	if err := json.Unmarshal([]byte(d.Get("json").(string)), &jInt); err != nil {
		return diag.Diagnostics{
			{
				Severity:      diag.Error,
				Summary:       "JSON Handling Failed",
				Detail:        "The data source received an unexpected error while attempting to parse the JSON.",
				AttributePath: cty.GetAttrPath("json"),
			},
		}
	}

	if specJson := d.Get("spec").(string); specJson != "" {
		schema := new(spec.Schema)
		if err := json.Unmarshal([]byte(specJson), schema); err != nil {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "JSON Spec Handling Failed",
					Detail: "The data source received an unexpected error while attempting to build the OpenAPI Specification." +
						fmt.Sprintf("\n\nError: %s", err),
					AttributePath: cty.GetAttrPath("spec"),
				},
			}
		}
		if err := validate.AgainstSchema(schema, jInt, strfmt.Default); err != nil {
			return diag.Diagnostics{
				{
					Severity:      diag.Error,
					Summary:       "JSON Spec Validation Failure",
					Detail:        fmt.Sprint(err),
					AttributePath: cty.GetAttrPath("json"),
				},
			}
		}
	}

	avs, err := attributevalue.MarshalMap(&jInt)
	if err != nil {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "DynamoDB JSON Marshalling Failed",
				Detail: "The data source received an unexpected error while attempting to transform the JSON into DynamoDB Attribute Values." +
					fmt.Sprintf("\n\nError: %s", err),
				AttributePath: cty.GetAttrPath("json"),
			},
		}
	}
	jsonBytes, err := SerializeAttributeMap(avs)
	if err != nil {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "DynamoDB JSON Serialization Failed",
				Detail: "The data source received an unexpected error while attempting to transform the DynamoDB Attribute Values into DynamoDB JSON Format." +
					fmt.Sprintf("\n\nError: %s", err),
				AttributePath: cty.GetAttrPath("json"),
			},
		}
	}

	d.Set("result", string(jsonBytes))

	d.SetId("-")
	return nil
}
