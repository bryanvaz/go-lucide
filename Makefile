OUT_PATH=templ

all: sync build

deps:
	go mod tidy

.PHONY: sync
sync:
	@bash scripts/sync_lucide.sh

.PHONY: build
build:
	@go run ./scripts/build_packages

.PHONY: clean
clean:
	@rm -rf lucide/*
	@rm -rf $(OUT_PATH)/*.go
	@rm -rf $(OUT_PATH)/*.templ

