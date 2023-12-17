package auth

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/yaza-putu/golang-starter-api/src/app/auth/entity"
	"github.com/yaza-putu/golang-starter-api/src/config"
	"github.com/yaza-putu/golang-starter-api/src/core"
	"github.com/yaza-putu/golang-starter-api/src/database"
	response2 "github.com/yaza-putu/golang-starter-api/src/http/response"
	"github.com/yaza-putu/golang-starter-api/src/utils"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type e2eTestSuite struct {
	suite.Suite
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	s.Require().NoError(core.EnvTesting())
	s.Require().NoError(core.DatabaseTesting())
	// run migration
	database.MigrationRegister(func(db *gorm.DB) error {
		return db.AutoMigrate(&entity.User{})
	}, func(db *gorm.DB) error {
		return db.Migrator().DropTable(&entity.User{})
	})

	// run seeder
	database.SeederRegister(func(db *gorm.DB) error {
		m := entity.Users{
			entity.User{
				ID:       utils.Uid(13),
				Name:     "User",
				Email:    "user@mail.com",
				Password: utils.Bcrypt("Password1"),
			},
		}

		return db.Create(&m).Error
	})

	database.MigrationUp()
	database.SeederUp()

	go core.HttpServerTesting()
}

func (s *e2eTestSuite) TearDownSuite() {
	database.MigrationDown()
}

func (s *e2eTestSuite) TestCreateToken() {
	reqStr := `{"email":"user@mail.com","password" : "Password1"}`

	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/token", config.Host().Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}

	response, err := client.Do(req)
	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)

	s.Equal(http.StatusOK, response.StatusCode)
	assert.Contains(s.T(), strings.Trim(string(byteBody), "\n"), "access_token")

	defer response.Body.Close()
}

func (s *e2eTestSuite) TestWrongCredintial() {
	reqStr := `{"email":"user@mail.com","password" : "1"}`

	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/token", config.Host().Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusUnauthorized, response.StatusCode)

	defer response.Body.Close()
}

func (s *e2eTestSuite) TestValidationPassword() {
	reqStr := `{"email":"user@mail.com","password" : ""}`

	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/token", config.Host().Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusUnprocessableEntity, response.StatusCode)

	defer response.Body.Close()
}

func (s *e2eTestSuite) TestValidationEmail() {
	reqStr := `{"email":"","password" : "Password1"}`

	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/token", config.Host().Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusUnprocessableEntity, response.StatusCode)

	defer response.Body.Close()
}

func (s *e2eTestSuite) TestRenewalToken() {
	// get token
	reqTokenStr := `{"email":"user@mail.com","password" : "Password1"}`

	reqToken, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/token", config.Host().Port), strings.NewReader(reqTokenStr))
	s.NoError(err)

	reqToken.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}

	resToken, err := client.Do(reqToken)
	bodyToken := response2.DataApi{}
	json.NewDecoder(resToken.Body).Decode(&bodyToken)
	s.NoError(err)
	token := bodyToken.Data.(map[string]any)

	reqRefTokenStr := fmt.Sprintf(`{"refresh_token":"%s"}`, token["refresh_token"].(string))

	reqRToken, err := http.NewRequest(echo.PUT, fmt.Sprintf("http://localhost:%d/api/token", config.Host().Port), strings.NewReader(reqRefTokenStr))
	s.NoError(err)

	reqRToken.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client2 := http.Client{}

	resRToken, err := client2.Do(reqRToken)
	s.NoError(err)

	s.Equal(http.StatusOK, resRToken.StatusCode)

	resRToken.Body.Close()
	resToken.Body.Close()
}

func (s *e2eTestSuite) TestFailedRenewalToken() {
	// get token
	reqTokenStr := `{"email":"user@mail.com","password" : "Password1"}`

	reqToken, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/token", config.Host().Port), strings.NewReader(reqTokenStr))
	s.NoError(err)

	reqToken.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}

	resToken, err := client.Do(reqToken)
	bodyToken := response2.DataApi{}
	json.NewDecoder(resToken.Body).Decode(&bodyToken)
	s.NoError(err)
	token := bodyToken.Data.(map[string]any)

	reqRefTokenStr := fmt.Sprintf(`{"refresh_token":"%s"}`, token["access_token"].(string))

	reqRToken, err := http.NewRequest(echo.PUT, fmt.Sprintf("http://localhost:%d/api/token", config.Host().Port), strings.NewReader(reqRefTokenStr))
	s.NoError(err)

	reqRToken.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client2 := http.Client{}

	resRToken, err := client2.Do(reqRToken)
	s.NoError(err)

	s.Equal(http.StatusInternalServerError, resRToken.StatusCode)

	resRToken.Body.Close()
	resToken.Body.Close()
}
