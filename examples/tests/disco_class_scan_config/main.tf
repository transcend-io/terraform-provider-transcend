terraform {
  required_providers {
    transcend = {
      version = "0.20.0"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.staging.transcen.dental/"
}

variable "title" {}
variable "disco_class_scan_config_vars" {
  type = object({
    enabled                    = bool
    type                       = string
    schedule_frequency_minutes = number
    schedule_start_at          = string
  })
}

resource "transcend_data_silo" "silo" {
  type            = "amazonDynamodb"
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

  // Enable item-level access
  plaintext_context {
    name  = "database"
    value = "true"
  }
}

resource "transcend_disco_class_scan_config" "disco_class_scan_config" {
  data_silo_id = transcend_data_silo.silo.id

  enabled                    = var.disco_class_scan_config_vars["enabled"]
  type                       = var.disco_class_scan_config_vars["type"]
  schedule_frequency_minutes = var.disco_class_scan_config_vars["schedule_frequency_minutes"]
  schedule_start_at          = var.disco_class_scan_config_vars["schedule_start_at"]

  depends_on = [transcend_data_silo_connection.connection]
}

output "awsExternalId" {
  value = transcend_data_silo.silo.aws_external_id
}

output "dataSiloId" {
  value = transcend_data_silo.silo.id
}
