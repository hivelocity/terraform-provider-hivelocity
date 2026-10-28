data "hivelocity_contact" "oncall" {
  filter {
    name   = "email"
    values = ["oncall@example.com"]
  }
}
