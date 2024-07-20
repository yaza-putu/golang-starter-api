package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
)

// Token / **************************************************************
type Token interface {
	Create(user entity.User) (string, string, error)
	Refresh(rToken string) (string, error)
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
)

// NewToken is constructor
func NewToken() Token {
	return &tokenService{}
}

// Create token
// return arg1 string token
// return arg2 string refresh token
// return arg3 error error
func (t *tokenService) Create(user entity.User) (token string, rToken string, e error) {
	token, err := t.generateToken(&jwtTokenClaims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	})
	if err != nil {
		logger.New(err)
		return "", "", err
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
		return "", "", err
	}

	return token, refresh, nil
}

func (t *tokenService) Refresh(rToken string) (token string, e error) {
	if !strings.Contains(rToken, ".") {
		return "", errors.New("invalid token")
	}

	// token string to slice
	sToken := strings.Split(rToken, ".")
	if len(sToken) != 3 {
		return "", errors.New("invalid token")
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
		return "", err
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
		return "", err
	}
	return nToken, nil
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
