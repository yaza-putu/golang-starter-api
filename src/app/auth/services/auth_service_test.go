package services

import (
	"context"
	"github.com/magiconair/properties/assert"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	a := NewAuthService()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	r := a.Login(ctx, "admin@mail.com", "Password1")

	token := r.GetData().(map[string]string)

	assert.Equal(t, r.Code, 200)
	assert.Equal(t, strings.Split(token["token"], "."), 3)
}
