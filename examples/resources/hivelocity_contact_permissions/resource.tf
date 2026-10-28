data "hivelocity_permissions" "all" {}

resource "hivelocity_contact_permissions" "oncall" {
  contact_id  = hivelocity_contact.oncall.contact_id
  permissions = ["device.view", "device.power"]
}
