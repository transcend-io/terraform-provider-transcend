terraform {
  required_providers {
    transcend = {
      version = "0.0.2"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
    // this is key for local dev environment
    url = "https://yo.com:4001/"
    key = "2efbb919ae0615431e04e4132976a79aa9528567fedc59b1cf3c908f560348c4"
}

resource "transcend_data_silo" "amazon" {
  type = "amazonS3"
  title = "Amazon"
  url = "https://"
  description = "This is a test"
  headers {
    name = "test"
    value = "what"
    is_secret = true
  }
}

output "amazon_title" {
    value = resource.transcend_data_silo.amazon.title
}
output "amazon_description" {
    value = resource.transcend_data_silo.amazon.description
}
output "amazon_link" {
    value = resource.transcend_data_silo.amazon.link
}
