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

resource "hivelocity_ignition" "demo" {
  name     = "This is my Terraform Ignition Config"
  contents = file("${path.module}/ignition.json")
}

output "demo_ignition" {
  value = hivelocity_ignition.demo
}
