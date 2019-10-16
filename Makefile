REPO_PATH=github.com/echoturing/luxury

GOFMT_FILES=$(shell find . -name '*.go' | grep -v vendor | xargs)

.PHONY: fmt
fmt:
	@goimports -l -local "$(REPO_PATH)" -w $(GOFMT_FILES)

.PHONY: build
build:
	@go build cmd/client/luxury_client.go