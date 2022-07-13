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
}

output "amazon" {
  value = resource.transcend_data_silo.amazon
}

# resource "transcend_api_key" "test" {
#   title = "testing this "
#   data_silos = []
# }

# output "test" {
#   value = resource.transcend_api_key.test
# }
