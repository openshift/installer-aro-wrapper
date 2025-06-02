SHELL = /bin/bash
TAG ?= $(shell git describe --exact-match 2>/dev/null)
COMMIT = $(shell git rev-parse --short=7 HEAD)$(shell [[ $$(git status --porcelain) = "" ]] || echo -dirty)
ARO_IMAGE_BASE = $(RP_IMAGE_ACR).azurecr.io/aroinstaller

ifneq ($(shell uname -s),Darwin)
    export CGO_CFLAGS=-Dgpgme_off_t=off_t
endif

ifeq ($(TAG),)
	VERSION = $(COMMIT)
else
	VERSION = $(TAG)
endif

# default to registry.access.redhat.com for build images on local builds and CI builds without $RP_IMAGE_ACR set.
ifeq ($(RP_IMAGE_ACR),arointsvc)
	REGISTRY = arointsvc.azurecr.io
else ifeq ($(RP_IMAGE_ACR),arosvc)
	REGISTRY = arosvc.azurecr.io
else
	REGISTRY = registry.access.redhat.com
endif

ARO_IMAGE ?= $(ARO_IMAGE_BASE):$(VERSION)

include .bingo/Variables.mk

build-all:
	go build -tags altinfra,fipscapable,aro,containers_image_openpgp ./...

aro: generate
	go build -tags altinfra,fipscapable,aro,containers_image_openpgp,codec.safe ./cmd/aro

clean:
	rm -rf aro
	find -type d -name 'gomock_reflect_[0-9]*' -exec rm -rf {} \+ 2>/dev/null

generate: install-tools
	go generate ./...

image-aro:
	docker pull $(REGISTRY)/ubi8/ubi-minimal
	docker build --network=host --no-cache -f Dockerfile.aro -t $(ARO_IMAGE) --build-arg REGISTRY=$(REGISTRY) .

publish-image-aro: image-aro
	docker push $(ARO_IMAGE)
ifeq ("${RP_IMAGE_ACR}-$(BRANCH)","arointsvc-master")
		docker tag $(ARO_IMAGE) arointsvc.azurecr.io/aroinstaller:latest
		docker push arointsvc.azurecr.io/aroinstaller:latest
endif

test-go: generate build-all validate-go lint-go unit-test-go

validate-go: $(GOIMPORTS)
	gofmt -s -w cmd hack pkg test
	$(GOIMPORTS) -w -local=github.com/openshift/installer-aro-wrapper cmd hack pkg test
	go run ./hack/validate-imports cmd hack pkg test
	go run ./hack/licenses
	@[ -z "$$(ls pkg/util/*.go 2>/dev/null)" ] || (echo error: go files are not allowed in pkg/util, use a subpackage; exit 1)
	@[ -z "$$(find -name "*:*")" ] || (echo error: filenames with colons are not allowed on Windows, please rename; exit 1)
	go vet -tags containers_image_openpgp ./...

validate-go-action:
	go run ./hack/licenses -validate -ignored-go vendor,pkg/client,.git -ignored-python python/client,vendor,.git
	go run ./hack/validate-imports cmd hack pkg test
	@[ -z "$$(ls pkg/util/*.go 2>/dev/null)" ] || (echo error: go files are not allowed in pkg/util, use a subpackage; exit 1)
	@[ -z "$$(find -name "*:*")" ] || (echo error: filenames with colons are not allowed on Windows, please rename; exit 1)

.PHONY: validate-fips
validate-fips: $(BINGO)
	GOFLAGS="-mod=mod" $(BINGO) get -l fips-detect
	GOFLAGS="-mod=mod" $(BINGO) get -l gojq
	hack/fips/validate-fips.sh ./aro

unit-test-go: $(GOTESTSUM)
	$(GOTESTSUM) --format pkgname --junitfile report.xml -- -tags=altinfra,aro,containers_image_openpgp -coverprofile=cover.out ./...

lint-go: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --verbose

vendor:
	# See comments in the script for background on why we need it
	hack/update-go-module-dependencies.sh

.PHONY: install-tools
install-tools: $(BINGO)
	GOFLAGS="-mod=mod" $(BINGO) get -l
# Fixes https://github.com/uber-go/mock/issues/185 for MacOS users
ifeq ($(shell uname -s),Darwin)
	codesign -f -s - ${GOPATH}/bin/mockgen
endif
