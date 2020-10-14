Terraform Provider Hivelocity 
=============================

- Website: https://www.terraform.io
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.13.x
-	[Go](https://golang.org/doc/install) 1.15+ (to build the provider plugin)

Building the provider
---------------------

You need to install the provider based on instructions for [installing 3rd-party provider plugins for Terraform 0.13](https://www.hashicorp.com/blog/automatic-installation-of-third-party-providers-with-terraform-0-13).

First you must download the binary for the provider.  You have two options:

1. You can download the latest release from the [release page](https://github.com/hivelocity/terraform-provider-hivelocity/releases).
2. Download the repository directly as a go binary `go get -u github.com/hivelocity/terraform-provider-hivelocity`

Next, create the Hivelocity provider plugin directory. 

```sh
mkdir -p ~/.terraform.d/plugins/hivelocity.net/prod/hivelocity/<VERSION>/<OS>_<ARCH>
```

Note: `<VERSION>` is the most recent version of this provider.

Note: `<OS>` and `<ARCH>` use the Go language's standard OS and architecture names; for example, darwin_amd64.

Finally, symlink your binary to the new plugin directory

```sh
ln -s $GOPATH/bin/terraform-provider-hivelocity ~/.terraform.d/plugins/hivelocity.net/prod/hivelocity/<VERSION>/<OS>_<ARCH>/terraform-provider-hivelocity
```

Set an environment variable containing the Hivelocity API key:
```
export HIVELOCITY_API_KEY=<your-api-key>
```
The API key can also be specified in the provider configuration as shown below.

Example Usage
-------------

You may set the environment variable `HIVELOCITY_API_KEY` instead of setting it on the provider itself.

```tf
terraform {
  required_providers {
    hivelocity = {
      versions = ["<VERSION>"]
      source = "hivelocity.net/prod/hivelocity"
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
    name   = "location"
    values = ["TPA1"]
  }

  filter {
    name   = "stock"
    values = ["limited", "available"]
  }
}

// Provision the product with CentOS 7.
resource "hivelocity_bare_metal_device" "tampa_server" {
    product_id = "${data.hivelocity_product.tampa_product.product_id}"
    os_name = "CentOS 7.x"
    location_name = "${data.hivelocity_product.tampa_product.location}"
    hostname = "hivelocity.terraform.test"
    tags = ["hello", "world"]
}
```

The above provider will deploy and provision the first found in stock (`limited` or `available`) device with `16GB` of memory in the `TPA1`.
It will be provisioned with `CentOS7.x` as the OS and received the tags `hello` and `world`.  The hostname will be set to
`hivelocitiy.terraform.test`.

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
    values = ["66.165.245.226"], "66.165.245.226"]
  }
}
```

Filters also work with fields that are arrays. In the above example the `ip_addresses` field returns an array of ip addresses.
The filter will check if `66.165.245.226` or `66.165.245.225` is in the `ip_addresses` array.

You can pass strings, integers, or booleans to the `values` array on the filter. The type must match the returned type of the
data source schema.

Developing the provider
---------------------------

Please see https://github.com/hivelocity/terraform-provider-hivelocity/blob/main/CONTRIBUTING.md



