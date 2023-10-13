package services

import (
	"context"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
	"time"
	"yaza/src/app/auth/entities"
	"yaza/src/app/auth/repositories"
	"yaza/src/utils"
)

type (
	mockUserRepository struct {
		mock.Mock
		repositories.UserInterface
	}

	mockTokenService struct {
		mock.Mock
		TokenInterface
	}
)

func (m *mockTokenService) Create(ctx context.Context, user entities.User) (string, string, error) {
	args := m.Called(ctx, user)

	return args.String(0), args.String(0), nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, mail string) (entities.User, error) {
	args := m.Called(ctx, mail)

	return args.Get(0).(entities.User), nil
}

func TestLogin(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	objUserRepository := new(mockUserRepository)
	objTokenService := new(mockTokenService)

	objUserRepository.On("FindByEmail", ctx, "admin@mail.com").Return(entities.User{ID: "xys", Name: "admin", Email: "admin@mail.com", Password: utils.Bcrypt("Password1")}, nil)
	// run method mocked for QA
	usr, _ := objUserRepository.FindByEmail(ctx, "admin@mail.com")

	objTokenService.On("Create", ctx, usr).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwuY29tIiwibmFtZSI6ImFkbWluIn0.QjV0U4Qy-5hfOc60dg997lyb1_sfKGtxx2aeO91L_sg", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwuY29tIiwibmFtZSI6ImFkbWluIiwidG9rZW4iOiJleUpoYkdjaU9pSklVekkxTmlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKbGJXRnBiQ0k2SW1Ga2JXbHVRRzFoYVd3dVkyOXRJaXdpYm1GdFpTSTZJbUZrYldsdUluMC5RalYwVTRReS01aGZPYzYwZGc5OTdseWIxX3NmS0d0eHgyYWVPOTFMX3NnIn0.QS16oXYsaYqZS7m-LnZVRDMhOwbEDJ8mEoxW80zl4yk", nil)
	// run method mocked for QA
	objTokenService.Create(ctx, usr)
	// test mock object
	objTokenService.AssertExpectations(t)
	objUserRepository.AssertExpectations(t)

	a := NewAuthService(objUserRepository, objTokenService)

	r := a.Login(ctx, "admin@mail.com", "Password1")
	token := r.Data

	assert.Equal(t, r.Code, 200)
	assert.Equal(t, len(strings.Split(token.(map[string]string)["access_token"], ".")), 3)
}

func TestLoginWithTimout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()

	objUserRepository := new(mockUserRepository)
	objTokenService := new(mockTokenService)

	objUserRepository.On("FindByEmail", ctx, "admin@mail.com").Return(entities.User{ID: "xys", Name: "admin", Email: "admin@mail.com", Password: utils.Bcrypt("Password1")}, nil)
	// run method mocked for QA
	usr, _ := objUserRepository.FindByEmail(ctx, "admin@mail.com")

	objTokenService.On("Create", ctx, usr).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwuY29tIiwibmFtZSI6ImFkbWluIn0.QjV0U4Qy-5hfOc60dg997lyb1_sfKGtxx2aeO91L_sg", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwuY29tIiwibmFtZSI6ImFkbWluIiwidG9rZW4iOiJleUpoYkdjaU9pSklVekkxTmlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKbGJXRnBiQ0k2SW1Ga2JXbHVRRzFoYVd3dVkyOXRJaXdpYm1GdFpTSTZJbUZrYldsdUluMC5RalYwVTRReS01aGZPYzYwZGc5OTdseWIxX3NmS0d0eHgyYWVPOTFMX3NnIn0.QS16oXYsaYqZS7m-LnZVRDMhOwbEDJ8mEoxW80zl4yk", nil)
	// run method mocked for QA
	objTokenService.Create(ctx, usr)
	// test mock object
	objTokenService.AssertExpectations(t)
	objUserRepository.AssertExpectations(t)

	a := NewAuthService(objUserRepository, objTokenService)

	r := a.Login(ctx, "admin@mail.com", "Password1")

	assert.Equal(t, r.Code, 408)
}

func TestLoginFailed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	objUserRepository := new(mockUserRepository)
	objTokenService := new(mockTokenService)

	objUserRepository.On("FindByEmail", ctx, "admin@mail.com").Return(entities.User{ID: "xys", Name: "admin", Email: "admin@mail.com", Password: utils.Bcrypt("Password1")}, nil)
	// run method mocked for QA
	usr, _ := objUserRepository.FindByEmail(ctx, "admin@mail.com")

	objTokenService.On("Create", ctx, usr).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwuY29tIiwibmFtZSI6ImFkbWluIn0.QjV0U4Qy-5hfOc60dg997lyb1_sfKGtxx2aeO91L_sg", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwuY29tIiwibmFtZSI6ImFkbWluIiwidG9rZW4iOiJleUpoYkdjaU9pSklVekkxTmlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKbGJXRnBiQ0k2SW1Ga2JXbHVRRzFoYVd3dVkyOXRJaXdpYm1GdFpTSTZJbUZrYldsdUluMC5RalYwVTRReS01aGZPYzYwZGc5OTdseWIxX3NmS0d0eHgyYWVPOTFMX3NnIn0.QS16oXYsaYqZS7m-LnZVRDMhOwbEDJ8mEoxW80zl4yk", nil)
	// run method mocked for QA
	objTokenService.Create(ctx, usr)
	// test mock object
	objTokenService.AssertExpectations(t)
	objUserRepository.AssertExpectations(t)

	a := NewAuthService(objUserRepository, objTokenService)

	r := a.Login(ctx, "admin@mail.com", "wrong")

	assert.Equal(t, r.Code, 401)
}
