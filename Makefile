name := gocert


clean:
	@ rm -rf bin coverage $(name)

run:
	@ go run main.go

build:
	@ cherry build -cross-compile=false

build-all:
	@ cherry build -cross-compile=true

test:
	@ go test -race ./...

test-short:
	@ go test -short ./...

coverage:
	@ cherry test


.PHONY: clean
.PHONY: run build build-all
.PHONY: test test-short coverage
