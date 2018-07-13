TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

.PHONY: test

fmt:
	gofmt -w $(GOFMT_FILES)

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test -v -timeout=30s -parallel=4
