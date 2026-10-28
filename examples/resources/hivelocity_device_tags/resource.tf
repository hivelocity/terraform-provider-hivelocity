# Tag a device that was not provisioned by Terraform.
resource "hivelocity_device_tags" "legacy" {
  device_id = 12345
  tags      = ["legacy", "do-not-reload"]
}
