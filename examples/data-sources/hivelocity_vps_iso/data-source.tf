data "hivelocity_vps_iso" "rescue" {
  filter {
    name   = "name"
    values = ["systemrescue"]
  }
}
