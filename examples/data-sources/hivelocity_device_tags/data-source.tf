# All tags in use on the account
data "hivelocity_device_tags" "all" {}

# Tags of one device
data "hivelocity_device_tags" "web" {
  device_id = 12345
}
