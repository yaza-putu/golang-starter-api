package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/repository"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
	"github.com/yaza-putu/golang-starter-api/pkg/unique"
)

// Token / **************************************************************
type Token interface {
	Create(user entity.User, ip string, device string) (string, string, error)
	Refresh(deviceId string) (string, string, error)
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
// return arg2 string deviceId
// return arg3 error error
func (t *tokenService) Create(user entity.User, ip string, device string) (token string, deviceId string, e error) {
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

	// gen device ID
	dId := unique.Key(52)

	// store refresh token
	_, err = repository.NewUserToken().Create(entity.UserToken{
		ID:        unique.Uid(),
		UserId:    user.ID,
		TokenType: entity.REFRESH_TOKEN,
		DeviceId:  dId,
		Token:     refresh,
		Device:    device,
		IP:        ip,
	})

	if err != nil {
		logger.New(err)
		return "", "", err
	}

	return token, dId, nil
}

func (t *tokenService) Refresh(deviceId string) (token string, devId string, e error) {
	// get refresh token by device ID
	userTokenRepository := repository.NewUserToken()

	userToken, errToken := userTokenRepository.FindByDeviceId(deviceId)
	if errToken != nil {
		return "", "", errToken
	}

	// token string to slice
	sToken := strings.Split(userToken.Token, ".")
	if len(sToken) != 3 {
		return "", "", errors.New("invalid token")
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
	var tx, err = jwt.ParseWithClaims(userToken.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Key().Refresh), nil
	})
	if err != nil {
		logger.New(err)
		return "", "", err
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
		return "", "", err
	}

	// generate new device ID
	dId := unique.Key(52)
	userToken.DeviceId = dId
	_, err = userTokenRepository.Update(userToken.ID, userToken)
	if err != nil {
		return "", "", err
	}

	return nToken, dId, nil
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
