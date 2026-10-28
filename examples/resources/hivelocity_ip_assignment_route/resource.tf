data "hivelocity_ip_assignment" "subnet" {
  filter {
    name   = "subnet"
    values = ["66.165.231.128/29"]
  }
}

# Statically route the subnet to an address configured on one of your devices.
resource "hivelocity_ip_assignment_route" "subnet" {
  assignment_id = data.hivelocity_ip_assignment.subnet.assignment_id
  next_hop_ip   = "66.165.231.10"
}
