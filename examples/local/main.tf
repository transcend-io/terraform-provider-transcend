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
    key = "2efbb919ae0615431e04e4132976a79aa9528567fedc59b1cf3c908f560348c4"
}

resource "transcend_enricher" "test" {
  title = "Basic Identity Enrichment"
  description = "Enrich an email address to the userId and phone number"
  url = "https://example.acme.com/transcend-enrichment-webhook"
  input_identifier = "a088fc27-ef3d-4396-872d-cda1870ce1bd"
  output_identifiers = ["8d16d4bf-cbc5-40ad-bdcd-17b1439456c2", "a088fc27-ef3d-4396-872d-cda1870ce1bd"]
  actions = ["ACCESS"]
  headers {
    name = "test"
    value = "what"
    is_secret = true
  }
}

output "test" {
  value = resource.transcend_enricher.test
}

# data "transcend_data_silo" "data_silos" {
#     text = ""
#     first = 15
#     offset = 0
# }

# resource "transcend_data_silo" "bigquery" {
#   type = "googleBigQuery"
#   title = "Google BigQuery"
#   url = "https://"
#   description = "This is a test"
#   headers {
#     name = "test"
#     value = "what"
#     is_secret = true
#   }
# }

# resource "transcend_data_point" "customer" {
#   data_silo_id = resource.transcend_data_silo.bigquery.id
#   name = "customer"
#   title = "whatever"
#   data_collection_tag = "test"
#   query_suggestions {
#     suggested_query = "testing"
#     request_type = "ACCESS"
#   }
#   sub_data_points {
#     name = "test"
#     description = "testing"
#     categories {
#       name = "Other"
#       category = "FINANCIAL"
#     }
#     categories {
#       name = "Biometrics"
#       category = "HEALTH"
#     }
#     purposes {
#       name = "essential"
#       purpose = "ESSENTIAL"
#     }
#     attributes {
#       key = "something"
#       values = ["something"]
#     }
#   }
# }

# output "amazon" {
#   value = resource.transcend_data_silo.bigquery
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
