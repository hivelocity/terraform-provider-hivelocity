resource "hivelocity_vps_snapshot_schedule" "nightly" {
  volume_id     = hivelocity_vps_volume.data.volume_id
  facility_code = hivelocity_vps_volume.data.facility_code
  interval_type = "DAILY"
  hour          = 2
  minute        = 30
  max_snapshots = 7
  timezone      = "America/New_York"
}
