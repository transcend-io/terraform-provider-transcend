terraform {
  required_providers {
    transcend = {
      version = "0.1"
      source  = "transcend.com/cli/transcend-io"
    }
  }
}

provider "transcend" {
    url = "https://yo.com:4001/"
    key = "5ab54b35263d9165fddd2f95a7646eeb63dd37c387d0ca2d4be448750fb43163"
}

data "transcend_data_silo" "all" {
    text = ""
    first = 15
    offset = 0
}

output "all" {
  value = data.transcend_data_silo.all
}
