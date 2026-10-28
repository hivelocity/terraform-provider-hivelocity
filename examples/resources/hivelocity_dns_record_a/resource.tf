resource "hivelocity_dns_domain" "example" {
  name = "example.com"
}

resource "hivelocity_dns_record_a" "www" {
  domain    = hivelocity_dns_domain.example.name
  name      = "www"
  addresses = ["66.165.231.122", "66.165.231.123"]
  ttl       = 3600
}
