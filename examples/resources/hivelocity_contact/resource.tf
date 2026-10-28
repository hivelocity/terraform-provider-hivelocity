resource "hivelocity_contact" "oncall" {
  full_name   = "On-call Engineer"
  email       = "oncall@example.com"
  phone       = "+1 555 0100"
  description = "Shared on-call login"
  password    = var.oncall_password
}

variable "oncall_password" {
  type      = string
  sensitive = true
}
