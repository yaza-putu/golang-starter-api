package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
	"strings"
	"time"
)

// Token / **************************************************************
type Token interface {
	Create(ctx context.Context, user entity.User) (string, string, error)
	Refresh(ctx context.Context, rToken string) (string, error)
	generateToken(claims jwt.Claims) (string, error)
	generateRefreshToken(claims jwt.Claims) (string, error)
}

type (
	tokenService   struct{}
	jwtTokenClaims struct {
		Email string `json:"email"`
		jwt.StandardClaims
	}
	jwtRefreshClaim struct {
		Email    string `json:"email"`
		OldToken string `json:"old_token"`
		jwt.StandardClaims
	}
	ResponseChannel struct {
		Token   string
		Refresh string
		Error   error
	}
)

// NewToken is constructor
func NewToken() Token {
	return &tokenService{}
}

// Create token
// return arg1 string token
// return arg2 string refresh token
// return arg3 error error
func (t *tokenService) Create(ctx context.Context, user entity.User) (string, string, error) {
	rc := make(chan ResponseChannel)
	go func() {
		// create token
		token, err := t.generateToken(&jwtTokenClaims{
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
			},
		})
		if err != nil {
			logger.New(err)
			rc <- ResponseChannel{Token: "", Refresh: "", Error: err}
		}

		// gen refresh token
		refresh, err := t.generateRefreshToken(&jwtRefreshClaim{
			Email:    user.Email,
			OldToken: token,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		})
		if err != nil {
			logger.New(err)
			rc <- ResponseChannel{Token: "", Refresh: "", Error: err}
		}

		rc <- ResponseChannel{Token: token, Refresh: refresh, Error: nil}

	}()

	for {
		select {
		case <-ctx.Done():
			return "", "", errors.New("request time out")
		case res := <-rc:
			return res.Token, res.Refresh, res.Error
		}
	}
}

func (t *tokenService) Refresh(ctx context.Context, rToken string) (string, error) {
	rc := make(chan ResponseChannel)

	go func() {
		if !strings.Contains(rToken, ".") {
			rc <- ResponseChannel{Token: "", Refresh: "", Error: errors.New("token invalid")}
		}

		defer close(rc)

		// token string to slice
		sToken := strings.Split(rToken, ".")
		if len(sToken) != 3 {
			rc <- ResponseChannel{Token: "", Refresh: "", Error: errors.New("token invalid")}
		}
		// decode base64 from token
		var decodedByte, errDecode = base64.StdEncoding.DecodeString(sToken[1])
		logger.New(errDecode, logger.SetType(logger.ERROR))

		var decodedString = string(decodedByte)
		var claims = jwt.MapClaims{}
		if err := json.Unmarshal([]byte(decodedString), &claims); err != nil {
			logger.New(err, logger.SetType(logger.ERROR))
		}
		// claim data from refresh token
		var tx, err = jwt.ParseWithClaims(rToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Key().Refresh), nil
		})
		if err != nil {
			logger.New(err)
			rc <- ResponseChannel{Token: "", Refresh: "", Error: err}
		}

		dataClaim := tx.Claims.(jwt.MapClaims)
		newClaimToken := &jwtTokenClaims{
			Email: dataClaim["email"].(string),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
			},
		}

		// gen new token
		nToken, err := t.generateToken(newClaimToken)
		if err != nil {
			logger.New(err)
			rc <- ResponseChannel{Token: "", Refresh: "", Error: err}
		}
		rc <- ResponseChannel{Token: nToken, Refresh: "", Error: nil}
	}()

	for {
		select {
		case <-ctx.Done():
			return "", errors.New("request time out")
		case res := <-rc:
			return res.Token, res.Error
		}
	}
}

func (t *tokenService) generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var tk, err = token.SignedString([]byte(config.Key().Token))
	if err != nil {
		logger.New(err)
		return "", err
	}

	return tk, nil
}

func (t *tokenService) generateRefreshToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var tk, err = token.SignedString([]byte(config.Key().Refresh))
	if err != nil {
		logger.New(err)
		return "", err
	}

	return tk, nil
}
