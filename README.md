# Golang Starter API
![golang](https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Go_Logo_Blue.svg/1200px-Go_Logo_Blue.svg.png)

## How to start
- install depedency
  ```bash
    make tidy
  ```
- copy environment
  ```bash
    make config
  ```
- run dev mode
  ```bash
    make run
  ```
- build
  ```bash
  make build
  ```
- test
  ```bash
    make test
  ```

## Zoro command
- make migration
  ```bash
    go run zoro.go make:migration file_name
  ```
- migration up
  ```bash
    go run zoro.go migrate:up
  ```
- migration down
  ```bash
    go run zoro.go migrate:down
  ```
  
## Validation Unique With Struct Tag
- unique
```go
type v struct {
	Name string `validate:"unique=table_name:column_name"`
}
// ecample
type v struct {
Name string `validate:"unique=users:name"`
}
```
- unique with ignore
```go
type v struct {
Name string `validate:"unique=table_name:column_name:ignore_with_field_name"`
ID   string `validate:"required"`
}
// example
type v struct {
Name string `validate:"unique=users:name:ID"`
ID   string `validate:"required"`
}
```
## Stack 
- [Echo](https://echo.labstack.com)
- [Gorm](https://gorm.io)
- [Env](https://github.com/spf13/viper)
