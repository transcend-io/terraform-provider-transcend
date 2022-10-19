terraform {
  required_providers {
    transcend = {
      version = "0.8.5"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.dev.trancsend.com/"
}

variable "title" {}
variable "type" { default = "amazonWebServices" }
variable "plugin_config" {
  type = set(object({
    enabled                    = bool
    type                       = string
    schedule_frequency_minutes = number
    schedule_start_at          = string
    schedule_now               = bool
  }))
}

resource "transcend_data_silo" "silo" {
  type                 = var.type
  title                = var.title
}

resource "transcend_data_silo_plugin" "plugin" {
  for_each = var.plugin_config
  enabled                    = each.key["enabled"]
  type                       = each.key["type"]
  schedule_frequency_minutes = each.key["schedule_frequency_minutes"]
  schedule_start_at          = each.key["schedule_start_at"]
  schedule_now               = each.key["schedule_now"]
}

output "awsExternalId" {
  value = transcend_data_silo.silo.aws_external_id
}

output "dataSiloId" {
  value = transcend_data_silo.silo.id
}
