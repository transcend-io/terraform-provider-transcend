terraform {
  required_providers {
    transcend = {
      version = "0.11.0"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.dev.trancsend.com/"
}

variable "title" {}
variable "scopes" {
  type    = list(string)
  default = []
}
variable "data_silo_type" { default = null }

resource "transcend_data_silo" "silo" {
  count           = var.data_silo_type != null ? 1 : 0
  type            = var.data_silo_type
  title           = var.title
  skip_connecting = true
  lifecycle { ignore_changes = [description] }
}

resource "transcend_api_key" "key" {
  title      = var.title
  scopes     = var.scopes
  data_silos = transcend_data_silo.silo.*.id
}

output "dataSiloId" {
  value = var.data_silo_type != null ? transcend_data_silo.silo[0].id : ""
}

output "apiKeyId" {
  value = transcend_api_key.key.id
}
