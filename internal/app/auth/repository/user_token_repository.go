package repository

import (
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-api/internal/database"
)

type (
	UserToken interface {
		Create(data entity.UserToken) (entity.UserToken, error)
		FindByDeviceId(deviceId string) (entity.UserToken, error)
		Revoke(deviceId string) error
		Update(deviceId string, data entity.UserToken) (entity.UserToken, error)
	}
	userTokenRepository struct {
		entity entity.UserToken
	}
)

func NewUserToken() *userTokenRepository {
	return &userTokenRepository{
		entity: entity.UserToken{},
	}
}

func (u *userTokenRepository) Create(data entity.UserToken) (entity.UserToken, error) {
	// find by userid
	e := u.entity
	database.Instance.Model(u.entity).Where("user_id", data.UserId).Where("ip", data.IP).Where("device", data.Device).First(&e)

	// update
	if e.ID != "" {
		db := database.Instance.Model(u.entity).Where("id", e.ID).Updates(&data)
		return data, db.Error
	}

	// create
	db := database.Instance.Create(&data)

	return data, db.Error
}

func (u *userTokenRepository) Revoke(deviceId string) error {
	db := database.Instance.Where("device_id", deviceId).Delete(u.entity)

	return db.Error
}

func (u *userTokenRepository) FindByDeviceId(deviceId string) (entity.UserToken, error) {
	e := u.entity
	db := database.Instance.Where("device_id", deviceId).First(&e)

	return e, db.Error
}

func (u *userTokenRepository) Update(id string, data entity.UserToken) (entity.UserToken, error) {
	db := database.Instance.Model(u.entity).Where("id", id).Updates(&data)

	return data, db.Error
}
