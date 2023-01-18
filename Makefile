fmt:
	go fmt ./...

vet:
	go vet ./...

run: fmt vet
	go run ./cmd/api/main.go

test: fmt vet
	go test ./...