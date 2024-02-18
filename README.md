# Golang Starter Project
![golang](https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Go_Logo_Blue.svg/1200px-Go_Logo_Blue.svg.png)

Project layout design influenced by [standard go project layout](https://github.com/golang-standards/project-layout)
### Note : Be sure to rename the name of module before beginning the project.
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
  
# How to mock gorm and redis
For this template using global variable for database connection and redis connection 
  - if you need mock database in unit test you can use
 ```go
// import
import "qila.co.id/internal/database"

dbMock :=  ...

database.Mock(dbMock)
```  
special setting connection gorm in unit test (mocking) <br>
set SkipInitializeWithVersion to true 
```go
db, err := gorm.Open(mysql.New(mysql.Config{
    ...
    SkipInitializeWithVersion: true,
}), &gorm.Config{})
```

- if you need mock redis in unit test you can use
```go
// import redis client
import redis_client "qila.co.id/internal/pkg/redis"

redisMock := ...

redis_client.Mock(redisMock)
```

# Default Login
email : admin@mail.com

pass : Password1

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
ID   string `validate:"required" json:"id"`
}
```
## Stack 
- [Echo](https://echo.labstack.com)
- [Gorm](https://gorm.io)
- [Env](https://github.com/spf13/viper)
- [Redis](https://github.com/redis/go-redis)

