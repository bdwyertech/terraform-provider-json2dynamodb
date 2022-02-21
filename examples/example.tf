terraform {
  required_providers {
    json2dynamodb = {
      source = "bdwyertech/json2dynamodb"
    }
  }
}

data "json2dynamodb" "test" {
  json = jsonencode(
    {
      name        = "briansenvtest"
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
          type      = "string"
          pattern   = "^[A-Za-z]+$"
          minLength = 1
        }
      }

      patternProperties = {
        "address-[0-9]+" = {
          type    = "string"
          pattern = "^[\\s|a-z]+$"
        }
      }
      required = [
        "name"
      ]
      additionalProperties = false
    }
  )
}

output "ddbjson" {
  value = data.json2dynamodb.test.result
}
