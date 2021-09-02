Terraform Provider Hivelocity 
=============================

- Website: https://www.terraform.io
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.13.x
-	[Go](https://golang.org/doc/install) 1.15+ (to build the provider plugin)

Example Usage
-------------

You may set the environment variable `HIVELOCITY_API_KEY` instead of setting it on the provider itself.

```tf
terraform {
  required_providers {
    hivelocity = {
      versions = ["<VERSION>"]
      source = "hivelocity/hivelocity"
    }
  }
}

provider "hivelocity" {
  api_key=<YOUR-API-KEY>
}

// Find first in stock product with 16GB of memory in Tampa.
data "hivelocity_product" "tampa_product" {
  first = true

  filter {
    name   = "product_memory"
    values = ["16GB"]
  }

  filter {
    name   = "data_center"
    values = ["TPA1"]
  }

  filter {
    name   = "stock"
    values = ["limited", "available"]
  }
}

data "hivelocity_ssh_key" "ssh_keys" {
  first = true
  
  filter {
    name   = "name"
    values = ["This is my Terraform SSH Key"]
  }
}

// Provision the product with CentOS 7.
resource "hivelocity_bare_metal_device" "tampa_server" {
    product_id = "${data.hivelocity_product.tampa_product.product_id}"
    os_name = "CentOS 7.x"
    location_name = "${data.hivelocity_product.tampa_product.data_center}"
    hostname = "hivelocity.terraform.test"
    tags = ["hello", "world"]
    script = file("${path.module}/cloud_init_example.yaml")
    public_ssh_key_id = "${data.hivelocity_ssh_key.ssh_keys.ssh_key_id}"
}
```

The above provider will deploy and provision the first found in stock (`limited` or `available`) device with `16GB` of memory in the `TPA1`.
It will be provisioned with `CentOS7.x` as the OS and received the tags `hello` and `world`.  The hostname will be set to
`hivelocitiy.terraform.test`. Will use a cloud-init user-data script to be executed on the device first boot.

Cloud-init script start with #cloud-init, must be a valid YAML script.
Post-install script start with #!/bin/bash.
Cloud-init or Post-install script just can be used with OS Ubuntu and Centos.

Note: all these values can be changed from the portal, which could make terraform states get out of sync.  We are working on 
a solution to help users avoid this issue.

### Filtering

All data sources are filterable.  You may create as many filters as you like as demonstrated above. The data source will return the value that matches all filters.

```tf
data "hivelocity_product" "tampa_product" {
  first = true

  filter {
    name   = "location"
    values = ["TPA1", "TPA2"]
  }
}
```

The above config has a single filter. The initial thing you will notice is the `first = true` variable.  This means
that if our filters match multiple objects return the first object that matched.  If you want Terraform to throw an error
 when multiple objects are matched set `first = false`.   Next the provider will look all `hivelocity_product`s that are available and return
the first product who's location is either `TPA1` or `TPA2`

```tf
data "hivelocity_device" "ip_filter" {
  first = false

  filter {
    name   = "ip_addresses"
    values = ["66.165.245.226", "66.165.245.226"]
  }
}
```

Filters also work with fields that are arrays. In the above example the `ip_addresses` field returns an array of ip addresses.
The filter will check if `66.165.245.226` or `66.165.245.225` is in the `ip_addresses` array.

You can pass strings, integers, or booleans to the `values` array on the filter. The type must match the returned type of the
data source schema.

How to test Cloud-Init script locally
---------------------------

Please see https://api-docs.hivelocity.net/cloud-init-test

Developing the provider
---------------------------

Please see https://github.com/hivelocity/terraform-provider-hivelocity/blob/main/CONTRIBUTING.md
