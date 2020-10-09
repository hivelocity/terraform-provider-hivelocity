terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1"]
      source = "hivelocity.net/dev/hivelocity"
    }
  }
}

data "hivelocity_product_options" "all" {}


output "all_product_options" {
  value = data.hivelocity_product_options.all.product_options
}

