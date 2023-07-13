build:
	@go build -o ./bin/api.exe

run: build
	@./bin/api.exe

test:
	@go test -v ./...
