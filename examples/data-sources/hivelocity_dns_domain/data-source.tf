data "hivelocity_dns_domain" "example" {
  filter {
    name   = "name"
    values = ["example.com"]
  }
}
