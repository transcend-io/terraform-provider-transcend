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

variable "name" {}
variable "title" {}

resource "transcend_data_silo" "silo" {
  type = "server"
  title = var.title
}

resource "transcend_data_point" "point" {
  data_silo_id = transcend_data_silo.silo.id
  name = var.name
  title = var.title
}

output "dataPointId" {
  value = transcend_data_point.point.id
}