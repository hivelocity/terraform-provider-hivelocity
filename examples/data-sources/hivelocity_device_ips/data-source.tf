data "hivelocity_device_ips" "web" {
  device_id = 12345
}

output "web_subnets" {
  value = data.hivelocity_device_ips.web.ip_assignments[*].subnet
}
