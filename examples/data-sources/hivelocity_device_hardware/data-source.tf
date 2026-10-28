data "hivelocity_device_hardware" "web" {
  device_id = 12345
}

output "web_drives" {
  value = data.hivelocity_device_hardware.web.drives
}
