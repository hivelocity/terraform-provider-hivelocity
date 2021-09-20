ifeq ($(HIVELOCITY_API_URL),)
	HIVELOCITY_API_URL:=https://core.hivelocity.net/api/v2
endif

ifeq ($(GOPATH),)
	GOPATH:=$(shell go env GOPATH)
endif

ifeq ($(GOOS),)
	GOOS:=$(shell go env GOHOSTOS)
endif

ifeq ($(GOARCH),)
	GOARCH:=$(shell go env GOHOSTARCH)
endif

OS_ARCH:=$(GOOS)_$(GOARCH)
BUILDPATH:=$(HOME)/.terraform.d/plugins/registry.terraform.io/hivelocity/hivelocity/0.1.0/$(OS_ARCH)
SWAGGER_CODEGEN_CLI:=https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/2.4.15/swagger-codegen-cli-2.4.15.jar

default: build

install:
	go install

build: client
	go build -o $(BUILDPATH)/terraform-provider-hivelocity

swagger-codegen-cli.jar:
	curl -o swagger-codegen-cli.jar $(SWAGGER_CODEGEN_CLI)

client: swagger-codegen-cli.jar
	rm -rf hivelocity-client-go
	java -jar swagger-codegen-cli.jar generate -i swagger.json -l go -o ./hivelocity-client-go

docs: build
	rm -rf docs
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

swagger:
	curl -o swagger.json $(HIVELOCITY_API_URL)/swagger.json?partner=1

test:
	go test github.com/hivelocity/terraform-provider-hivelocity/hivelocity

testacc:
	TF_ACC=1 go test github.com/hivelocity/terraform-provider-hivelocity/hivelocity
