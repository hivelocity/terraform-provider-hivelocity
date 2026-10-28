resource "hivelocity_ssh_key" "demo" {
  name       = "This is my Terraform SSH Key"
  public_key = file("~/.ssh/id_ed25519.pub")
}

output "demo_ssh_key_id" {
  value = hivelocity_ssh_key.demo.ssh_key_id
}
