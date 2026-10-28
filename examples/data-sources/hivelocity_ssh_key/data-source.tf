data "hivelocity_ssh_key" "ssh_keys_demo" {
  first = true
  filter {
    name   = "name"
    values = ["This is my Terraform SSH Key"]
  }
}

output "ssh_keys_demo" {
  value = data.hivelocity_ssh_key.ssh_keys_demo
}
