terraform {
  required_providers {
    hivelocity = {
      source   = "hivelocity/hivelocity"
      version  = "0.1.0"
    }
  }
}

resource "hivelocity_bare_metal_device" "webserver" {
  product_id    = 525
  os_name       = "CentOS 7.x"
  location_name = "DAL1"
  hostname      = "webserver.terraform.test"
}

resource "hivelocity_bare_metal_device" "database" {
  product_id    = 525
  os_name       = "CentOS 7.x"
  location_name = "DAL1"
  hostname      = "database.terraform.test"
}

resource "hivelocity_vlan" "demo" {
}

output "demo_vlan" {
  value = hivelocity_vlan.demo
}
