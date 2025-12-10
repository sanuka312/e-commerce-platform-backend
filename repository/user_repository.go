package repository

import (
	"shophub-backend/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetOrCreateUser(keycloakUserID string) (*model.User, error)
}

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserRepository(Db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{Db: Db}
}

func (r *UserRepositoryImpl) GetOrCreateUser(keycloakUserID string) (*model.User, error) {
	user := &model.User{
		KeycloakUserID: keycloakUserID,
	}
	// Use FirstOrCreate to ensure user exists without overwriting existing data
	err := r.Db.FirstOrCreate(user, model.User{KeycloakUserID: keycloakUserID}).Error
	return user, err
}
