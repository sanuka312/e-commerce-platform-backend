package service

import (
	"shophub-backend/model"
	"shophub-backend/repository"
)

type AddressService interface {
	GetAddressesByUser(keycloakUserID string) ([]model.Address, error)
	CreateAddress(address *model.Address) error
}

type AddressServiceImpl struct {
	AddressRepository repository.AddressRepository
}

func NewAddressServiceImpl(AddressRepository repository.AddressRepository) (service AddressService, err error) {
	return &AddressServiceImpl{
		AddressRepository: AddressRepository,
	}, err
}

func (s *AddressServiceImpl) GetAddressesByUser(keycloakUserID string) ([]model.Address, error) {
	return s.AddressRepository.GetAddressesByUser(keycloakUserID)
}

func (s *AddressServiceImpl) CreateAddress(address *model.Address) error {
	return s.AddressRepository.CreateAddress(address)
}
