terraform {
  required_providers {
    transcend = {
      version = "0.9.1"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.dev.trancsend.com/"
}

variable "title" {}
variable "plugin_config" {
  type = list(object({
    enabled                    = bool
    schedule_frequency_minutes = number
    schedule_start_at          = string
    schedule_now               = bool
  }))
}

resource "transcend_data_silo" "silo" {
  type            = "amazonWebServices"
  title           = var.title
  skip_connecting = true
}

resource "transcend_data_silo_connection" "connection" {
  data_silo_id = transcend_data_silo.silo.id

  plaintext_context {
    name  = "role"
    value = "TranscendAWSIntegrationRole"
  }

  plaintext_context {
    name  = "accountId"
    value = "590309927493"
  }
}

resource "transcend_schema_discovery_plugin" "plugin" {
  for_each = {
    for config in var.plugin_config :
    config.type => config
  }

  data_silo_id = transcend_data_silo.silo.id

  enabled                    = each.value["enabled"]
  schedule_frequency_minutes = each.value["schedule_frequency_minutes"]
  schedule_start_at          = each.value["schedule_start_at"]
  schedule_now               = each.value["schedule_now"]

  depends_on = [transcend_data_silo_connection.connection]
}

output "awsExternalId" {
  value = transcend_data_silo.silo.aws_external_id
}

output "dataSiloId" {
  value = transcend_data_silo.silo.id
}
