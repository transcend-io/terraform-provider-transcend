terraform {
  required_providers {
    transcend = {
      version = "0.4.1"
      source = "transcend.com/cli/transcend"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0.0"
    }
  }
}

# Set TRANSCEND_KEY and TRANSCEND_URL locally, or define in this block
provider "transcend" {
  url = "https://api.dev.trancsend.com/"
}

provider "aws" {}
