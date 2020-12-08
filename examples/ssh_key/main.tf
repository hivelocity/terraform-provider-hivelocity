terraform {
  required_providers {
    hivelocity = {
      versions = ["0.0.1"]
      source = "hivelocity.net/prod/hivelocity"
    }
  }
}

resource "hivelocity_ssh_key" "demo" {
  name = "This is my Terraform SSH Key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDYYGqXPIJ1yQ7pMTSmdB+LrjvnPh2jSJewvinxYwuk/CBvbsUgFxOTBtjzfU0kHIlo5QyIhJ7tjh8PhT8LdJPGB86ItwTg3Lmt1q5UFxbHmZ0kPmmzaDI/9aakOal3P93D14HDCzBnkTHfC8/JZ5JpDxp86XM+gWQ9sFMkLx83ZOwONNu3E8PowTEpsp0jMx2B2aFeSM+T4bLkJJQtA5Cp6lgRAc5AklXaqmpdAil/fIL/+gvRf8kJIccAe/2oVj4flaMK7mgZ39qhzYYUjTEokEYvJf17QbdtFTxxtIQ+hTxzKzwT6p8cMu6DNQLfq6oxzbuBVGTvKD79MR5vjx+RNpaIru8wzIspHTez9eGDzdR0316GWDcxmMwVIDMM+3pjopDJV6DILfs6hVlAuH11yCX8YwwGHYpsdzLLd00FEEaGLGVRDr/hvduZ1caQIvdNln6Gr7k6W51U1VTC3NRq49yoxYSsXxn30PfTe2IxFaZyhQXHunCLaMCF+TrAOc0= someone@somewhere.com"
}

output "demo_device" {
  value = hivelocity_ssh_key.demo
}
