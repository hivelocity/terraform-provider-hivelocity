data "hivelocity_webhook_events" "all" {}

resource "hivelocity_webhook" "device_power" {
  name  = "Notify ops on power changes"
  event = "device.power"
  url   = "https://ops.example.com/hooks/hivelocity"
  headers = {
    Authorization = "Bearer ${var.hook_token}"
  }
}

variable "hook_token" {
  type      = string
  sensitive = true
}
