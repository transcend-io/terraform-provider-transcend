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
variable "outer_type" { default = null }
variable "type" { default = "amazonWebServices" }
variable "description" { default = "some description" }
variable "owner_emails" {
  type = list(string)
  default = []
}
variable "is_live" {
  type = bool
  default = false
}
variable "url" { default = null }
variable "notify_email_address" { default = null }
# variable "data_silo_identifiers" {
#   type = list(string)
#   default = []
# }
# variable "data_subject_block_list_ids" {
#   type = list(string)
#   default = []
# }
variable "headers" {
  type = list(object({
    name = string
    value = string
    is_secret = bool
  }))
  default = []
}

resource "transcend_data_silo" "silo" {
  type = var.type
  title = var.title
  description = var.description
  owner_emails = var.owner_emails
  is_live = var.is_live

  dynamic "headers" {
    for_each = var.headers
    content {
      name = headers.value["name"]
      value = headers.value["value"]
      is_secret = headers.value["is_secret"]
    }
  }
  url = var.url
  notify_email_address = var.notify_email_address
  outer_type = var.outer_type

  // TODO: Add tests for changing these
  # identifiers = var.data_silo_identifiers
  # data_subject_block_list_ids = var.data_subject_block_list_ids
}

output "awsExternalId" {
  value = transcend_data_silo.silo.aws_external_id
}

output "dataSiloId" {
  value = transcend_data_silo.silo.id
}