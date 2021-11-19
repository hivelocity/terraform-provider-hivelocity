terraform {
  required_providers {
    hivelocity = {
      source   = "hivelocity/hivelocity"
      version  = "0.1.0"
    }
  }
}

resource "hivelocity_order_group" "web_rack" {
  name = "Web Rack Group"
  same_rack = true

  bare_metal_device {
    product_id    = 580
    os_name       = "CentOS 7.x"
    hostname      = "database.terraform.test"
    location_name = "DAL1"
  }

  bare_metal_device {
    product_id    = 580
    os_name       = "CentOS 7.x"
    hostname      = "storage.terraform.test"
    location_name = "DAL1"
  }

  bare_metal_device {
    product_id    = 580
    os_name       = "CentOS 7.x"
    hostname      = "webserver.terraform.test"
    location_name = "DAL1"
  }
}

output "output_order_group" {
  value = hivelocity_order_group.web_rack
}
