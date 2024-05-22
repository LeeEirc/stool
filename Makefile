NAME=stool
BUILDDIR=build

BASEPATH := $(shell pwd)
BRANCH := $(shell git symbolic-ref HEAD 2>/dev/null | cut -d"/" -f 3)
BUILD := $(shell git rev-parse --short HEAD)


VERSION ?= $(BRANCH)-$(BUILD)
BuildTime:= $(shell date -u '+%Y-%m-%d %I:%M:%S%p')
COMMIT:= $(shell git rev-parse HEAD)
GOVERSION:= $(shell go version)
TARGETARCH ?= amd64

LDFLAGS=-w -s

BUILDCMD=CGO_ENABLED=0 go build -trimpath -ldflags "${LDFLAGS}"

CURRENT_OS_ARCH = $(shell go env GOOS)-$(shell go env GOARCH)

define make_artifact_full
	GOOS=$(1) GOARCH=$(2) $(BUILDCMD) -o $(BUILDDIR)/$(NAME)-$(1)-$(2) .
endef

build:
	$(KOKOBUILD) -o $(BUILDDIR)/$(NAME)-$(CURRENT_OS_ARCH) $(KOKOSRCFILE)


all:
	$(call make_artifact_full,darwin,amd64)
	$(call make_artifact_full,darwin,arm64)
	$(call make_artifact_full,linux,amd64)
	$(call make_artifact_full,linux,arm64)
	$(call make_artifact_full,linux,mips64le)
	$(call make_artifact_full,linux,ppc64le)
	$(call make_artifact_full,linux,s390x)
	$(call make_artifact_full,linux,riscv64)

local:
	$(call make_artifact_full,$(shell go env GOOS),$(shell go env GOARCH))

darwin-amd64:
	$(call make_artifact_full,darwin,amd64)

darwin-arm64:
	$(call make_artifact_full,darwin,arm64)

linux-amd64:
	$(call make_artifact_full,linux,amd64)

linux-arm64:
	$(call make_artifact_full,linux,arm64)

linux-loong64:
	$(call make_artifact_full,linux,loong64)

linux-mips64le:
	$(call make_artifact_full,linux,mips64le)

linux-ppc64le:
	$(call make_artifact_full,linux,ppc64le)

linux-s390x:
	$(call make_artifact_full,linux,s390x)

linux-riscv64:
	$(call make_artifact_full,linux,riscv64)


.PHONY: clean
clean:
	-rm -rf $(BUILDDIR)
