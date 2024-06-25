build:
	@go build -o bin/fs cmd/main.go

run: build
	@./bin/fs

test:
	@go test ./... -v -cover -race