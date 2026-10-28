resource "hivelocity_vps_network_attachment" "app_backend" {
  device_id  = hivelocity_vps.app.device_id
  network_id = hivelocity_vps_network.backend.network_id
}
