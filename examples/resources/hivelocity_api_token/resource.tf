resource "hivelocity_api_token" "ci" {
  name        = "ci-pipeline"
  expires_at  = "2027-01-01T00:00:00Z"
  allowed_ips = ["203.0.113.0/24"]
}

output "ci_token" {
  value     = hivelocity_api_token.ci.token
  sensitive = true
}
