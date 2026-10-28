data "hivelocity_permissions" "all" {}

output "permission_names" {
  value = data.hivelocity_permissions.all.permissions
}
