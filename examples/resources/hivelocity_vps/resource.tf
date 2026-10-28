# Pick a VPS plan that is in stock in Tampa.
data "hivelocity_product" "vps" {
  first = true

  filter {
    name   = "is_vps"
    values = ["true"]
  }

  filter {
    name   = "data_center"
    values = ["TPA1"]
  }

  filter {
    name   = "stock"
    values = ["available", "limited"]
  }
}

resource "hivelocity_ssh_key" "admin" {
  name       = "admin"
  public_key = file("~/.ssh/id_ed25519.pub")
}

resource "hivelocity_vps" "app" {
  product_id        = data.hivelocity_product.vps.product_id
  location_name     = data.hivelocity_product.vps.data_center
  os_name           = "Ubuntu 22.04"
  hostname          = "app.terraform.test"
  public_ssh_key_id = hivelocity_ssh_key.admin.ssh_key_id
  tags              = ["app"]
  script            = file("${path.module}/cloud_init.yaml")
}

output "app_ip" {
  value = hivelocity_vps.app.primary_ip
}
