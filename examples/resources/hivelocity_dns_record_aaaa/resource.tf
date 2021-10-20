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

resource "hivelocity_dns_record_aaaa" "demo_record" {
  domain_id = hivelocity_dns_domain.demo_domain.id
  name = "www.domain.com"
  address = "fe80::202:b3ff:fe1e:8330"
  ttl = 3600
}

output "demo_dns_record" {
  value = hivelocity_dns_record_aaaa.demo_record
}
