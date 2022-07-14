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
    key = "cdd83ca4557e935669726a4a518d82fdf173c40f7561a06523191cfa20144dd1"
}

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

}

output "amazon" {
  value = resource.transcend_data_silo.amazon
}

output "customer" {
  value = resource.transcend_data_point.customer.id
}

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
