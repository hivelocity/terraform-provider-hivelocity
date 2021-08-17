terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1.0"]
      source   = "hivelocity/hivelocity"
    }
  }
}

data "hivelocity_device" "reg_device" {
  first = true
}

output "reg_device" {
  value = data.hivelocity_device.reg_device
}

data "hivelocity_device" "ip_find" {
  first = true
  filter {
    name   = "ip_addresses"
    values = ["66.165.231.122"]
  }
}

output "ip_find" {
  value = data.hivelocity_device.ip_find
}
