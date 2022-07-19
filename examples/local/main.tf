terraform {
  required_providers {
    transcend = {
      version = "0.1"
      source  = "transcend.com/cli/transcend-io"
    }
  }
}

provider "transcend" {
    // this is key for local dev environment
    url = "https://yo.com:4001/"
    key = "d1b5c3d2320df3b3ca9a454cf2f9d6a8d4d7bc876e356352e4bcac48164de4a0"
}

# resource "transcend_enricher" "test" {
#   title = "Basic Identity Enrichment"
#   description = "Enrich an email address to the userId and phone number"
#   url = "https://example.acme.com/transcend-enrichment-webhook"
#   input_identifier = "7e91915f-c7c1-45a8-b67c-40f1e262fa27"
#   output_identifiers = ["7e91915f-c7c1-45a8-b67c-40f1e262fa27"]
#   actions = ["ACCESS"]
#   headers {
#     name = "test"
#     value = "what"
#     is_secret = true
#   }
# }

# data "transcend_data_silo" "data_silos" {
#     text = ""
#     first = 15
#     offset = 0
# }

resource "transcend_data_silo" "amazon" {
  type = "amazonS3"
  title = "Amazon"
  url = "https://"
  headers {
    name = "test"
    value = "what"
    is_secret = true
  }
}

resource "transcend_data_point" "customer" {
  data_silo_id = resource.transcend_data_silo.amazon.id
  name = "customer"
  title = "test"
  query_suggestions {
    suggested_query = "testing"
    request_type = "ACCESS"
  }
  sub_data_points {
    name = "test"
    description = "testing"
    categories {
      name = "financial"
      category = "FINANCIAL"
    }
    purposes {
      name = "essential"
      purpose = "ESSENTIAL"
    }
    attributes {
      key = "something"
      values = ["something"]
    }
  }
}

# }

# output "amazon" {
#   value = resource.transcend_data_silo.amazon
# }

# output "customer" {
#   value = resource.transcend_data_point.customer.id
# }

# resource "transcend_api_key" "test" {
#   title = "testing this "
#   data_silos = []
# }

# resource "transcend_api_key" "test" {
#   title = "test!"
#   data_silos = ["09bae972-a340-4cc9-a590-51715ee6d413"]
#   scopes = ["makeDataSubjectRequest", "connectDataSilos"]
# }
  
# output "test" {
#   value = resource.transcend_api_key.test
# }
