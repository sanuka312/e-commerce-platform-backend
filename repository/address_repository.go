package repository

import (
	"shophub-backend/model"

	"gorm.io/gorm"
)

type AddressRepository interface {
	GetAddressesByUser(keycloakUserID string) ([]model.Address, error)
	CreateAddress(address *model.Address) error
	//GetAddressById(addressId uint) (*model.Address, error)
}

type AddressRepositoryImpl struct {
	Db *gorm.DB
}

func NewAddressRepository(Db *gorm.DB) AddressRepository {
	return &AddressRepositoryImpl{Db: Db}
}

func (r *AddressRepositoryImpl) GetAddressesByUser(keycloakUserID string) ([]model.Address, error) {
	var addresses []model.Address
	err := r.Db.Where("keycloak_user_id=?", keycloakUserID).Find(&addresses).Error
	return addresses, err
}

func (r *AddressRepositoryImpl) CreateAddress(address *model.Address) error {
	return r.Db.Create(address).Error
}

func (r *AddressRepositoryImpl) GetAddressById(addressId uint) (*model.Address, error) {
	var address model.Address
	err := r.Db.First(&address, addressId).Error
	return &address, err
}
