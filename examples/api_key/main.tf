resource "transcend_data_silo" "silo" {
  type            = "server"
  skip_connecting = true
}

resource "transcend_api_key" "test" {
  title = "server-key"
  data_silos = [transcend_data_silo.silo.id]
  scopes = ["makeDataSubjectRequest", "connectDataSilos"]
}