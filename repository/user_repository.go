package repository

import (
	"e-commerce-platform-backend/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserById(userId uint) (*model.User, error)
	CreateUser(user *model.User) error
}

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{Db: db}
}

func (r *UserRepositoryImpl) GetUserById(userId uint) (*model.User, error) {
	var user model.User
	if err := r.Db.First(&user, userId).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) CreateUser(user *model.User) error {
	return r.Db.Create(user).Error
}
