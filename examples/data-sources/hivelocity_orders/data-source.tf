data "hivelocity_orders" "pending" {
  first = true
  filter {
    name   = "status"
    values = ["verification", "provisioning"]
  }
}
