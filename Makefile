build-backend: redis build-go
test-backend: redis test-go

redis:
	redis-server

test-go:
	go test -v ./...

build-go:
	go build -o patrickarvatu ./cmd