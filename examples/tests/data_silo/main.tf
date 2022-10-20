terraform {
  required_providers {
    transcend = {
      version = "0.9.0"
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
  type    = list(string)
  default = []
}
variable "is_live" {
  type    = bool
  default = false
}
variable "skip_connecting" {
  type    = bool
  default = true
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
    name      = string
    value     = string
    is_secret = bool
  }))
  default = []
}
variable "secret_context" {
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}
variable "plugin_config" {
  type = list(object({
    enabled                    = bool
    type                       = string
    schedule_frequency_minutes = number
    schedule_start_at          = string
    schedule_now               = bool
  }))
  default = []
}

resource "transcend_data_silo" "silo" {
  type                 = var.type
  title                = var.title
  description          = var.description
  owner_emails         = var.owner_emails
  is_live              = var.is_live
  url                  = var.url
  notify_email_address = var.notify_email_address
  outer_type           = var.outer_type
  skip_connecting      = var.skip_connecting

  dynamic "plugin_configuration" {
    for_each = var.plugin_config
    content {
      enabled                    = plugin_configuration.value["enabled"]
      type                       = plugin_configuration.value["type"]
      schedule_frequency_minutes = plugin_configuration.value["schedule_frequency_minutes"]
      schedule_start_at          = plugin_configuration.value["schedule_start_at"]
      schedule_now               = plugin_configuration.value["schedule_now"]
    }
  }

  dynamic "plaintext_context" {
    for_each = (var.type == "amazonWebServices" || var.type == "amazonS3") && !var.skip_connecting ? [
      { name = "role", value = "TranscendAWSIntegrationRole" },
      { name = "accountId", value = "590309927493" },
    ] : []
    content {
      name  = plaintext_context.value["name"]
      value = plaintext_context.value["value"]
    }
  }

  dynamic "secret_context" {
    for_each = var.secret_context
    content {
      name  = secret_context.value["name"]
      value = secret_context.value["value"]
    }
  }

  dynamic "headers" {
    for_each = var.headers
    content {
      name      = headers.value["name"]
      value     = headers.value["value"]
      is_secret = headers.value["is_secret"]
    }
  }

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
