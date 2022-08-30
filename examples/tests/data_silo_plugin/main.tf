terraform {
  required_providers {
    transcend = {
      version = "0.7.0"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.dev.trancsend.com/"
}

variable "type" {
  type = string
  default = "DATA_SILO_DISCOVERY"
}

variable "schedule_frequency_minutes" {
  type = string
  default = "3000"
}

variable "schedule_start_at" {
  type = string
  default = "2030-08-16T07:00:00.000Z"
}

variable "enabled" {
  type = bool
  default = false
}

resource "transcend_data_silo" "gradle" {
  type = "gradle"
}


resource "transcend_data_silo_plugin" "gradle" {
  data_silo_id = resource.transcend_data_silo.gradle.id
  type = var.type
  schedule_frequency_minutes = var.schedule_frequency_minutes
  schedule_start_at = var.schedule_start_at
  enabled = var.enabled
}

output "gradlePluginId" {
  value = resource.transcend_data_silo_plugin.gradle.id
}

output "gradlePluginDataSiloId" {
  value = resource.transcend_data_silo_plugin.gradle.data_silo_id
}

output "gradlePluginType" {
  value = resource.transcend_data_silo_plugin.gradle.type
}

output "gradlePluginEnabled" {
  value = resource.transcend_data_silo_plugin.gradle.enabled
}
