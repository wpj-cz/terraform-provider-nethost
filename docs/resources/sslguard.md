---
page_title: "nethost_sslguard Resource"
---

# nethost_sslguard Resource

The `nethost_sslguard` resource manages SSLGuard service entries in Nethost.

## Example Usage

```hcl
resource "nethost_sslguard" "example" {
  domain = "example.cz"
  port   = 443
}
```

## Argument Reference

* `domain` - (Required) The domain to protect. Changing this attribute requires replacement.
* `port` - (Optional) The target port. Defaults to `443`. Changing this attribute requires replacement.

## Attribute Reference

* `id` - The resource ID.
