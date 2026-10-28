resource "hivelocity_bare_metal_device" "web" {
  product_id    = 580
  os_name       = "Ubuntu 22.04"
  location_name = "DAL1"
  hostname      = "web.terraform.test"
  bonded        = true
}

resource "hivelocity_device_port_bond" "web" {
  device_id = hivelocity_bare_metal_device.web.device_id
}

# The bond can then be placed in a VLAN.
resource "hivelocity_vlan" "private" {
  facility_code = "DAL1"
  type          = "private"
  port_ids      = [hivelocity_device_port_bond.web.bond_port_id]
}
