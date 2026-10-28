data "hivelocity_device_ipmi" "web" {
  device_id = 12345
}

output "web_temperatures" {
  value = [for s in data.hivelocity_device_ipmi.web.sensors : s if s.group == "Temperature"]
}
