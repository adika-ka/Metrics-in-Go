package service

import (
	"task4.2.3/internal/models"
	"task4.2.3/internal/repository"
)

type AddressService interface {
	Search(query string) ([]models.Address, error)
	Geocode(lat, lng string) ([]models.Address, error)
}

type addressService struct {
	repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) Search(query string) ([]models.Address, error) {
	return s.repo.Search(query)
}

func (s *addressService) Geocode(lat, lng string) ([]models.Address, error) {
	return s.repo.Geocode(lat, lng)
}
