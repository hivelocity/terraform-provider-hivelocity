terraform {
  required_providers {
    hivelocity = {
      source = "hivelocity/hivelocity"
      version = "0.1.0"
    }
  }
}

resource "hivelocity_ip_assignment" "demo" {
  purpose = "IP Address Asked to Test Terraform Provider SP-292"
  facility_code = "DAL1"
  prefix_length = "27"
}

output "ip_assignment" {
  value = hivelocity_ip_assignment.demo
}