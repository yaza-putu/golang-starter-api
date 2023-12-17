package goods

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	eAuth "github.com/yaza-putu/golang-starter-api/src/app/auth/entity"
	eCat "github.com/yaza-putu/golang-starter-api/src/app/category/entity"
	"github.com/yaza-putu/golang-starter-api/src/app/goods/entity"
	"github.com/yaza-putu/golang-starter-api/src/config"
	"github.com/yaza-putu/golang-starter-api/src/core"
	"github.com/yaza-putu/golang-starter-api/src/database"
	response2 "github.com/yaza-putu/golang-starter-api/src/http/response"
	"github.com/yaza-putu/golang-starter-api/src/utils"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"testing"
)

type e2eTestSuite struct {
	suite.Suite
	Token string
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	s.Require().NoError(core.EnvTesting())
	s.Require().NoError(core.DatabaseTesting())

	// run migration
	database.MigrationRegister(func(db *gorm.DB) error {
		return db.AutoMigrate(&eAuth.User{})
	}, func(db *gorm.DB) error {
		return db.Migrator().DropTable(&eAuth.User{})
	})

	// run seeder
	database.SeederRegister(func(db *gorm.DB) error {
		m := eAuth.Users{
			eAuth.User{
				ID:       utils.Uid(13),
				Name:     "User",
				Email:    "user@mail.com",
				Password: utils.Bcrypt("Password1"),
			},
		}

		return db.Create(&m).Error
	})

	database.MigrationRegister(func(db *gorm.DB) error {
		return db.AutoMigrate(&eCat.Category{})
	}, func(db *gorm.DB) error {
		return db.Migrator().DropTable(&eCat.Category{})
	})

	database.MigrationRegister(func(db *gorm.DB) error {
		return db.AutoMigrate(&entity.Goods{})
	}, func(db *gorm.DB) error {
		return db.Migrator().DropTable(&entity.Goods{})
	})

	database.MigrationUp()
	database.SeederUp()

	core.Redis()
	go core.HttpServerTesting()
	Token(s)
}

func (s *e2eTestSuite) TearDownSuite() {
	database.MigrationDown()
}

func Token(s *e2eTestSuite) {
	reqStr := `{"email":"user@mail.com","password" : "Password1"}`

	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/token", config.Host().Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	bodyToken := response2.DataApi{}
	json.NewDecoder(response.Body).Decode(&bodyToken)
	s.NoError(err)
	token := bodyToken.Data.(map[string]any)
	s.Token = token["access_token"].(string)
	defer response.Body.Close()
}

func (s *e2eTestSuite) TestValidationForm() {
	reqStr := `{"name":"", "category_id":""}`
	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/goods", config.Host().Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", s.Token))

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusUnprocessableEntity, response.StatusCode)
	defer response.Body.Close()
}

func (s *e2eTestSuite) TestTokenEmpty() {
	reqStr := `{"name":"", "category_id":""}`

	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/goods", config.Host().Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusBadRequest, response.StatusCode)

	defer response.Body.Close()
}

func (s *e2eTestSuite) TestWrongToken() {
	reqStr := `{"name":"", "category_id":""}`

	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/goods", config.Host().Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"))

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusUnauthorized, response.StatusCode)

	defer response.Body.Close()
}

func (s *e2eTestSuite) TestSuccessCreate() {
	reqStr := fmt.Sprintf(`{"name":"GD 1", "category_id": "CAT X"}`)
	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/goods", config.Host().Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", s.Token))

	client := http.Client{}

	response, err := client.Do(req)

	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	defer response.Body.Close()
}

func (s *e2eTestSuite) create(name string) (string, string) {
	reqStr := fmt.Sprintf(`{"name":"%s", "category_id": "CAT X"}`, name)
	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/goods", config.Host().Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", s.Token))

	client := http.Client{}

	response, err := client.Do(req)
	bodyToken := response2.DataApi{}
	json.NewDecoder(response.Body).Decode(&bodyToken)

	data := bodyToken.Data.(map[string]any)

	s.NoError(err)

	defer response.Body.Close()

	return data["id"].(string), data["category_id"].(string)
}

func (s *e2eTestSuite) TestSuccessUpdate() {
	id, CatId := s.create("GD 1")
	reqStr := fmt.Sprintf(`{"name":"GD 2", "category_id" : "%s"}`, CatId)
	req, err := http.NewRequest(echo.PUT, fmt.Sprintf("http://localhost:%d/api/goods/%s", config.Host().Port, id), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", s.Token))

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, response.StatusCode)

	defer response.Body.Close()
}

func (s *e2eTestSuite) TestSuccessFindById() {
	id, _ := s.create("GD 1")

	reqStr := ``
	req, err := http.NewRequest(echo.GET, fmt.Sprintf("http://localhost:%d/api/goods/%s", config.Host().Port, id), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", s.Token))

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, response.StatusCode)

	defer response.Body.Close()
}

func (s *e2eTestSuite) TestNotFoundData() {

	reqStr := ``
	req, err := http.NewRequest(echo.GET, fmt.Sprintf("http://localhost:%d/api/goods/%s", config.Host().Port, "zs"), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", s.Token))

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusNotFound, response.StatusCode)
	defer response.Body.Close()
}

func (s *e2eTestSuite) TestUpdateStock() {
	id, _ := s.create("GD 1")

	reqStr := `{"stock" : 4}`
	req, err := http.NewRequest(echo.PATCH, fmt.Sprintf("http://localhost:%d/api/goods/%s", config.Host().Port, id), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", s.Token))

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, response.StatusCode)

	defer response.Body.Close()
}
