terraform {
  required_providers {
    hivelocity = {
      versions = ["0.0.1"]
      source = "hivelocity.net/prod/hivelocity"
    }
  }
}

data "hivelocity_product" "tampa" {
  filter {
    name = "data_center"
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

data "hivelocity_bare_metal_device" "mine" {}

output "my_bare_metal" {
  value = data.hivelocity_bare_metal_device.mine
}

data "hivelocity_device" "reg_device" {
  first = true
}

output "reg_device" {
  value = data.hivelocity_device.reg_device
}

data "hivelocity_device" "ip_find" {
  first = true
  filter {
    name = "ip_addresses"
    values = ["66.165.231.122"]
  }
}

output "ip_find" {
  value = data.hivelocity_device.ip_find
}

data "hivelocity_device_initial_creds" "initial_creds" {
  device_id = 13796
}

output "initial_creds" {
  value = data.hivelocity_device_initial_creds.initial_creds
}

data "hivelocity_ssh_key" "ssh_keys_demo" {
  first = true
  filter {
    name = "name"
    values = ["This is my Terraform SSH Key"]
  }
}

output "ssh_keys_demo" {
  value = data.hivelocity_ssh_key.ssh_keys_demo
}

