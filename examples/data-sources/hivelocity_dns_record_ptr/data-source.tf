data "hivelocity_dns_record_ptr" "web" {
  filter {
    name   = "address"
    values = ["66.165.231.122"]
  }
}
