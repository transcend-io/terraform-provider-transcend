terraform {
  required_providers {
    transcend = {
      version = "0.18.15"
      source  = "transcend.com/cli/transcend"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0.0"
    }
  }
}

# Set TRANSCEND_KEY and TRANSCEND_URL locally, or define in this block
provider "transcend" {
  url = "https://api.staging.transcen.dental/"
}

provider "aws" {
  region = "us-east-1"
}

# To use the sombra module, you must declare the AWS and Vault providers explicitly
# Your settings will very likely be different here. 
provider "vault" {
  # You are more than welcome to use real vault credentials here.
  # See https://github.com/hashicorp/terraform-provider-vault/issues/666
  # for an explanation of why a "fake" set of settings is required when using
  # modules that optionally use the vault provider
  address          = "https://vault.dev.trancsend.com"
  token            = "not-a-real-token"
  skip_tls_verify  = true
  skip_child_token = true
}
