GITHUB_USER          ?= jlorgal
GITHUB_REPO          ?= odor
GITHUB_API           ?= 
GITHUB_TOKEN         ?= 

PROJECT_NAME         ?= odor
DOCKER_REGISTRY_AUTH ?=
DOCKER_REGISTRY      ?= dockerhub.hi.inet
DOCKER_ORG           ?= awazza
DOCKER_API_VERSION   ?=

PRODUCT_VERSION      ?= $(get_version)
PRODUCT_REVISION     ?= $(get_revision)
BUILD_VERSION        ?= $(PRODUCT_VERSION)-$(PRODUCT_REVISION)
LDFLAGS              ?= -X main.Version=$(BUILD_VERSION)
DOCKER_IMAGE         ?= $(if $(DOCKER_REGISTRY),$(DOCKER_REGISTRY)/$(DOCKER_ORG)/$(PROJECT_NAME),$(DOCKER_ORG)/$(PROJECT_NAME))
PACKAGES             := $(shell go list ./... | grep -v /vendor/)

# Get the environment and import the settings.
# If the make target is pipeline-xxx, the environment is obtained from the target.
ifeq ($(patsubst pipeline-%,%,$(MAKECMDGOALS)),$(MAKECMDGOALS))
	ENVIRONMENT ?= pull
else
	override ENVIRONMENT := $(patsubst pipeline-%,%,$(MAKECMDGOALS))
endif
include delivery/env/$(ENVIRONMENT)

define help
Usage: make <command>
Commands:
  help:            Show this help information
  dep:             Ensure dependencies with dep tool
  build:           Build the application
  test-acceptance: Pass component tests
  release:         Create a new release (tag and release notes)
  run:             Launch the service with docker-compose (for testing purposes)
  clean:           Clean the project
  pipeline-pull:   Launch pipeline to handle a pull request
  pipeline-dev:    Launch pipeline to handle the merge of a pull request
  pipeline:        Launch the pipeline for the selected environment
  develenv-up:     Launch the development environment with a docker-compose of the service
  develenv-sh:     Access to a shell of a launched development environment
  develenv-down:   Stop the development environment
endef
export help

.PHONY: help dep build-deps build-config build test-acceptance release-deps release run clean \
		pipeline-pull pipeline-dev pipeline \
		develenv-up develenv-sh develenv-down

help:
	@echo "$$help"

dep:
	$(info) "Installing golang dependencies"
	dep ensure

build-deps:
	$(info) "Installing golang build dependencies"
	go get -v github.com/golang/lint/golint github.com/golang/dep/cmd/dep

build-config:
	$(info) "Copying configuration"
	mkdir -p build/bin
	cp odor/cmd/odor/config.json build/bin/

build: build-config build-deps dep
	$(info) "Building version: $(BUILD_VERSION)"
	GOBIN=$$PWD/build/bin/ go install -ldflags="$(LDFLAGS)" ./...
	go fmt $(PACKAGES)
	go vet $(PACKAGES)
	golint $(PACKAGES)
	go test $(PACKAGES)

test-acceptance:
	$(info) "Passing acceptance tests"

release-deps:
	$(info) "Installing golang release dependencies"
	go get github.com/aktau/github-release

release: release-deps
ifeq ($(RELEASE),true)
	$(info) "Creating release: $(PRODUCT_VERSION)"
	GITHUB_API="$(GITHUB_API)" GITHUB_TOKEN="$(GITHUB_TOKEN)" github-release release \
		--user $(GITHUB_USER) \
		--repo $(GITHUB_REPO) \
		--tag $(PRODUCT_VERSION) \
		--name $(PRODUCT_VERSION) \
		--description "$(get_release_notes)"
endif

run: build
	$(info) "Launching the service"
	cd build/bin && ./odor 

clean:
	$(info) "Cleaning the project"
	go clean
	rm -rf build/ vendor/

pipeline-pull: build test-acceptance
	$(info) "Completed successfully pipeline-pull"

pipeline-dev:  build test-acceptance release
	$(info) "Completed successfully pipeline-dev"

pipeline:      pipeline-$(ENVIRONMENT)

develenv-up:
	$(info) "Launching the development environment"
	docker-compose -p develenv -f delivery/docker/dev/docker-compose.yml build
	docker-compose -p develenv -f delivery/docker/dev/docker-compose.yml up -d

develenv-sh:
	docker exec -it develenv_odor_1 bash

develenv-down:
	$(info) "Shutting down the development environment"
	docker-compose -p develenv -f delivery/docker/dev/docker-compose.yml down

# Functions
info := @printf "\033[32;01m%s\033[0m\n"
get_version  := $$(delivery/scripts/github.sh get_version)
get_revision := $$(delivery/scripts/github.sh get_revision)
get_release_notes := $$(delivery/scripts/github.sh get_release_notes)
