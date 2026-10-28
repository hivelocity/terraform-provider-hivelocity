resource "hivelocity_vps_network" "backend" {
  facility_code = "TPA1"
  name          = "backend"
}

# Bridge a VPS network to a VLAN shared with bare metal devices.
resource "hivelocity_vlan" "shared" {
  facility_code = "TPA1"
  type          = "private"
}

resource "hivelocity_vps_network" "bridged" {
  facility_code = "TPA1"
  name          = "bridged"
  vlan_id       = hivelocity_vlan.shared.id
}
