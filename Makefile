OUT_PATH=templ
VERSION_FILE = dist/go-templ-lucide-icons/VERSION
VERSION := $(shell cat $(VERSION_FILE))

deps:
	go mod tidy

.PHONY: build
build:
	@go run ./scripts/build_packages

.PHONY: test
test:
	@cd ./dist/go-templ-lucide-icons && go test -v ./test

.PHONY: commit
commit:
	@cd ./dist/go-templ-lucide-icons && git add . && git commit -m "chore: update icons to $(VERSION)" -m "Based on lucide@v$(VERSION). See https://github.com/lucide-icons/lucide/tree/$(VERSION)"

.PHONY: publish
publish:
	@cd ./dist/go-templ-lucide-icons && git tag v$(VERSION) && git push origin v$(VERSION)
	@cd ./dist/go-templ-lucide-icons && git push origin main
	@cd ./dist/go-templ-lucide-icons && GOPROXY=proxy.golang.org go list -m github.com/bryanvaz/go-templ-lucide-icons@v$(VERSION)
	@cd ./dist/go-templ-lucide-icons && gh release create -d -t v$(VERSION) --notes-from-tag v$(VERSION)

.PHONY: clean
clean:
	@rm -rf dist/*

