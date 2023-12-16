package repository

import (
	"black-jack/models"
	"errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetByName(name string) (*models.User, error) {
	var existedUser models.User
	err := u.db.First(&existedUser, "name = ?", name).Error
	return &existedUser, err
}

func (u *UserRepository) Create(user *models.User) error {
	_, err := u.GetByName(user.Name)
	if err == nil {
		return errors.New("暱稱已被使用, 請替換")
	} else if err == gorm.ErrRecordNotFound {
		err = u.db.Create(user).Error
		return err
	}
	return err
}
