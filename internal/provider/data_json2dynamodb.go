package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &JSON2DynamoDBDataSource{}

func NewJSON2DynamoDBDataSource() datasource.DataSource {
	return &JSON2DynamoDBDataSource{}
}

// JSON2DynamoDBDataSource defines the data source implementation.
type JSON2DynamoDBDataSource struct {
	// client *http.Client
}

// JSON2DynamoDBDataSourceModel describes the data source data model.
type JSON2DynamoDBDataSourceModel struct {
	JSON   jsontypes.Normalized `tfsdk:"json"`
	Spec   jsontypes.Normalized `tfsdk:"spec"`
	Result jsontypes.Normalized `tfsdk:"result"`
	Id     types.String         `tfsdk:"id"`
}

func (d *JSON2DynamoDBDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName // + "_data"
}

func (d *JSON2DynamoDBDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "JSON into DynamoDB JSON format",

		Attributes: map[string]schema.Attribute{
			"json": schema.StringAttribute{
				MarkdownDescription: "JSON String",
				Required:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"spec": schema.StringAttribute{
				MarkdownDescription: "OpenAPI Schema specification in JSON format to validate the JSON against.",
				Optional:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"result": schema.StringAttribute{
				MarkdownDescription: "JSON rendered as DynamoDB JSON",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this data source",
				Computed:            true,
			},
		},
	}
}

func (d *JSON2DynamoDBDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	// if req.ProviderData == nil {
	// 	return
	// }
	// client, ok := req.ProviderData.(*http.Client)
	//if !ok {
	//	resp.Diagnostics.AddError(
	//		"Unexpected Data Source Configure Type",
	//		fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	//	)
	//
	//	return
	//}

	// d.client = client
}

func (d *JSON2DynamoDBDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data JSON2DynamoDBDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var jInt interface{}
	if err := json.Unmarshal([]byte(data.JSON.ValueString()), &jInt); err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("json"),
			"JSON Handling Failed",
			"The data source received an unexpected error while attempting to parse the JSON.",
		)
		return
	}

	if data.Spec.ValueString() != "" {
		schema := new(spec.Schema)
		if err := schema.UnmarshalJSON([]byte(data.Spec.ValueString())); err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("spec"),
				"JSON Spec Handling Failed",
				fmt.Sprintf("The data source received an unexpected error while attempting to build the OpenAPI Specification.\n\nError: %s", err),
			)
			return
		}

		if err := validate.AgainstSchema(schema, jInt, strfmt.Default); err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("json"),
				"JSON Spec Validation Failure",
				fmt.Sprint(err),
			)
			return
		}
	}
	avs, err := attributevalue.MarshalMap(&jInt)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("json"),
			"DynamoDB JSON Marshalling Failed",
			fmt.Sprintf("The data source received an unexpected error while attempting to transform the JSON into DynamoDB Attribute Values.\n\nError: %s", err),
		)
		return
	}
	jsonBytes, err := SerializeAttributeMap(avs)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("json"),
			"DynamoDB JSON Serialization Failed",
			fmt.Sprintf("The data source received an unexpected error while attempting to transform the DynamoDB Attribute Values into DynamoDB JSON Format.\n\nError: %s", err),
		)
		return
	}
	data.Result = jsontypes.NewNormalizedValue(string(jsonBytes))
	data.Id = types.StringValue("-")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
