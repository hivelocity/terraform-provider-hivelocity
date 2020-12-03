terraform {
  required_providers {
    hivelocity = {
      versions = ["0.0.1"]
      source = "hivelocity.net/prod/hivelocity"
    }
  }
}

resource "hivelocity_ssh_key" "demo" {
  name = "My Terraform SSH Key"
  public_key = "My Terraform SSH Public Key"
}

output "demo_device" {
  value = hivelocity_ssh_key.demo
}
