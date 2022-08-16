terraform {
  required_providers {
    transcend = {
      version = "0.5.1"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.dev.trancsend.com/"
}

variable "schedule_frequency" {
  type = string
  default = "3000"
}

variable "schedule_start_at" {
  type = string
  default = "2022-08-16T07:00:00.000Z"
}

variable "enabled" {
  type = bool
}

resource "transcend_data_silo" "gradle" {
  type = "gradle"
}

data "transcend_data_silo_plugin" "gradlePlugin" {
  data_silo_id = resource.transcend_data_silo.gradle.id
  type = "DATA_SILO_DISCOVERY"
}

resource "transcend_data_silo_plugin" "gradle" {
  data_silo_id = data.transcend_data_silo_plugin.gradlePlugin.data_silo_id
  type = data.transcend_data_silo_plugin.gradlePlugin.type
  schedule_frequency = var.schedule_frequency
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
