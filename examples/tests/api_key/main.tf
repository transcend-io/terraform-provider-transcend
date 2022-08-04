terraform {
  required_providers {
    transcend = {
      version = "0.3.0"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.dev.trancsend.com/"
}

variable "title" {}
variable "scopes" {
  type = list(string)
  default = []
}

resource "transcend_api_key" "key" {
  title = var.title
  scopes = var.scopes
}

output "apiKeyId" {
  value = transcend_api_key.key.id
}