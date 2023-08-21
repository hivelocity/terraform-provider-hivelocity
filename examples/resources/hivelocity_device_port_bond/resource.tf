terraform {
    required_providers {
        hivelocity = {
            source  = "hivelocity/hivelocity"
            version = "0.1.0"
        }
    }
}

resource "hivelocity_bare_metal_device" "webserver" {
    product_id    = 580
    os_name       = "CentOS 7.x"
    location_name = "DAL1"
    hostname      = "webserver.terraform.testbond"
    bonded        = true 
}

resource "hivelocity_device_port_bond" "bond" {
    device_id = hivelocity_bare_metal_device.webserver.device_id
}

output "bond_device" {
    value = hivelocity_device_port_bond.bond
}