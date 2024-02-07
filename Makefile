build:
	GOOS=linux GOARCH=amd64 go build -o build/main cmd/main.go

run:
	go run cmd/main.go

config:
	cp .env.example .env && cp .env.example .env.test

tidy:
	go mod tidy
key:
	go run cmd/zoro.go key:generate

gotest:
	go test ./test/...

migration:
	go run cmd/zoro.go make:migration ${table}

migrate-up:
	go run cmd/zoro.go migrate:up

migrate-down:
	go run cmd/zoro.go migrate:down

seed-up:
	go run cmd/zoro.go seed:up

seeder:
	go run cmd/zoro.go make:seeder ${name}