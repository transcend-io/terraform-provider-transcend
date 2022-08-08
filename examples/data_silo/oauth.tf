resource "transcend_data_silo" "oauth" {
  type            = "slack"
  description     = "Slack is a team communication application providing real-time messaging, archiving, and search for modern teams."
  skip_connecting = true
}