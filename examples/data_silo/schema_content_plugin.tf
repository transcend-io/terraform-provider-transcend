resource "transcend_data_silo" "aws" {
  type        = "amazonWebServices"
  description = "Amazon Web Services (AWS) provides information technology infrastructure services to businesses in the form of web services."

  schema_discovery_plugin {
    enabled                    = true
    schedule_frequency_minutes = 1440 # 1 day
    schedule_start_at          = "2022-09-06T17:51:13.000Z"
  }

  content_classification_plugin {
    enabled                    = true
    schedule_frequency_minutes = 1440 # 1 day
    schedule_start_at          = "2022-09-06T17:51:13.000Z"
  }

  # ...other fields...
}

# ...other resources...
