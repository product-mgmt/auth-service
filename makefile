BINARY_NAME=$(notdir $(shell pwd))
.DEFAULT := run

dev:
	@go run cmd/main.go

build:
	@GOARCH=amd64 GOOS=darwin go build -o ./target/${BINARY_NAME}-darwin cmd/main.go
	@GOARCH=amd64 GOOS=linux go build -o ./target/${BINARY_NAME}-linux cmd/main.go
	@GOARCH=amd64 GOOS=windows go build -o ./target/${BINARY_NAME}-windows cmd/main.go

run: build
	@./target/${BINARY_NAME}-darwin

clean:
	@go clean
	@rm ./target/${BINARY_NAME}-darwin
	@rm ./target/${BINARY_NAME}-linux
	@rm ./target/${BINARY_NAME}-windows

test:
	@go test ./...

test_coverage:
	@go test ./... -coverageprofile=coverage.out

dep:
	@go mod download

vet:
	@go vet

lint:
	@golangci-lint run --enable-all