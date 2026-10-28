resource "hivelocity_dns_domain" "example" {
  name = "example.com"
}

resource "hivelocity_dns_record_aaaa" "www" {
  domain_id = hivelocity_dns_domain.example.domain_id
  name      = "www"
  address   = "2001:db8::1"
  ttl       = 3600
}
