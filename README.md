⚠️⚠️⚠️ WARNING! This provider is currently in active development. ⚠️⚠️⚠️

⚠️⚠️⚠️ Documentation may lag functionality, and vice-versa. ⚠️⚠️⚠️

⚠️⚠️⚠️ Use at your own risk in a production environment. ⚠️⚠️⚠️

Terraform Provider Hivelocity 
=============================

- Website: https://www.terraform.io
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.13+ (to build the provider plugin)

Building the provider
---------------------

Clone repository to: `$GOPATH/src/github.com/hivelocity/terraform-provider-hivelocity`

```sh
$ mkdir -p $GOPATH/src/github.com/hivelocity; cd $GOPATH/src/github.com/hivelocity
$ git clone git@github.com:hivelocity/terraform-provider-hivelocity
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/hivelocity/terraform-provider-hivelocity
$ make build
```

Example Usage
-------------

You may set the environment variable `HIVELOCITY_API_KEY` instead of setting it on the provider itself.

```tf
terraform {
  required_providers {
    hivelocity = {
      versions = ["0.1"]
      source = "hivelocity.net/dev/hivelocity"
    }
  }
}

// Find first in stock product with 16GB of memory in Tampa.
data "hivelocity_products" "tampa_product" {
  first = true

  filter {
    name   = "product_memory"
    values = ["16GB"]
  }

  filter {
    name   = "location"
    values = ["TPA1"]
  }

  filter {
    name   = "stock"
    values = ["limited", "available"]
  }
}

// Provision the product with CentOS 7.
resource "hivelocity_bare_metal_devices" "tampa_server" {
    product_id = "${data.hivelocity_products.tampa_product.product_id}"
    os_name = "CentOS 7.x"
    location_name = "${data.hivelocity_products.tampa_product.location}"
    hostname = "hivelocity.terraform.test"
    tags = ["hello", "world"]
}
```

Developing the provider
---------------------------

Please see https://github.com/hivelocity/terraform-provider-hivelocity/blob/main/CONTRIBUTING.md



