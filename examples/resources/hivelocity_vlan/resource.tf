terraform {
  required_providers {
    hivelocity = {
      source  = "hivelocity/hivelocity"
      version = "0.1.0"
    }
  }
}

resource "hivelocity_bare_metal_device" "webserver" {
  product_id    = 525
  os_name       = "CentOS 7.x"
  location_name = "TPA2"
  hostname      = "webserver.terraform.test"
}

resource "hivelocity_bare_metal_device" "database" {
  product_id    = 525
  os_name       = "CentOS 7.x"
  location_name = "TPA2"
  hostname      = "database.terraform.test"
}

data "hivelocity_device_port" "webserver_private_port" {
  first = true
  device_id = hivelocity_bare_metal_device.webserver.device_id
  filter {
    name  = "name"
    values = ["eth1"]
  }
}

data "hivelocity_device_port" "database_private_port" {
  first = true
  device_id = hivelocity_bare_metal_device.database.device_id
  filter {
    name  = "name"
    values = ["eth1"]
  }
}

resource "hivelocity_vlan" "private_vlan" {
  facility_code = "TPA2"
  type          = "private"
  port_ids      = [
    data.hivelocity_device_port.webserver_private_port.port_id,
    data.hivelocity_device_port.database_private_port.port_id,
  ]
}

output "vlan_info" {
  value = hivelocity_vlan.private_vlan
}
