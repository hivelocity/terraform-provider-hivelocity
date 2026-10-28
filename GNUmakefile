ifndef HIVELOCITY_API_URL
HIVELOCITY_API_URL:=https://core.hivelocity.net/api/v2
endif

ifndef GOPATH
GOPATH:=$(shell go env GOPATH)
endif

ifndef GOOS
GOOS:=$(shell go env GOHOSTOS)
endif

ifndef GOARCH
GOARCH:=$(shell go env GOHOSTARCH)
endif

OS_ARCH:=$(GOOS)_$(GOARCH)
BUILDPATH:=$(HOME)/.terraform.d/plugins/registry.terraform.io/hivelocity/hivelocity/0.1.0/$(OS_ARCH)
# Beware that changing the version of swagger-codegen may produce an incompatible version of the client.
# Requires a Java 8+ runtime (e.g. `brew install openjdk@17`).
SWAGGER_CODEGEN_CLI:=https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/2.4.15/swagger-codegen-cli-2.4.15.jar

default: build

install:
	go install

build: client
	go build -o $(BUILDPATH)/terraform-provider-hivelocity

swagger-codegen-cli.jar:
	curl -o swagger-codegen-cli.jar $(SWAGGER_CODEGEN_CLI)

# Regenerate the Go API client from swagger.json. Run `make swagger` first to
# pull the latest spec from the API.
client: swagger-codegen-cli.jar
	rm -rf hivelocity-client-go
	java -jar swagger-codegen-cli.jar generate -i swagger.json -l go -o ./hivelocity-client-go

# Regenerate docs/ from the provider schema, templates/ and examples/.
# tfplugindocs downloads a Terraform binary if one is not on PATH.
docs: client
	rm -rf docs
	go tool tfplugindocs generate --provider-name hivelocity

swagger:
	curl -o swagger.json "$(HIVELOCITY_API_URL)/swagger.json?partner=1"

lint:
	golangci-lint run ./...

test:
	go test github.com/hivelocity/terraform-provider-hivelocity/hivelocity $(TESTARGS)

testacc:
	TF_ACC=1 go test github.com/hivelocity/terraform-provider-hivelocity/hivelocity -v -timeout 120m $(TESTARGS)

example_init: build
	rm -f $(EXAMPLE)/.terraform.lock.hcl
	terraform -chdir=$(EXAMPLE) init -plugin-dir ~/.terraform.d/plugins/

example_apply: example_init
	terraform -chdir=$(EXAMPLE) apply

example_destroy: example_init
	terraform -chdir=$(EXAMPLE) destroy

.PHONY: default install build client docs swagger lint test testacc example_init example_apply example_destroy
