fmt:
	go fmt ./...

vet:
	go vet ./...

run: fmt vet
	go run ./cmd/api/main.go

clear:
	go clean -testcache

test: fmt vet clear
	go test ./... -v