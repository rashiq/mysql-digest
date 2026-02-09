.PHONY: build

build:
	git checkout main
	GOOS=js GOARCH=wasm go build -o /tmp/digest.wasm ./cmd/wasm/
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" /tmp/wasm_exec.js
	git checkout gh-pages
	mv /tmp/digest.wasm .
	mv /tmp/wasm_exec.js .
	@echo "✓ Built digest.wasm and wasm_exec.js"
