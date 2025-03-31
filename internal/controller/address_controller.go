package controller

import (
	"task4.2.3/internal/models"
	"task4.2.3/internal/service"
)

type AddressController struct {
	service service.AddressService
}

func NewAddressController(service service.AddressService) *AddressController {
	return &AddressController{service: service}
}

func (a *AddressController) Search(query string) ([]models.Address, error) {
	return a.service.Search(query)
}

func (a *AddressController) Geocode(lat, lng string) ([]models.Address, error) {
	return a.service.Geocode(lat, lng)
}
