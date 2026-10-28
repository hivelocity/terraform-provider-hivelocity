data "hivelocity_vps_snapshot" "before_upgrade" {
  filter {
    name   = "name"
    values = ["before-upgrade"]
  }
}
