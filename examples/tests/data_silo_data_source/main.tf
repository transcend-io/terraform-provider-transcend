terraform {
  required_providers {
    transcend = {
      version = "0.18.12"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.staging.transcen.dental/"
}

variable "title" {}

resource "transcend_data_silo" "silo" {
  title                = var.title

  type                 = "server"
  description          = "Send a webhook to a server and POST back through our API."
  skip_connecting      = true
}

data "transcend_data_silo" "silo" {
  id = transcend_data_silo.silo.id
}

output "dataSiloId" {
  value = data.transcend_data_silo.silo.id
}

output "dataSiloTitle" {
  value = data.transcend_data_silo.silo.title
}

output "dataSiloLink" {
  value = data.transcend_data_silo.silo.link
}