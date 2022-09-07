resource "transcend_data_silo" "aws" {
  type        = "amazonWebServices"
  description = "Amazon Web Services (AWS) provides information technology infrastructure services to businesses in the form of web services."

  plugin_configuration {
    enabled                    = true
    type                       = "DATA_SILO_DISCOVERY"
    schedule_frequency_minutes = 1440 # 1 day
    schedule_start_at          = "2022-09-06T17:51:13.000Z"
    schedule_now               = false
  }

  # ...other fields...
}

# ...other resources...