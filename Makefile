SHELL := $(shell which bash) # set default shell
# OS / Arch we will build our binaries for
OSARCH ?= "linux/amd64 linux/386 darwin/amd64 darwin/386"
ENV = /usr/bin/env

.SHELLFLAGS = -c # Run commands in a -c flag
.SILENT: ; # no need for @
.ONESHELL: ; # recipes execute in same shell
.NOTPARALLEL: ; # wait for this target to finish
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

.PHONY: all # All targets are accessible for user
.DEFAULT: help # Running Make will run the help target

help: ## Show Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dep: ## Get build dependencies
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/mitchellh/gox

cross-build: clean dep lint ## Build the app for multiple os/arch
	dep ensure && gox -osarch=$(OSARCH) -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

clean: ## Clean the dist directory
	rm -rf dist/

lint: dep ## Lint the code
	gometalinter --deadline=10m --config=.gometalinter.json --exclude=^vendor\/ ./...
