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
  contents = jsonencode({"ignition": {"version": "2.4.0"},"systemd": {"units": [{"name": "example.service","enabled": true,"contents": "[Service]\nType=oneshot\nExecStart=/usr/bin/echo Hello World\n\n[Install]\nWantedBy=multi-user.target"}]}})
}

output "demo_ignition" {
  value = hivelocity_ignition.demo
}
