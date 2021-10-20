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

resource "hivelocity_dns_record_a" "demo_record" {
  domain_id = hivelocity_dns_domain.demo_domain.id
  name = "www.domain.com"
  address = "149.255.34.100"
  ttl = 3600
}

resource "hivelocity_dns_record_ptr" "demo_record" {
  address = "149.255.34.100"
  name = "www.domain.com"
  ttl = 3600

  # Ensure that the A record is created before the PTR
  depends_on = [
    hivelocity_dns_record_a.demo_record
  ]
}

output "demo_dns_record" {
  value = hivelocity_dns_record_ptr.demo_record
}
