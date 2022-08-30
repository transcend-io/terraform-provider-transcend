resource "transcend_data_silo" "silo" {
  type            = "server"
  skip_connecting = true
}

resource "transcend_data_point" "customer" {
  data_silo_id = transcend_data_silo.silo.id
  name = "customer"
  title = "whatever"

  properties {
    name = "test"
    description = "testing"

    categories {
      name = "Other"
      category = "FINANCIAL"
    }
    categories {
      name = "Biometrics"
      category = "HEALTH"
    }

    purposes {
      name = "essential"
      purpose = "ESSENTIAL"
    }

    attributes {
      key = "something"
      values = ["something"]
    }
  }
}