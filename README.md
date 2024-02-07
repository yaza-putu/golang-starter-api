# Golang Starter Project
![golang](https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Go_Logo_Blue.svg/1200px-Go_Logo_Blue.svg.png)

Project layout design influenced by [standard go project layout](https://github.com/golang-standards/project-layout)
### Note : Be sure to rename the module before beginning the project.
## How to start

- install depedency
  ```bash
  make tidy
  # or
  go mod tidy
  ```
- copy environment dev, test and set app_status=test for .env.test
  ```bash
  make config
  #or
  cp .env.example .env
  cp .env.example .env.test
  ```


- generate key
  ```bash
    make key
    # or
    go run zoro.go key:generate
  ```

- run dev mode
  ```bash
    make run
  ```
- build
  ```bash
  make build
  ```

- run test
  ```bash
   make gotest
  ```

- make migration
  ```bash
   make migration table="name_of_table"
  ```
  
- run migration
  ```bash
   make migrate-up
  ```

- make seeder
   ```bash
   make seeder name="name_of_seeder"
  ```
- run seeder
  ```bash
   make seed-up
  ```
  

# Default Login
email : user@mail.com

pass : Password1

# Alasan memakai service & repository pattern
Saya biasanya memakai pattern service dan repository bertujuan memisahkan business logic dengan query logic serta memfungsikan handler untuk fokus pada pengarahan lalulintas data
kalau tidak memakai pattern cendrung kita menulis logika programing pada satu file katakan handler kalau project kecil tidak masalah tetapi jika project besar ini akan sulit di maintenace karena kode program yang banyak dan bercampur

selain pattern di project ini saya menambahkan otomatisasi pembuatan migration dan seeder untuk memudahkan development


  

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
- [Redis](https://github.com/redis/go-redis)

