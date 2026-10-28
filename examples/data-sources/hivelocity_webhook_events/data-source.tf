data "hivelocity_webhook_events" "all" {}

output "event_names" {
  value = data.hivelocity_webhook_events.all.events[*].name
}
