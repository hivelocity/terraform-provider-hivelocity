data "hivelocity_ignition" "ignition_config_demo" {
  first = true

  filter {
    name   = "name"
    values = ["my-ignition-config"]
  }
}

output "ignition_config_demo" {
  value = data.hivelocity_ignition.ignition_config_demo
}
