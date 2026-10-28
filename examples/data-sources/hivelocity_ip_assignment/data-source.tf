# By ID
data "hivelocity_ip_assignment" "by_id" {
  assignment_id = 4321
}

# By an address inside the assignment
data "hivelocity_ip_assignment" "by_address" {
  ip_address = "66.165.231.122"
}

# By filtering the list
data "hivelocity_ip_assignment" "tampa_v6" {
  first = true
  filter {
    name   = "facility_code"
    values = ["TPA1"]
  }
  filter {
    name   = "version"
    values = ["6"]
  }
}
