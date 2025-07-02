terraform {
  required_providers {
    transcend = {
      version = "0.18.17"
      source  = "transcend.com/cli/transcend"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0.0"
    }
  }
}

# Set TRANSCEND_KEY and TRANSCEND_URL locally, or define in this block
provider "transcend" {}

provider "aws" {
  profile = "aws-sso" # A temporary, local profile name from `aws sso login` output
}
