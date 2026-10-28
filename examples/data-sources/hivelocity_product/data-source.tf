data "hivelocity_product" "tampa" {
  filter {
    name   = "data_center"
    values = ["TPA1"]
  }
}

output "tampa_product" {
  value = data.hivelocity_product.tampa
}

data "hivelocity_product" "all" {}

output "first_product" {
  value = data.hivelocity_product.all
}
