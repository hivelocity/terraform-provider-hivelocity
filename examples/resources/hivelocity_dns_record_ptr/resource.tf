# PTR records already exist for your IPs; look the record up by address
# and manage its name.
data "hivelocity_dns_record_ptr" "web" {
  filter {
    name   = "address"
    values = ["66.165.231.122"]
  }
}

resource "hivelocity_dns_record_ptr" "web" {
  record_id = data.hivelocity_dns_record_ptr.web.record_id
  name      = "web.example.com"
  ttl       = 3600
}
