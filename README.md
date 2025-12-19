# Terraform Provider for Porkbun DNS

> **Note:** This project was generated entirely by AI using [pi-coder](https://github.com/mariozechner/pi-coder) and Claude Opus 4.5. It has been reviewed by a human.

I created this provider on December 18, 2025 because none of the existing Porkbun Terraform providers were working for me. The Porkbun API has aggressive rate limiting that returns 503 errors, which likely caused issues with other implementations. This provider handles rate limiting with automatic retries and exponential backoff.

[![Tests](https://github.com/neenaoffline/terraform-provider-porkbun/actions/workflows/test.yml/badge.svg)](https://github.com/neenaoffline/terraform-provider-porkbun/actions/workflows/test.yml)
[![Release](https://github.com/neenaoffline/terraform-provider-porkbun/actions/workflows/release.yml/badge.svg)](https://github.com/neenaoffline/terraform-provider-porkbun/actions/workflows/release.yml)

This is a Terraform/OpenTofu provider for managing DNS records on [Porkbun](https://porkbun.com).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0 or [OpenTofu](https://opentofu.org/) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (to build the provider)
- A Porkbun account with API access enabled

## Getting Started with Porkbun API

1. Log in to your Porkbun account
2. Go to **Account** → **API Access**
3. Create a new API key and save both the **API Key** and **Secret Key**
4. Enable API access for each domain you want to manage:
   - Go to **Account** → **Domain Management**
   - Click **Details** on your domain
   - Enable the **API Access** toggle

For more details, see [Porkbun's API documentation](https://porkbun.com/api/json/v3/documentation).

## Installation

### From GitHub Releases (Recommended)

1. Download the latest release for your platform from [GitHub Releases](https://github.com/neenaoffline/terraform-provider-porkbun/releases)

2. Extract and install:
   ```bash
   # Linux (amd64)
   unzip terraform-provider-porkbun_*_linux_amd64.zip
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/neenaoffline/porkbun/0.3.0/linux_amd64/
   mv terraform-provider-porkbun_* ~/.terraform.d/plugins/registry.terraform.io/neenaoffline/porkbun/0.3.0/linux_amd64/terraform-provider-porkbun
   
   # macOS (arm64/Apple Silicon)
   unzip terraform-provider-porkbun_*_darwin_arm64.zip
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/neenaoffline/porkbun/0.3.0/darwin_arm64/
   mv terraform-provider-porkbun_* ~/.terraform.d/plugins/registry.terraform.io/neenaoffline/porkbun/0.3.0/darwin_arm64/terraform-provider-porkbun
   ```

### Building from Source

```bash
git clone https://github.com/neenaoffline/terraform-provider-porkbun.git
cd porkbun-terraform-provider
make install
```

## Usage

### Provider Configuration

```hcl
terraform {
  required_providers {
    porkbun = {
      source  = "neenaoffline/porkbun"
      version = "~> 0.3"
    }
  }
}

provider "porkbun" {
  api_key        = var.porkbun_api_key        # or set PORKBUN_API_KEY env var
  secret_api_key = var.porkbun_secret_api_key # or set PORKBUN_SECRET_API_KEY env var
}
```

### Environment Variables

You can configure the provider using environment variables:

```bash
export PORKBUN_API_KEY="pk1_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export PORKBUN_SECRET_API_KEY="sk1_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

### Creating DNS Records

```hcl
# A record for root domain
resource "porkbun_dns_record" "root" {
  domain  = "example.com"
  type    = "A"
  content = "192.168.1.1"
  ttl     = "600"
}

# A record for subdomain
resource "porkbun_dns_record" "www" {
  domain  = "example.com"
  name    = "www"
  type    = "A"
  content = "192.168.1.1"
  ttl     = "600"
}

# MX record
resource "porkbun_dns_record" "mx" {
  domain  = "example.com"
  type    = "MX"
  content = "mail.example.com"
  prio    = "10"
  ttl     = "600"
}

# TXT record for SPF
resource "porkbun_dns_record" "spf" {
  domain  = "example.com"
  type    = "TXT"
  content = "v=spf1 include:_spf.google.com ~all"
  ttl     = "600"
}

# CNAME record
resource "porkbun_dns_record" "blog" {
  domain  = "example.com"
  name    = "blog"
  type    = "CNAME"
  content = "myblog.netlify.app"
  ttl     = "600"
}

# AAAA record (IPv6)
resource "porkbun_dns_record" "ipv6" {
  domain  = "example.com"
  type    = "AAAA"
  content = "2001:db8::1"
  ttl     = "600"
}
```

### Reading DNS Records (Data Source)

```hcl
data "porkbun_dns_record" "existing" {
  domain = "example.com"
  id     = "123456789"
}

output "record_content" {
  value = data.porkbun_dns_record.existing.content
}
```

### Managing Domain Name Servers

```hcl
# Use custom name servers
resource "porkbun_domain_nameservers" "custom" {
  domain = "example.com"
  nameservers = [
    "ns1.example.net",
    "ns2.example.net",
  ]
}

# Use Porkbun's default name servers
resource "porkbun_domain_nameservers" "porkbun" {
  domain = "example.com"
  nameservers = [
    "curitiba.ns.porkbun.com",
    "fortaleza.ns.porkbun.com",
    "maceio.ns.porkbun.com",
    "salvador.ns.porkbun.com",
  ]
}
```

### Importing Existing Records

You can import existing DNS records using the format `domain/record_id`:

```bash
terraform import porkbun_dns_record.www example.com/123456789
```

You can import existing domain name server configuration:

```bash
terraform import porkbun_domain_nameservers.custom example.com
```

## Resource: porkbun_dns_record

### Argument Reference

| Attribute | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `domain`  | string | Yes      | The domain name (e.g., `example.com`) |
| `name`    | string | No       | The subdomain. Leave empty for root domain. Use `*` for wildcard. |
| `type`    | string | Yes      | Record type: `A`, `MX`, `CNAME`, `ALIAS`, `TXT`, `NS`, `AAAA`, `SRV`, `TLSA`, `CAA`, `HTTPS`, `SVCB` |
| `content` | string | Yes      | The record content/value |
| `ttl`     | string | No       | Time to live in seconds (minimum/default: `600`) |
| `prio`    | string | No       | Priority for MX/SRV records (default: `0`) |
| `notes`   | string | No       | Notes for the record |

### Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | The ID of the DNS record |

## Data Source: porkbun_dns_record

### Argument Reference

| Attribute | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `domain`  | string | Yes      | The domain name |
| `id`      | string | Yes      | The record ID |

### Attribute Reference

All attributes from the resource are available as computed values.

## Resource: porkbun_domain_nameservers

Manages the name servers for a domain. 

**Note:** When this resource is destroyed, the domain's name servers will be reset to Porkbun's default name servers.

### Argument Reference

| Attribute     | Type         | Required | Description |
|---------------|--------------|----------|-------------|
| `domain`      | string       | Yes      | The domain name (e.g., `example.com`) |
| `nameservers` | list(string) | Yes      | List of name server hostnames |

### Attribute Reference

| Attribute | Type   | Description |
|-----------|--------|-------------|
| `id`      | string | The domain name (used as identifier) |

### Porkbun Default Name Servers

If you want to use Porkbun's name servers:
- `curitiba.ns.porkbun.com`
- `fortaleza.ns.porkbun.com`
- `maceio.ns.porkbun.com`
- `salvador.ns.porkbun.com`

## Testing

### Unit Tests

```bash
go test -v ./...
```

### Acceptance Tests

Acceptance tests create real resources on Porkbun. You need valid API credentials and a domain with API access enabled.

```bash
export PORKBUN_API_KEY="pk1_..."
export PORKBUN_SECRET_API_KEY="sk1_..."
export PORKBUN_TEST_DOMAIN="yourdomain.com"
make testacc
```

Or run a specific test:

```bash
make testacc-one TEST=TestAccDNSRecordResource_A
```

## Development

```bash
# Build the provider
make build

# Install locally for testing
make install

# Format code
make fmt

# Run all checks
make check
```

## License

MIT License
