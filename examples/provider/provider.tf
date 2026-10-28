terraform {
  required_providers {
    hivelocity = {
      source  = "hivelocity/hivelocity"
      version = "~> 0.8"
    }
  }
}

// The API key can also be supplied with the HIVELOCITY_API_KEY environment variable.
provider "hivelocity" {
  api_key = var.hivelocity_api_key
}

variable "hivelocity_api_key" {
  type      = string
  sensitive = true
}

// Find a plan with 16GB of memory in Tampa.
data "hivelocity_product" "tampa_product" {
  first = true

  filter {
    name   = "product_memory"
    values = ["16GB"]
  }

  filter {
    name   = "data_center"
    values = ["TPA1"]
  }

  filter {
    name   = "stock"
    values = ["limited", "available"]
  }
}

// Create a SSH Key
resource "hivelocity_ssh_key" "my_ssh_key" {
  name       = "This is my Terraform SSH Key"
  public_key = file("~/.ssh/id_ed25519.pub")
}

// Provision your devices with Ubuntu.
resource "hivelocity_bare_metal_device" "webserver" {
  product_id        = data.hivelocity_product.tampa_product.product_id
  os_name           = "Ubuntu 22.04"
  location_name     = data.hivelocity_product.tampa_product.data_center
  hostname          = "webserver.terraform.test"
  tags              = ["hello", "world"]
  script            = file("${path.module}/cloud_init_example.yaml")
  public_ssh_key_id = hivelocity_ssh_key.my_ssh_key.ssh_key_id
}

resource "hivelocity_bare_metal_device" "database" {
  product_id        = data.hivelocity_product.tampa_product.product_id
  os_name           = "Ubuntu 22.04"
  location_name     = data.hivelocity_product.tampa_product.data_center
  hostname          = "database.terraform.test"
  tags              = ["hello", "world"]
  script            = file("${path.module}/cloud_init_example.yaml")
  public_ssh_key_id = hivelocity_ssh_key.my_ssh_key.ssh_key_id
}

// Look up each device's private port (defaults to the bond, or eth1 on unbonded devices).
data "hivelocity_device_port" "webserver" {
  device_id = hivelocity_bare_metal_device.webserver.device_id
}

data "hivelocity_device_port" "database" {
  device_id = hivelocity_bare_metal_device.database.device_id
}

// Create a private VLAN connecting the servers.
resource "hivelocity_vlan" "private_vlan" {
  facility_code = data.hivelocity_product.tampa_product.data_center
  type          = "private"
  port_ids = [
    data.hivelocity_device_port.webserver.port_id,
    data.hivelocity_device_port.database.port_id,
  ]
}
