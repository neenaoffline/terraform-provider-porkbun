terraform {
  required_providers {
    porkbun = {
      source  = "neena/porkbun"
      version = "0.1.0"
    }
  }
}

# Configure the provider using environment variables:
# PORKBUN_API_KEY and PORKBUN_SECRET_API_KEY
provider "porkbun" {}

variable "domain" {
  description = "The domain to manage DNS records for"
  type        = string
  default     = "example.com"
}

# A record for root domain
resource "porkbun_dns_record" "root_a" {
  domain  = var.domain
  type    = "A"
  content = "192.168.1.1"
  ttl     = "600"
}

# A record for www subdomain
resource "porkbun_dns_record" "www" {
  domain  = var.domain
  name    = "www"
  type    = "A"
  content = "192.168.1.1"
  ttl     = "600"
}

# MX record
resource "porkbun_dns_record" "mx_primary" {
  domain  = var.domain
  type    = "MX"
  content = "mail.${var.domain}"
  prio    = "10"
  ttl     = "600"
}

# TXT record for SPF
resource "porkbun_dns_record" "spf" {
  domain  = var.domain
  type    = "TXT"
  content = "v=spf1 mx -all"
  ttl     = "600"
}

# CNAME record for blog
resource "porkbun_dns_record" "blog" {
  domain  = var.domain
  name    = "blog"
  type    = "CNAME"
  content = "myblog.example.net"
  ttl     = "600"
}

# Output the created record IDs
output "root_record_id" {
  value = porkbun_dns_record.root_a.id
}

output "www_record_id" {
  value = porkbun_dns_record.www.id
}
