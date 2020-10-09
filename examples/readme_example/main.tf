terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1"]
      source = "hivelocity.net/dev/hivelocity"
    }
  }
}

// Find a plan with 16GB of memory in Tampa.
data "hivelocity_products" "tampa_product" {
  first = true

  filter {
    name   = "product_memory"
    values = ["16GB"]
  }

  filter {
    name   = "location"
    values = ["TPA1"]
  }

  filter {
    name   = "stock"
    values = ["limited", "available"]
  }
}

// Provision your device with CentOS 7.
resource "hivelocity_bare_metal_devices" "tampa_server" {
    product_id = "${data.hivelocity_products.tampa_product.product_id}"
    os_name = "CentOS 7.x"
    location_name = "${data.hivelocity_products.tampa_product.location}"
    hostname = "hivelocity.terraform.test"
    tags = ["hello", "world"]
}
