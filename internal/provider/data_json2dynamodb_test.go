package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
			  default_attributes = {
				  wildfly = {
					  config = {
						  abc = 123
					  }
				  }
			  }
        }
    )

	spec = jsonencode(
		{
			properties = {
				name = {
					type = "string"
					pattern = "^[A-Za-z]+$"
					minLength = 1
				}
			}
			
			patternProperties = {
				"address-[0-9]+" = {
					type = "string"
					pattern = "^[\\s|a-z]+$"
				}
			}
			required = [
				"name"
			]
			// additionalProperties = false
		}
	)
}

output "ddbjson" {
  value = "${data.json2dynamodb.test.result}"
}
`

var basicExpectedOutput = `{"chef_type":{"S":"environment"},"cookbook_versions":{"M":{"wildfly":{"S":"\u003e 0.0.0"}}},"default_attributes":{"M":{"wildfly":{"M":{"config":{"M":{"abc":{"N":"123"}}}}}}},"description":{"S":"Brian's Test Environment"},"json_class":{"S":"Chef::Environment"},"name":{"S":"briansenvtest"}}`

func TestDataSource_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
