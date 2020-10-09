terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1"]
      source = "hivelocity.net/dev/hivelocity"
    }
  }
}

data "hivelocity_devices" "all" {}


output "all_devices" {
  value = data.hivelocity_devices.all.devices
}

