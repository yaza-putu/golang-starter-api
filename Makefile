build:
	GOOS=linux GOARCH=amd64 go build -o bin/main main.go

run:
	go run main.go

config:
	cp .env.example .env

tidy:
	go mod tidy
key:
	go run zoro.go key:generate

gotest:
	go test ./test/...