resource "hivelocity_order_group" "web_rack" {
  name      = "Web Rack Group"
  same_rack = true

  bare_metal_device {
    product_id    = 580
    os_name       = "Ubuntu 22.04"
    hostname      = "database.terraform.test"
    location_name = "DAL1"
  }

  bare_metal_device {
    product_id    = 580
    os_name       = "Ubuntu 22.04"
    hostname      = "storage.terraform.test"
    location_name = "DAL1"
  }

  bare_metal_device {
    product_id    = 580
    os_name       = "Ubuntu 22.04"
    hostname      = "webserver.terraform.test"
    location_name = "DAL1"
  }
}

output "output_order_group" {
  value = hivelocity_order_group.web_rack
}
