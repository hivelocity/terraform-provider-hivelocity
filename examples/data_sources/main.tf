terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1"]
      source = "hivelocity.net/dev/hivelocity"
    }
  }
}

module "devices" {
  source = "./devices"
}

output "devices" {
  value = module.devices.all_devices
}

module "products" {
  source = "./products"
}

output "products" {
  value = module.products.all_products
}

module "tampa_products" {
  source = "./products"
}

output "tampa_products" {
  value = module.tampa_products.tampa_products
}


module "product_operating_systems" {
  source = "./product_operating_systems"
}

output "product_operating_systems" {
  value = module.product_operating_systems.all_product_operating_systems
}

