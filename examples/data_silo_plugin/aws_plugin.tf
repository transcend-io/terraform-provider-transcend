resource "transcend_data_silo" "aws" {
  type        = "amazonWebServices"
  description = "Amazon Web Services (AWS) provides information technology infrastructure services to businesses in the form of web services."

  # Normally, Data Silos are connected in this resource. But for AWS, we want to delay connecting until after
  # we create the IAM Role, which must use the `aws_external_id` output from this resource. So instead, we set
  # `skip_connecting` to `true` here and use a `transcend_data_silo_connection` resource below
  skip_connecting = true
  lifecycle { ignore_changes = [plaintext_context, data_silo_discovery_plugin] }
}

data "aws_caller_identity" "current" {}
resource "transcend_data_silo_connection" "connection" {
  data_silo_id = transcend_data_silo.aws.id

  plaintext_context {
    name  = "role"
    value = "TranscendAWSIntegrationRole"
  }

  plaintext_context {
    name  = "accountId"
    value = "590309927493"
  }
}

resource "transcend_data_silo_discovery_plugin" "plugin" {
  data_silo_id = transcend_data_silo.aws.id

  enabled                    = true
  schedule_frequency_minutes = 120
  schedule_start_at          = "2122-09-06T17:51:13.000Z"
  schedule_now               = false

  depends_on = [transcend_data_silo_connection.connection]
}

output "plugin_info" {
  value = transcend_data_silo_discovery_plugin.plugin
}
