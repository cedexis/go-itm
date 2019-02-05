TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

.PHONY: test

fmt:
	gofmt -w $(GOFMT_FILES)

test:
	go test ./...
