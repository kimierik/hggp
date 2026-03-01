go build -C ./frontend -o ./build/frontend
go build -C ./backend -o ./build/backend

GOOS=js GOARCH=wasm go build \
    -C ./frontend \
    -o static/wasm/app.wasm \
    -trimpath \
    -ldflags="-s -w" \
    wasm.go
