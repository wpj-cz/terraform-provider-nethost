---
page_title: "nethost_spongefinger Resource"
---

# nethost_spongefinger Resource

The `nethost_spongefinger` resource manages Spongefinger services in Nethost.

## Example Usage

```hcl
resource "nethost_spongefinger" "example" {
  name       = "My WAF"
  variant_id = 1
  domain     = "example.cz"
  subdomains = [
    "www.example.cz",
    "shop.example.cz",
  ]
  backends = [
    "origin.example.cz:80",
  ]
}
```

## Argument Reference

* `name` - (Required) Service name.
* `variant_id` - (Required) Variant or plan ID. Changing this attribute requires replacement.
* `domain` - (Required) Protected root domain.
* `subdomains` - (Optional) Additional subdomains. Do not include the root domain.
* `backends` - (Required) Backend targets in `host:port` format.

## Attribute Reference

* `id` - The resource ID.
* `name` - The service name.
* `variant_id` - The selected variant ID.
* `domain` - The protected root domain.
* `subdomains` - The configured subdomains.
* `backends` - The configured backend targets.
