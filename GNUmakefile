ifeq ($(GOPATH),)
	GOPATH:=$(shell go env GOPATH)
endif

default: build

build:
	go install
