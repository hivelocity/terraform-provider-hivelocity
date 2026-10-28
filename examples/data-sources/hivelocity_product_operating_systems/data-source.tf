data "hivelocity_product_operating_systems" "plan" {
  product_id = 580
}

output "os_names" {
  value = data.hivelocity_product_operating_systems.plan.operating_systems[*].name
}
