terraform {
  required_providers {
    hivelocity = {
      versions = ["0.0.1"]
      source = "hivelocity.net/prod/hivelocity"
    }
  }
}

resource "hivelocity_bare_metal_device" "demo" {
  product_id = 525
  os_name = "CentOS 7.x"
  location_name = "DAL1"
  hostname = "hivelocity.terraform.test"
  tags = ["hello", "world"]
  script = "#cloud-config\npackage_update: true\npackages:\n - vim"
}

output "demo_device" {
  value = hivelocity_bare_metal_device.demo
}
