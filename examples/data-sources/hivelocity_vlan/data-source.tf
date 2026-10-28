data "hivelocity_vlan" "private" {
  filter {
    name   = "vlan_tag"
    values = ["1234"]
  }
}
