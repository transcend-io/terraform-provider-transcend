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
variable "description" { default = null }
variable "data_silo_type" { default = "server" }
variable "properties" {
  type = list(object({
    name = string
    description = string
    categories = list(object({
      name = string
      category = string
    }))
  }))
  default = []
}

resource "transcend_data_silo" "silo" {
  type = var.data_silo_type
  title = var.title
}

resource "transcend_data_point" "point" {
  data_silo_id = transcend_data_silo.silo.id
  name = var.name
  title = var.title
  description = var.description

  dynamic "properties" {
    for_each = var.properties
    content {
      name = properties.value["name"]
      description = properties.value["description"]

      dynamic "categories" {
        for_each = properties.value["categories"]
        content {
          name = categories.value["name"]
          category = categories.value["category"]
        }
      }
    }
  }
}

output "properties" {
  value = transcend_data_point.point.properties
}

output "dataSiloId" {
  value = transcend_data_silo.silo.id
}

output "dataPointId" {
  value = transcend_data_point.point.id
}