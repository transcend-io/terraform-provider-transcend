---
page_title: "transcend_api_key Resource - terraform-provider-transcend"
subcategory: ""
description: |-
  
---

# transcend_api_key (Resource)



## Example Usage

```terraform
resource "transcend_data_silo" "silo" {
  type            = "server"
  skip_connecting = true
}

resource "transcend_api_key" "test" {
  title      = "server-key"
  data_silos = [transcend_data_silo.silo.id]
  scopes     = ["makeDataSubjectRequest", "connectDataSilos"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `title` (String) The title used to identify the API key

### Optional

- `data_silos` (List of String) The ids of the data silos to assign to
- `scopes` (List of String) The names of the scopes to add

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import transcend_api_key.key <api_key_id>
```