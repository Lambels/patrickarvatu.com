build-backend: generate redis build-go
test-backend: generate redis test-go

redis:
	redis-server

generate:
	go generate ./...

test-go:
	go test -v ./...

build-go:
	go build -o patrickarvatu ./cmd