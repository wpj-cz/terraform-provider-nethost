# terraform-provider-nethost

Terraform provider for the Nethost API.

## Supported Resources

- `nethost_sslguard`
- `nethost_spongefinger`

## Provider Configuration

```hcl
provider "nethost" {
  endpoint = "https://klient-api-46qvm.nethost.cz/api/v3"
  api_key  = var.nethost_api_key
}
```

The `api_key` value can also be supplied through:

- `NETHOST_API_KEY`
- `API_KEY`

## Documentation

Resource documentation is in `docs/`:

- `docs/index.md`
- `docs/resources/sslguard.md`
- `docs/resources/spongefinger.md`

## Development

Generate docs:

```bash
go generate ./...
```

Build the provider:

```bash
go build ./...
```

Release artifacts are produced with GoReleaser:

```bash
goreleaser release --snapshot --clean
```
