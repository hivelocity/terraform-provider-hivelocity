# Without a filter the private port (the bond, or eth1 on unbonded devices) is returned.
data "hivelocity_device_port" "private" {
  device_id = 12345
}

data "hivelocity_device_port" "public" {
  device_id = 12345
  filter {
    name   = "name"
    values = ["eth0"]
  }
}
