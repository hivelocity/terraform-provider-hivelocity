resource "hivelocity_vps_snapshot" "before_upgrade" {
  volume_id     = hivelocity_vps_volume.data.volume_id
  name          = "before-upgrade"
  facility_code = hivelocity_vps_volume.data.facility_code
}
