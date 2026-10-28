resource "hivelocity_bare_metal_device" "demo" {
  product_id    = 525
  os_name       = "Ubuntu 22.04"
  location_name = "DAL1"
  hostname      = "hivelocity.terraform.test"
  tags          = ["hello", "world"]
  script        = file("${path.module}/cloud_init_example.yaml")
}

output "demo_device" {
  value = hivelocity_bare_metal_device.demo
}
