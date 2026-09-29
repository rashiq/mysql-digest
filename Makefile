.PHONY: build

SOURCE ?= ../mysql-digest

build:
	cd "$(SOURCE)" && GOOS=js GOARCH=wasm go build -trimpath -o "$(CURDIR)/digest.wasm" ./cmd/wasm
	cp "$$(cd "$(SOURCE)" && go env GOROOT)/lib/wasm/wasm_exec.js" wasm_exec.js
