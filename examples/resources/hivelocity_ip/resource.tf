terraform {
  required_providers {
    hivelocity = {
      source = "hivelocity/hivelocity"
      version = "0.1.0"
    }
  }
}

resource "hivelocity_ip" "demo" {
  purpose = "The intended use of this IP address"
  facility_code = "TPA1"
  prefix_length = 27
}

output "ip_assignment" {
  value = hivelocity_ip.demo
}