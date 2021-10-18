# Contributing

This document contains instructions and best practices for developing against `terraform-provider-hivelocity`.

**NOTE: This project is still in Alpha, development best practices are still being refined**

## Assumptions

1. We assume you are developing on OSX.  This document will be updated to support Microsoft and Linux distros in the future.
2. We assume you are working against our local development version of the core api.  If you are not a Hivelocity employee, replace `http://localhost:5065` with `https://core.hivelocity.net`. Be aware you will be charged when using `core.hivelocity.net`. If you'd like to improve the provider without worrying about being charged, please open an issue and leave us an email address to contact you at. 
3. You are familiar with the basics of using Terraform to create custom API calls: https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers
4. You are familiar with the best practices for extending Terraform: https://www.terraform.io/docs/extend/best-practices/index.html

## ENV VARS

You may set `HIVELOCITY_API_URL` to `http://localhost:5065/api/v2` for development.
You may set `HIVELOCITY_API_KEY` if you'd like.

## Generating the go client sources for local development

The provider uses [swagger-codegen](https://github.com/swagger-api/swagger-codegen) to auto-generate sources for the API's client. A copy of Hivelocity's API swagger file is kept in this repository for use during the build process.

To generate the go client sources:

```sh
$ make client
```

## Updating the API swagger file and client's source code

If API endpoints have been added to the core, you will want to update the API swagger file and rebuild the go client. Bear in mind that the updated swagger needs to be checked into the repository.

To fetch an updated version of the swagger file and generate the client:

```sh
$ make swagger
$ make client
```

The swagger file will be fetched from the URL set on HIVELOCITY_API_URL environment variable.

## Generating the documentation for the registry

The documentation for the [Hashicorp registry](https://registry.terraform.io/providers/hivelocity/hivelocity/latest/docs) is auto-generated on build-time and placed into the docs folder. It's based on the examples and descriptions of fields present on the provider's source code.

To generate a local version of the docs:

```sh
$ make docs
```

## Rebuilding the project

Whenever you update the project run the following cmd to add the changes to your Terraform plugins:

```sh
$ make
```

It is recommend you set the environment variable `export TF_LOG=DEBUG` so that you can see and debug API calls while developing new functionality.

## Data Sources

All data sources should be added to the `examples/data-sources/` folder for development.

Once you have updated the example, you can test your new data source from the root of the repo:

```sh
# terraform init && terraform apply
$ EXAMPLE=examples/data-sources/<data-source> make example_apply
```

## Resources

All resources should be added to the `examples/resources/` folder for development.

Once you have updated the example, you can test your new data source from the root of the repo:

To Create/Update:

```sh
# terraform init && terraform apply
$ EXAMPLE=examples/resources/<resource> make example_apply
```

To Delete:

```sh
# terraform init && terraform destroy
$ EXAMPLE=examples/resources/<resource> make example_destroy
```

## Testing provider code

We will mostly acceptance test the provider. So you should not run them all unless you want to create tons of resources.
If you make a change, when you are satisfied with the results, run the one covering whatever you changed.

NOTE: The acceptance test run will create real resources, and potentially charge you if you are developing against the prod API `core.hivelocity.net`.
Please open an issue and share your email address if you would like to run tests without being charged.

```sh
$ make testacc TESTARGS='-run=TestAccHivelocityVlan_basic'
```

To run the full suite of acceptance tests:

```sh
$ make testacc
```

To run the suite of unit tests:

```sh
$ make test
```

## Testing the provider with Terraform

Once you've built the plugin binary (see [Developing the provider](#developing-the-provider) above), it can be incorporated within your Terraform environment using the `-plugin-dir` option. Subsequent runs of Terraform will then use the plugin from your development environment.

```sh
$ terraform init -plugin-dir ~/.terraform.d/plugins/
```

## Cross compiling the provider

If you need to provide a compiled version of the provider under development for usage in another platform, it's going to be necessary to cross-compile the project. To check all available platforms you can use the command `go tool dist list`.

For example, if you're on Linux and you need to build a binary to be used in MacOS:

```sh
$ GOOS=darwin GOARCH=amd64 go build
```

The output file is going to be generated in the in the local directory: `terraform-provider-hivelocity`.

To use this version, place the file in the target machine in the following path:
`~/.terraform.d/plugins/registry.terraform.io/hivelocity/hivelocity/0.1.0/darwin_amd64/terraform-provider-hivelocity`.

Then run `terraform init` command to import the updated version:

```sh
$ terraform init -plugin-dir ~/.terraform.d/plugins/
```
