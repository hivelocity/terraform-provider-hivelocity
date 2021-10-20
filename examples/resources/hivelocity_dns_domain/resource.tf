terraform {
  required_providers {
    hivelocity = {
      source   = "hivelocity/hivelocity"
      version  = "0.1.0"
    }
  }
}

resource "hivelocity_dns_domain" "demo" {
  name = "domain.com"
}

output "demo_dns_domain" {
  value = hivelocity_dns_domain.demo
}
