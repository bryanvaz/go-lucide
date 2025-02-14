OUT_PATH=templ

all: build

deps:
	go mod tidy

.PHONY: build
build:
	@go run ./scripts/build_packages

.PHONY: test
test:
	@cd ./packages/go-templ-lucide-icons && go test -v ./test

.PHONY: clean
clean:
	@rm -rf lucide/*
	@rm -rf $(OUT_PATH)/*.go
	@rm -rf $(OUT_PATH)/*.templ

