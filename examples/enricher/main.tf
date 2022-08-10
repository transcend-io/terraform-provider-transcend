data "transcend_identifier" "email" {
  text = "email"
}

data "transcend_identifier" "coreIdentifier" {
  text = "coreIdentifier"
}

resource "transcend_enricher" "enricher" {
  title              = "someEnricher"
  description        = "some description"
  actions            = ["ACCESS"]
  input_identifier   = data.transcend_identifier.email.id
  output_identifiers = [data.transcend_identifier.coreIdentifier.id]
  type               = "SERVER"
  url                = "https://some.api.endpoint"
}

output "enricherId" {
  value = transcend_enricher.enricher.id
}
