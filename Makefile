name := gocert


build:
	@ cherry build -cross-compile=false

build-all:
	@ cherry build -cross-compile=true

test:
	@ go test -race ./...

test-short:
	@ go test -short ./...

coverage:
	@ go test -covermode=atomic -coverprofile=c.out ./...
	@ go tool cover -html=c.out -o coverage.html


.PHONY: build build-all
.PHONY: test test-short coverage
