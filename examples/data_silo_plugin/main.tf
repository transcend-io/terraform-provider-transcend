data "transcend_data_silo_plugin" "bigquery" {
  data_silo_id = "73b59e00-0c39-4f69-ba57-d1a1cd34fb6e"
  type         = "DATA_POINT_DISCOVERY"
}

resource "transcend_data_silo_plugin" "bigquery" {
  data_silo_id               = data.transcend_data_silo_plugin.bigquery.data_silo_id
  type                       = "DATA_POINT_DISCOVERY"
  schedule_frequency_minutes = "3000"
  schedule_start_at          = "2022-08-16T07:00:00.000Z"
  enabled                    = true
}