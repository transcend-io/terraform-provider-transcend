terraform {
  required_providers {
    transcend = {
      version = "0.13.0"
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
variable "path" {
  type    = list(string)
  default = []
}
variable "properties" {
  type = list(object({
    name        = string
    description = string
    categories = list(object({
      name     = string
      category = string
    }))
    purposes = list(object({
      name    = string
      purpose = string
    }))
    attributes = list(object({
      key    = string
      values = list(string)
    }))
    access_request_visibility_enabled = bool
    erasure_request_redaction_enabled = bool
  }))
  default = [{
    name                              = "test"
    description                       = "test"
    categories                        = []
    purposes                          = []
    attributes                        = []
    access_request_visibility_enabled = false
    erasure_request_redaction_enabled = false
  }]
}

resource "transcend_data_silo" "silo" {
  type            = var.data_silo_type
  title           = var.title
  description     = "Send a webhook to a server and POST back through our API."
  skip_connecting = true
}

resource "transcend_data_point" "point" {
  data_silo_id = transcend_data_silo.silo.id
  name         = var.name
  title        = var.title
  description  = var.description
  path         = var.path

  dynamic "properties" {
    for_each = var.properties
    content {
      name                              = properties.value["name"]
      description                       = properties.value["description"]
      access_request_visibility_enabled = properties.value["access_request_visibility_enabled"]
      erasure_request_redaction_enabled = properties.value["erasure_request_redaction_enabled"]

      dynamic "categories" {
        for_each = properties.value["categories"]
        content {
          name     = categories.value["name"]
          category = categories.value["category"]
        }
      }

      dynamic "purposes" {
        for_each = properties.value["purposes"]
        content {
          name    = purposes.value["name"]
          purpose = purposes.value["purpose"]
        }
      }

      dynamic "attributes" {
        for_each = properties.value["attributes"]
        content {
          key    = attributes.value["key"]
          values = attributes.value["values"]
        }
      }
    }
  }
}

output "properties" {
  value = transcend_data_point.point.properties
}

output "path" {
  value = transcend_data_point.point.path
}

output "dataSiloId" {
  value = transcend_data_silo.silo.id
}

output "dataPointId" {
  value = transcend_data_point.point.id
}
