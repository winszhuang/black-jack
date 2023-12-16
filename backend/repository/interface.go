package repository

import "black-jack/models"

type IUserRepository interface {
	Create(user *models.User) error
	GetByName(name string) (*models.User, error)
}
