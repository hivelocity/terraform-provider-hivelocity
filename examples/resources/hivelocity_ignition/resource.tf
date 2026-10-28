resource "hivelocity_ignition" "demo" {
  name     = "This is my Terraform Ignition Config"
  contents = file("${path.module}/ignition.json")
}

output "demo_ignition" {
  value = hivelocity_ignition.demo
}
