REPO_PATH=github.com/echoturing/luxury

GOFMT_FILES=$(shell find . -name '*.go' | grep -v vendor | xargs)

.PHONY: fmt
fmt:
	@goimports -l -local "$(REPO_PATH)" -w $(GOFMT_FILES)

.PHONY: build
build:
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build -o ./.build/luxury_windows.exe cmd/client/luxury_client.go
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64  go build -o ./.build/luxury_mac cmd/client/luxury_client.go