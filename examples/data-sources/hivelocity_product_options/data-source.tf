data "hivelocity_product_options" "plan" {
  product_id = 580
}

output "ram_upgrades" {
  value = [for o in data.hivelocity_product_options.plan.options : o if o.category == "RAM"]
}
