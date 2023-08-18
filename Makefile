build:
	GOOS=linux GOARCH=amd64 go build -o bin/main main.go

run:
	go run main.go

test:
	go test -race ./...

config:
	cp .sample.yml config.yml

tidy:
	go mod tidy
