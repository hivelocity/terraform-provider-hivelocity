terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1"]
      source = "hivelocity.net/dev/hivelocity"
    }
  }
}

data "hivelocity_product_operating_systems" "all" {
  product_id = 589
}


output "all_product_operating_systems" {
  value = data.hivelocity_product_operating_systems.all.product_operating_systems
}
