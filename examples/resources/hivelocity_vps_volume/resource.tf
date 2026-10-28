resource "hivelocity_vps_volume" "data" {
  device_id = hivelocity_vps.app.device_id
  size      = 100 # GiB
}
