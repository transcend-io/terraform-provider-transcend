terraform {
  required_providers {
    transcend = {
      version = "0.18.4"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.staging.transcen.dental/"
}

variable "title" {}
variable "plugin_config" {
  type = object({
    enabled                    = bool
    schedule_frequency_minutes = number
    schedule_start_at          = string
  })
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

  plaintext_context {
    name  = "region"
    value = "eu-west-1"
  }
}

resource "transcend_data_silo_discovery_plugin" "plugin" {
  data_silo_id = transcend_data_silo.silo.id

  enabled                    = var.plugin_config["enabled"]
  schedule_frequency_minutes = var.plugin_config["schedule_frequency_minutes"]
  schedule_start_at          = var.plugin_config["schedule_start_at"]

  depends_on = [transcend_data_silo_connection.connection]
}

output "awsExternalId" {
  value = transcend_data_silo.silo.aws_external_id
}

output "dataSiloId" {
  value = transcend_data_silo.silo.id
}
