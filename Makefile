path := $(shell pwd)
reports := $(path)/reports

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


clean:
	@ rm -rf gocert reports

run:
	@ go run main.go

build:
	@ go build $(ldflags)

install:
	@ go install $(ldflags)

test:
	@ go test -v -race ./...

test-short:
	@ go test -v -race -short ./...

coverage:
	@ mkdir -p $(reports) && \
	  echo "mode: atomic" > $(reports)/cover.out
	@ $(foreach package, $(go_packages), \
	    go test -covermode=atomic -coverprofile=cover.out $(package) || exit 1; \
	    tail -n +2 cover.out >> $(reports)/cover.out;)
	@ go tool cover -html=$(reports)/cover.out -o $(reports)/cover.html && \
	  rm cover.out $(reports)/cover.out

release:
	@ ./release.sh patch

release-minor:
	@ ./release.sh minor

release-major:
	@ ./release.sh major


.PHONY: clean
.PHONY: run build install
.PHONY: test test-short coverage
.PHONY: release release-minor release-major
