reports := $(shell pwd)/reports
packages := $(shell go list ./...)

test:
	@ go test -v -race ./...

benchmark:
	@ go test -run=none -bench=. -benchmem ./...

coverage:
	@ mkdir -p $(reports) && \
	  echo "mode: atomic" > $(reports)/cover.out
	@ $(foreach package, $(packages), \
	    go test -covermode=atomic -coverprofile=cover.out $(package) || exit 1; \
	    tail -n +2 cover.out >> $(reports)/cover.out;)
	@ go tool cover -html=$(reports)/cover.out -o $(reports)/cover.html && \
	  rm cover.out $(reports)/cover.out

.PHONY: test benchmark coverage
