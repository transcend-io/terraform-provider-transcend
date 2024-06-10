terraform {
  required_providers {
    transcend = {
      version = "0.18.1"
      source  = "transcend.com/cli/transcend"
    }
  }
}

provider "transcend" {
  url = "https://api.staging.transcen.dental/"
}

variable "title" {}

data "transcend_identifier" "email" {
  text = "email"
}

data "transcend_identifier" "coreIdentifier" {
  text = "coreIdentifier"
}

resource "transcend_enricher" "enricher" {
  title              = var.title
  description        = "some description"
  actions            = ["ACCESS"]
  input_identifier   = data.transcend_identifier.email.id
  output_identifiers = [data.transcend_identifier.coreIdentifier.id]
  type               = "SERVER"
  url                = "https://api.transcend.io/info" # This is not a real enricher endpoing
}

output "enricherId" {
  value = transcend_enricher.enricher.id
}
