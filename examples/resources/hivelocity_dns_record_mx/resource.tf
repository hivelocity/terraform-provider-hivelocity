resource "hivelocity_dns_domain" "example" {
  name = "example.com"
}

resource "hivelocity_dns_record_mx" "mail" {
  domain_id  = hivelocity_dns_domain.example.domain_id
  name       = "@"
  exchange   = "mail.example.com"
  preference = 10
  ttl        = 3600
}
