terraform {
  required_providers {
    transcend = {
      version = "0.5.1"
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
  url = "https://yo.com:4001/"
}

provider "aws" {}
