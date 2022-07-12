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
    key = "8a7a93d488eca202cf00d1e71f818df4fc10453f8671ee81381f751f37c86b27"
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
}
