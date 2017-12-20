path := $(shell pwd)
binary := gocert
build_dir := $(path)/artifacts
report_dir := $(path)/report_dir

go_packages := $(shell go list ./... | grep -v //)
version_package := $(shell go list ./version)

version := $(shell head -n1 version/VERSION)
revision := $(shell git rev-parse --short HEAD)
branch := $(shell git rev-parse --abbrev-ref HEAD)
buildtime := $(shell date -u "+%Y-%m-%dT%H:%M:%SZ%z")

flag_version := -X $(version_package).Version=$(version)
flag_revision := -X $(version_package).Revision=$(revision)
flag_branch := -X $(version_package).Branch=$(branch)
flag_buildtime := -X $(version_package).BuildTime=$(buildtime)
ldflags := -ldflags "$(flag_version) $(flag_revision) $(flag_branch) $(flag_buildtime)"

platforms := linux-386 linux-amd64 darwin-386 darwin-amd64 windows-386 windows-amd64


define cross_compile
	GOOS=$(shell echo $(1) | cut -d- -f1) \
	GOARCH=$(shell echo $(1) | cut -d- -f2) \
	go build $(ldflags) -o $(build_dir)/$(binary)-$(1);
	printf "\033[1;32m ✓ \033[1;35m $(binary)-$(1) \033[0m\n";
endef


clean:
	@ rm -rf $(binary) $(build_dir) $(report_dir)

run:
	@ go run main.go

install:
	@ go install $(ldflags)

build:
	@ go build $(ldflags)

build-all:
	@ mkdir -p $(build_dir)
	@ $(foreach platform, $(platforms), $(call cross_compile,$(platform)))

test:
	@ go test -v -race ./...

test-short:
	@ go test -v -race -short ./...

coverage:
	@ mkdir -p $(report_dir) && \
	  echo "mode: atomic" > $(report_dir)/cover.out
	@ $(foreach package, $(go_packages), \
	    go test -covermode=atomic -coverprofile=cover.out $(package) || exit 1; \
	    tail -n +2 cover.out >> $(report_dir)/cover.out;)
	@ go tool cover -html=$(report_dir)/cover.out -o $(report_dir)/cover.html && \
	  rm cover.out $(report_dir)/cover.out


.PHONY: clean
.PHONY: run install build build-all
.PHONY: test test-short coverage
