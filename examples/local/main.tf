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
    key = "97ced11a53792dd210427191eb12e137d5ade1cd0bb7fc2ba0d5bccf343d2250"
}

data "transcend_data_silo" "data_silos" {
    text = ""
    first = 15
    offset = 0
}

# resource "transcend_data_silo" "amazon" {
#   type = "amazonS3"
#   title = "Amazon"
#   url = "https://"
# }

# output "amazon" {
#   value = resource.transcend_data_silo.amazon
# }

resource "transcend_api_key" "test" {
  title = "test!"
  data_silos = ["09bae972-a340-4cc9-a590-51715ee6d413"]
  scopes = ["fullAdmin"]
}
