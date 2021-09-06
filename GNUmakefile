ifeq ($(GOPATH),)
	GOPATH:=$(shell go env GOPATH)
endif

default: build

build:
	go install

test:
	go test github.com/hivelocity/terraform-provider-hivelocity/hivelocity

testacc:
	TF_ACC=1 go test github.com/hivelocity/terraform-provider-hivelocity/hivelocity
