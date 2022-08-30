variable "dd_api_key" { sensitive = true }
variable "dd_app_key" { sensitive = true}

resource "transcend_data_silo" "datadog" {
  type            = "datadog"
  skip_connecting = false

  secret_context {
    name  = "apiKey"
    value = var.dd_api_key
  }
  secret_context {
    name  = "applicationKey"
    value = var.dd_app_key
  }
  secret_context {
    name  = "queryTemplate"
    value = "service:programmatic-remote-seeding AND @email:{{identifier}}"
  }
}
