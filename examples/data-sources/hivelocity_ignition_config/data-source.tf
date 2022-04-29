terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1.0"]
      source   = "hivelocity/hivelocity"
    }
  }
}

provider "hivelocity" {
  api_key="<YOUR-API-KEY>"
}

data "hivelocity_ignition" "ignition_config_demo" {
  first = true
}

output "ignition_config_demo" {
  value = data.hivelocity_ignition.ignition_config_demo
}
