.PHONY: format
format:
	$(call format)

.PHONY: generate.mock
generate.mock:
	@go generate ./domain/...
	$(call format)

define format
	@go run mvdan.cc/gofumpt -l -w .
	@go run golang.org/x/tools/cmd/goimports -w .
	@go mod tidy
endef
