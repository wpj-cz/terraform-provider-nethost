---
page_title: "nethost Provider"
---

# nethost Provider

The nethost provider manages services exposed by the Nethost API, including `sslguard` and `spongefinger` resources.

## Example Usage

```hcl
provider "nethost" {
  endpoint = "https://klient-api-46qvm.nethost.cz/api/v3"
  api_key  = var.nethost_api_key
}
```

## Argument Reference

* `endpoint` - (Required) The URL endpoint for the Nethost API.
* `api_key` - (Optional) The API key for Nethost authentication. If omitted, the provider also reads `NETHOST_API_KEY` or `API_KEY` from the environment.
