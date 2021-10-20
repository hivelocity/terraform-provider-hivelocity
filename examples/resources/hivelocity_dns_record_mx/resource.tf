terraform {
  required_providers {
    hivelocity = {
      source   = "hivelocity/hivelocity"
      version  = "0.1.0"
    }
  }
}

resource "hivelocity_dns_domain" "demo_domain" {
  name = "domain.com"
}

resource "hivelocity_dns_record_mx" "demo_record" {
  domain_id = hivelocity_dns_domain.demo_domain.id
  name = "domain.com"
  exchange = "mail.domain.com"
  preference = 10
  ttl = 3600
}

output "demo_dns_record" {
  value = hivelocity_dns_record_mx.demo_record
}
