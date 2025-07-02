terraform {
  required_providers {
    transcend = {
      version = "0.18.18"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.staging.transcen.dental/"
}

resource "transcend_identifier" "test" {
  name = var.name
}

output "identifierId" {
  value = transcend_identifier.test.id
}

variable "name" {
  type = string
}

