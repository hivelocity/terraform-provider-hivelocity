resource "hivelocity_vps_iso" "rescue" {
  name          = "systemrescue"
  description   = "SystemRescue 11"
  url           = "https://example.com/images/systemrescue-11.iso"
  facility_code = "TPA1"
  bootable      = true
}

# Attach it to an instance.
resource "hivelocity_vps" "app" {
  product_id    = 1234
  location_name = "TPA1"
  os_name       = "Ubuntu 22.04"
  hostname      = "app.terraform.test"
  iso_id        = hivelocity_vps_iso.rescue.iso_id
}
