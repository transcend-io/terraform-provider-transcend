resource "transcend_data_silo" "aws" {
  type        = "amazonS3"

  plugin_configuration {
    enabled                    = true
    type                       = "DATA_POINT_DISCOVERY"
    schedule_frequency_minutes = 1440 # 1 day
    schedule_start_at          = "2022-09-06T17:51:13.000Z"
    schedule_now               = false
  }

  # ...other fields...
}

# ...other resources...