data "hivelocity_device_port" "public" {
  device_id = 12345
  filter {
    name   = "name"
    values = ["eth0"]
  }
}

data "hivelocity_ip_assignment" "extra" {
  ip_address = "66.165.231.130"
}

# Route an additional IP assignment to the device's public port.
resource "hivelocity_device_port" "public" {
  device_id         = 12345
  port_id           = data.hivelocity_device_port.public.port_id
  ip_assignment_ids = [data.hivelocity_ip_assignment.extra.assignment_id]
}
