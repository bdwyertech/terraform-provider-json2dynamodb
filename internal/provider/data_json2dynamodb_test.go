package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testDataSourceConfig_basic = `
data "json2dynamodb" "test" {
    json = jsonencode(
  	    {
			  name = "briansenvtest"
			  description = "Brian's Test Environment"
			  cookbook_versions = {
				  "wildfly" = "> 0.0.0"
			  }
			  json_class = "Chef::Environment"
			  chef_type = "environment"
			  default_attributes = {
				  wildfly = {
					  config = {
						  abc = 123
					  }
				  }
			  }
        }
    )

	spec = <<EOF
	{
		"$schema": "http://json-schema.org/draft-07/schema",
		"$id": "http://example.com/example.json",
		"type": "object",
		"title": "Chef Environment",
		"description": "Chef Environment",
		"default": {},
		"examples": [
			{
				"name": "ics-ims_helmbuilder_a_int",
				"chef_type": "environment",
				"description": "The Integration environment",
				"default_attributes": {
					"wildfly": {
						"version": "1.2.3"
					}
				},
				"cookbook_versions": {
					"wildfly": "~> 1.0.0"
				},
				"override_attributes": {},
				"json_class": "Chef::Environment"
			}
		],
		"required": [
			"name",
			"json_class",
			"chef_type"
		],
		"properties": {
			"name": {
				"$id": "#/properties/name",
				"type": "string",
				"title": "The name schema",
				"description": "Name of the Chef environment",
				"pattern": "^[A-Za-z0-9_-]+$",
				"minLength": 1,
				"examples": [
					"myapp_dev"
				]
			},
			"chef_type": {
				"$id": "#/properties/chef_type",
				"type": "string",
				"title": "The chef_type schema",
				"description": "An explanation about the purpose of this instance.",
				"pattern": "environment"
			},
			"description": {
				"$id": "#/properties/description",
				"type": "string",
				"title": "The description schema",
				"description": "Description of the Chef Environment",
				"default": "",
				"examples": [
					"The Development environment"
				]
			},
			"default_attributes": {
				"$id": "#/properties/default_attributes",
				"type": "object",
				"title": "The default_attributes schema",
				"description": "Optional. A set of attributes to be applied to all nodes, assuming the node does not already have a value for the attribute. This is useful for setting global defaults that can then be overridden for specific nodes.",
				"default": {},
				"examples": [
					{
						"my": {
							"cool": "value",
							"another": {
								"nested": "value"
							}
						}
					}
				],
				"additionalProperties": true
			},
			"cookbook_versions": {
				"$id": "#/properties/cookbook_versions",
				"type": "object",
				"title": "The cookbook_versions schema",
				"description": "Cookbook versions for the environment",
				"default": {},
				"examples": [
					{
						"wildfly": "~> 0.1.0"
					}
				],
				"additionalProperties": {
					"type": "string"
				}
			},
			"override_attributes": {
				"$id": "#/properties/override_attributes",
				"type": "object",
				"title": "The override_attributes schema",
				"description": "Optional. A set of attributes to be applied to all nodes, even if the node already has a value for an attribute. This is useful for ensuring that certain attributes always have specific values.",
				"default": {},
				"additionalProperties": true
			},
			"json_class": {
				"$id": "#/properties/json_class",
				"type": "string",
				"title": "The json_class schema",
				"description": "An explanation about the purpose of this instance.",
				"pattern": "Chef::Environment"
			}
		},
		"additionalProperties": false
	}
	  EOF
}

output "ddbjson" {
  value = "${data.json2dynamodb.test.result}"
}
`

var basicExpectedOutput = `{"chef_type":{"S":"environment"},"cookbook_versions":{"M":{"wildfly":{"S":"\u003e 0.0.0"}}},"default_attributes":{"M":{"wildfly":{"M":{"config":{"M":{"abc":{"N":"123"}}}}}}},"description":{"S":"Brian's Test Environment"},"json_class":{"S":"Chef::Environment"},"name":{"S":"briansenvtest"}}`

func TestDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConfig_basic,
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.json2dynamodb.test"]
					if !ok {
						return fmt.Errorf("missing data resource")
					}

					outputs := s.RootModule().Outputs

					if o := outputs["ddbjson"].Value.(string); o != basicExpectedOutput {
						return fmt.Errorf("output does not match desired:\n %s", o)
					}

					return nil
				},
			},
		},
	})
}
