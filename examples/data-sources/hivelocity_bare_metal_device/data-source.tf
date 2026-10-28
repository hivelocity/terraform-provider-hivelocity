data "hivelocity_bare_metal_device" "mine" {
  first = true

  filter {
    name   = "hostname"
    values = ["webserver.terraform.test"]
  }
}

output "my_bare_metal" {
  value = data.hivelocity_bare_metal_device.mine
}
