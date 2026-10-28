# By device ID
data "hivelocity_vps" "app" {
  device_id = 12345
}

# By filtering the instance list
data "hivelocity_vps" "tampa" {
  first = true
  filter {
    name   = "facility_code"
    values = ["TPA1"]
  }
}
