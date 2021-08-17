terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1.0"]
      source   = "hivelocity/hivelocity"
    }
  }
}

data "hivelocity_device_initial_creds" "initial_creds" {
  device_id = 13796
}

output "initial_creds" {
  value = data.hivelocity_device_initial_creds.initial_creds
}
