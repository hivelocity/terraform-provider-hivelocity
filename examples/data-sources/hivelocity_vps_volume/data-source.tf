data "hivelocity_vps_volume" "root" {
  filter {
    name   = "device_id"
    values = ["12345"]
  }
  filter {
    name   = "type"
    values = ["ROOT"]
  }
}
