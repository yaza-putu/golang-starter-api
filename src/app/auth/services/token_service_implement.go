package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
	"strings"
	"time"
	"yaza/src/app/auth/entities"
	"yaza/src/config"
)

type (
	tokenService   struct{}
	jwtTokenClaims struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		jwt.StandardClaims
	}
	jwtRefreshClaim struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		OldToken string `json:"old_token"`
		jwt.StandardClaims
	}
)

// NewToken is constructor
func NewToken() TokenInterface {
	return &tokenService{}
}

// Create token
// return arg1 string token
// return arg2 string refresh token
// return arg3 error error
func (t *tokenService) Create(user entities.User) (string, string, error) {
	// create token
	token, err := t.generateToken(&jwtTokenClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	})
	if err != nil {
		return "", "", err
	}

	// gen refresh token
	refresh, err := t.generateRefreshToken(&jwtRefreshClaim{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		OldToken: token,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})
	if err != nil {
		return "", "", err
	}

	return token, refresh, nil
}

func (t *tokenService) Refresh(rToken string) (string, error) {
	// token string to slice
	sToken := strings.Split(rToken, ".")
	if len(sToken) != 3 {
		return "", errors.New("token Invalid")
	}

	// decode base64 from token
	var decodedByte, _ = base64.StdEncoding.DecodeString(sToken[1])
	var decodedString = string(decodedByte)
	var claims = jwt.MapClaims{}
	if err := json.Unmarshal([]byte(decodedString), &claims); err != nil {
		log.Error(err)
	}
	// claim data from refresh token
	var tx, err = jwt.ParseWithClaims(rToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Key().Refresh), nil
	})
	if err != nil {
		return "", err
	}

	dataClaim := tx.Claims.(jwt.MapClaims)
	newClaimToken := &jwtTokenClaims{
		ID:    dataClaim["id"].(string),
		Name:  dataClaim["name"].(string),
		Email: dataClaim["email"].(string),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	}

	// gen new token
	nToken, err := t.generateToken(newClaimToken)
	if err != nil {
		return "", err
	}

	return nToken, nil
}

func (t *tokenService) generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var tk, err = token.SignedString([]byte(config.Key().Token))
	if err != nil {
		return "", err
	}

	return tk, nil
}

func (t *tokenService) generateRefreshToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var tk, err = token.SignedString([]byte(config.Key().Refresh))
	if err != nil {
		return "", err
	}

	return tk, nil
}
