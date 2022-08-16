terraform {
  required_providers {
    transcend = {
      version = "0.5.1"
      # This next line uses the locally built code from `make install`
      source = "transcend.com/cli/transcend"
      # This next line uses the published version from the terraform registry
      # source = "transcend-io/transcend"
    }
  }
}

# Set TRANSCEND_KEY and TRANSCEND_URL locally, or define in this block
provider "transcend" {
  url = "https://api.dev.trancsend.com/"
}

// TODO: Add identifier data source and maybe resource
# resource "transcend_enricher" "test" {
#   title = "Basic Identity Enrichment"
#   description = "Enrich an email address to the userId and phone number"
#   url = "https://example.acme.com/transcend-enrichment-webhook"
#   input_identifier = "a088fc27-ef3d-4396-872d-cda1870ce1bd"
#   output_identifiers = ["8d16d4bf-cbc5-40ad-bdcd-17b1439456c2", "a088fc27-ef3d-4396-872d-cda1870ce1bd"]
#   actions = ["ACCESS"]
#   headers {
#     name = "test"
#     value = "what"
#     is_secret = false
#   }
# }

# resource "transcend_data_silo" "aws" {
#   type = "amazonWebServices"
#   title = "AWS (terraform test)"
#   description = "This is a test"
#   owner_emails = ["david@transcend.io"]
#   is_live = false
# }

# output "awsExternalId" {
#   value = transcend_data_silo.aws.aws_external_id
# }

# # TODO: Support plaintext paths
# # TranscendAWSIntegrationRole
# # 590309927493

resource "transcend_data_silo" "server" {
  type = "server"
  title = "User Data Webhook"
  url = "https://your.company.domain/user/lookup"
  description = "Fetches user data from our internal API"
  owner_emails = ["david@transcend.io"]
  headers {
    name = "someHeaderSentWithWebhook"
    value = "someSecret"
    is_secret = false
  }
}

resource "transcend_data_point" "server" {
  data_silo_id = transcend_data_silo.server.id
  name = "User"
  title = "User Data"

  # properties {
  #   name = "Email"
  #   description = "The email address of a customer"

  #   categories {
  #     name = "Email"
  #     category = "CONTACT"
  #   }
  #   purposes {
  #     name = "Other"
  #     purpose = "ESSENTIAL"
  #   }
  # }

  # properties {
  #   name = "Location"
  #   description = "The user's estimated location"

  #   categories {
  #     name = "Approximate Geolocation"
  #     category = "LOCATION"
  #   }
  #   purposes {
  #     name = "Other"
  #     purpose = "ADDITIONAL_FUNCTIONALITY"
  #   }
  # }
}

# resource "transcend_api_key" "test" {
#   title = "server-key"
#   data_silos = [transcend_data_silo.server.id]
#   scopes = ["makeDataSubjectRequest", "connectDataSilos"]
# }

# resource "transcend_data_silo" "avc" {
#   type                 = "promptAPerson"
#   outer_type           = "coupa"
#   notify_email_address = "dpo@coupa.com"
#   description          = "Coupa is a cloud platform for business spend that offers a fully unified suite of financial applications for business spend management"
#   is_live              = true
# }

# resource "transcend_data_silo" "oauth" {
#   type            = "slack"
#   description     = "Slack is a team communication application providing real-time messaging, archiving, and search for modern teams."
#   skip_connecting = true
# }

# resource "transcend_data_silo" "aws" {
#   type        = "amazonWebServices"
#   description = "Amazon Web Services (AWS) provides information technology infrastructure services to businesses in the form of web services."

#   plaintext_context {
#     name  = "role"
#     value = "TranscendAWSIntegrationRole"
#   }

#   plaintext_context {
#     name  = "accountId"
#     value = "590309927493"
#   }
# }

resource "transcend_data_silo" "gradle" {
  type = "gradle"
}

data "transcend_data_silo_plugin" "gradlePlugin" {
  data_silo_id = resource.transcend_data_silo.gradle.id
  type = "DATA_SILO_DISCOVERY"
}

resource "transcend_data_silo_plugin" "gradle" {
  plugin_id = data.transcend_data_silo_plugin.gradlePlugin.id
  data_silo_id = data.transcend_data_silo_plugin.gradlePlugin.data_silo_id
  type = data.transcend_data_silo_plugin.gradlePlugin.type
  enabled = true
}
