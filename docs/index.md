---
page_title: "Provider: Transcend"
description: |-
  Provider for Transcend.io
---

# Transcend Provider

You can create an API Key to use with this provider in [the admin dashboard](https://app.transcend.io/infrastructure/api-keys)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `internal_sombra_key` (String) The API Key to use to talk to a self-hosted sombra. Only used for enterprises with the self-hosted option
- `internal_sombra_url` (String) If set, this URL will be used for sombra operations instead of querying the backend. Useful for reverse proxy instances.
- `key` (String) The API Key to use to talk to Transcend. Ensure it has the scopes to perform whatever actions you need. Can be set using the TRANSCEND_KEY environment variable.
- `url` (String) The custom Transcend backend URL to talk to. Typically can be left to the default production URL.