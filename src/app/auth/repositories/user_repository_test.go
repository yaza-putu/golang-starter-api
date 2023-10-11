package repositories

import (
	"context"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
	"yaza/src/app/auth/entities"
)

func TestFinByEmail(t *testing.T) {
	u := NewUserRepository()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	r, err := u.FindByEmail(ctx, "admin@mail.com")
	assert.Equal(t, err, nil)

	assert.Equal(t, r, entities.User{})
}
