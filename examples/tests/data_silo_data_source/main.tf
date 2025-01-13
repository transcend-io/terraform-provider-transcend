terraform {
  required_providers {
    transcend = {
      version = "0.18.15"
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
  owner_emails         = ["david@transcend.io"]
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

output "dataSiloOwners" {
  value = data.transcend_data_silo.silo.owner_emails
}

output "dataSiloDescription" {
  value = data.transcend_data_silo.silo.description
}