terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1.0"]
      source   = "hivelocity/hivelocity"
    }
  }
}

data "hivelocity_bare_metal_device" "mine" {}

output "my_bare_metal" {
  value = data.hivelocity_bare_metal_device.mine
}
