default: dev

dev:
    go run main.go

wasmrun:
    go run github.com/hajimehoshi/wasmserve@latest .

build-wasm:
    env GOOS=js GOARCH=wasm go build -o yourgame.wasm mygame

copy-js:
    cp $(go env GOROOT)/lib/wasm/wasm_exec.js

test-server:
    uv run python -m http.server
