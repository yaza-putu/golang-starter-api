package services

import (
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
	"github.com/magiconair/properties/assert"
	"strings"
	"testing"
	"time"
	"yaza/src/app/auth/entities"
	"yaza/src/config"
)

func TestCreateToken(t *testing.T) {
	token := NewToken()
	u := entities.User{
		ID:    "xyz",
		Name:  "anyone",
		Email: "anyone@gmail.com",
	}

	nToken, rToken, _ := token.Create(u)

	assert.Equal(t, len(strings.Split(nToken, ".")), 3)
	assert.Equal(t, len(strings.Split(rToken, ".")), 3)
}

func TestValidToken(t *testing.T) {
	token := NewToken()
	u := entities.User{
		ID:    "xyz",
		Name:  "anyone",
		Email: "anyone@gmail.com",
	}

	nToken, _, _ := token.Create(u)

	sToken := strings.Split(nToken, ".")
	assert.Equal(t, len(sToken), 3)

	// verify token
	var decodedByte, _ = base64.StdEncoding.DecodeString(sToken[1])
	var decodedString = string(decodedByte)
	var claims = jwt.MapClaims{}
	if err := json.Unmarshal([]byte(decodedString), &claims); err != nil {
		log.Error(err)
	}
	// claim data from refresh token
	var _, err = jwt.ParseWithClaims(nToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Key().Refresh), nil
	})

	assert.Equal(t, err, nil)
}

func TestExpiredToken(t *testing.T) {
	token := NewToken()
	user := entities.User{
		ID:    "xyz",
		Name:  "anyone",
		Email: "anyone@gmail.com",
	}

	newToken, err := token.generateToken(&jwtTokenClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Minute * 5).Unix(),
		},
	})

	assert.Equal(t, err, nil)

	sToken := strings.Split(newToken, ".")
	var decodedByte, _ = base64.StdEncoding.DecodeString(sToken[1])
	var decodedString = string(decodedByte)
	var claims = jwt.MapClaims{}
	if err = json.Unmarshal([]byte(decodedString), &claims); err != nil {
		log.Error(err)
	}
	// check valid token
	var _, err2 = jwt.ParseWithClaims(newToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Key().Refresh), nil
	})

	assert.Equal(t, err2.Error(), "Token is expired")
}

func TestRefreshToken(t *testing.T) {
	token := NewToken()
	user := entities.User{
		ID:    "xyz",
		Name:  "anyone",
		Email: "anyone@gmail.com",
	}

	oldToken, err := token.generateToken(&jwtTokenClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	})
	assert.Equal(t, err, nil)

	rToken, err := token.generateRefreshToken(&jwtRefreshClaim{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		OldToken: oldToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	newToken, err := token.Refresh(rToken)
	assert.Equal(t, err, nil)

	sToken := strings.Split(newToken, ".")
	var decodedByte, _ = base64.StdEncoding.DecodeString(sToken[1])
	var decodedString = string(decodedByte)
	var claims = jwt.MapClaims{}
	if err = json.Unmarshal([]byte(decodedString), &claims); err != nil {
		log.Error(err)
	}
	// check valid token
	var _, err2 = jwt.ParseWithClaims(newToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Key().Refresh), nil
	})

	// write logic test here
	assert.Equal(t, err2, nil)
}

func TestFailedClaimRefreshToken(t *testing.T) {
	token := NewToken()
	u := entities.User{
		ID:    "xyz",
		Name:  "anyone",
		Email: "anyone@gmail.com",
	}

	// wrong refresh token
	oldToken, _, _ := token.Create(u)

	newToken, _ := token.Refresh(oldToken)

	// i want failed to create new token if not use refresh token to claim new token
	assert.Equal(t, newToken, "")
}
