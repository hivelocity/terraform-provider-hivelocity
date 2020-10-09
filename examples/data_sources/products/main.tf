terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1"]
      source = "hivelocity.net/dev/hivelocity"
    }
  }
}

data "hivelocity_products" "tampa" {
  filter {
    name   = "location"
    values = ["TPA1"]
  }
}


output "tampa_products" {
  value = data.hivelocity_products.tampa
}

data "hivelocity_products" "all" {}

output "all_products" {
  value = data.hivelocity_products.all
}
