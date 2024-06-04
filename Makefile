build:
	GOOS=linux GOARCH=amd64 go build -o build/main cmd/main.go

serve:
	go run cmd/main.go
	
server:
	go run cmd/main.go

config:
	cp .env.example .env && cp .env.example .env.test

tidy:
	go mod tidy
key:
	go run cmd/zoro.go key:generate

init:
	go run cmd/zoro.go configure:module ${module}

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