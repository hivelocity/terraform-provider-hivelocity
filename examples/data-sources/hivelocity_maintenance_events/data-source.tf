data "hivelocity_maintenance_events" "upcoming" {}

output "affected_devices" {
  value = distinct(flatten(data.hivelocity_maintenance_events.upcoming.events[*].device_ids))
}
