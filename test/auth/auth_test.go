package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/core"
	response2 "github.com/yaza-putu/golang-starter-api/internal/http/response"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/i18n"
)

type (
	e2eTestSuite struct {
		suite.Suite
	}
	expect struct {
		code    int
		include []string
	}
	testTable struct {
		name   string
		data   string
		expect expect
	}
)

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	s.Require().NoError(core.EnvTesting())
	i18n.New(i18n.SetLocale("en"))
}

// this function executes after all tests executed
func (s *e2eTestSuite) TearDownSuite() {
	core.EnvRollback()
}

func (s *e2eTestSuite) TestToken() {

	testTable := []testTable{
		{
			name: "create token",
			data: `{"email":"admin@mail.com","password" : "Password1"}`,
			expect: expect{
				code:    http.StatusOK,
				include: []string{"access_token"},
			},
		},
		{
			name: "credential not match",
			data: `{"email":"admin@mail.com","password" : "1"}`,
			expect: expect{
				code:    http.StatusUnauthorized,
				include: []string{},
			},
		},
		{
			name: "validation email",
			data: `{"email":"","password" : "sdsfsd"}`,
			expect: expect{
				code:    http.StatusUnprocessableEntity,
				include: []string{"Email is a required"},
			},
		},
		{
			name: "validation password",
			data: `{"email":"m@mail.com","password" : ""}`,
			expect: expect{
				code:    http.StatusUnprocessableEntity,
				include: []string{"Password is a required"},
			},
		},
	}

	for _, t := range testTable {
		s.Run(t.name, func() {
			req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/v1/token", config.Host().Port), strings.NewReader(t.data))
			s.NoError(err)

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			client := http.Client{}

			response, err := client.Do(req)
			byteBody, err := ioutil.ReadAll(response.Body)
			s.NoError(err)

			s.Equal(t.expect.code, response.StatusCode)

			if len(t.expect.include) > 0 {
				for _, include := range t.expect.include {
					assert.Contains(s.T(), strings.Trim(string(byteBody), "\n"), include)
				}
			}

			defer response.Body.Close()
		})
	}
}

func (s *e2eTestSuite) TestRenewalToken() {
	// get token
	reqTokenStr := `{"email":"admin@mail.com","password" : "Password1"}`

	reqToken, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/api/v1/token", config.Host().Port), strings.NewReader(reqTokenStr))
	s.NoError(err)

	reqToken.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}

	resToken, err := client.Do(reqToken)
	bodyToken := response2.DataApi{}
	json.NewDecoder(resToken.Body).Decode(&bodyToken)
	s.NoError(err)
	token := bodyToken.Data.(map[string]any)

	reqRefTokenStr := fmt.Sprintf(`{"device_id":"%s"}`, token["device_id"].(string))

	reqRToken, err := http.NewRequest(echo.PUT, fmt.Sprintf("http://localhost:%d/api/v1/token", config.Host().Port), strings.NewReader(reqRefTokenStr))
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

	reqRefTokenStr := fmt.Sprintf(`{"device_id":"%s"}`, "xys")

	reqRToken, err := http.NewRequest(echo.PUT, fmt.Sprintf("http://localhost:%d/api/v1/token", config.Host().Port), strings.NewReader(reqRefTokenStr))
	s.NoError(err)

	reqRToken.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client2 := http.Client{}

	resRToken, err := client2.Do(reqRToken)
	s.NoError(err)

	s.Equal(http.StatusBadRequest, resRToken.StatusCode)

	resRToken.Body.Close()
}
