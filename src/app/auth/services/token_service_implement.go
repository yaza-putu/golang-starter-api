package services

import (
	"context"
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
	ResponseChannel struct {
		Token   string
		Refresh string
		Error   error
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
func (t *tokenService) Create(ctx context.Context, user entities.User) (string, string, error) {
	rc := make(chan ResponseChannel)
	go func() {
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
			rc <- ResponseChannel{Token: "", Refresh: "", Error: err}
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
		// token string to slice
		sToken := strings.Split(rToken, ".")
		if len(sToken) != 3 {
			rc <- ResponseChannel{Token: "", Refresh: "", Error: errors.New("token invalid")}
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
			rc <- ResponseChannel{Token: "", Refresh: "", Error: err}
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
