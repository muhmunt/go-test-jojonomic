package service

import (
	"input-harga-storage-service/model"
	"input-harga-storage-service/repository"
)

type Service interface {
	StorePrice(request model.Price) (model.Price, error)
}

type service struct {
	repository repository.PriceRepository
}

func NewService(repository repository.PriceRepository) *service {
	return &service{repository}
}

func (s *service) StorePrice(request model.Price) (model.Price, error) {
	price := model.Price{}
	price.AdminID = request.AdminID
	price.HargaTopup = request.HargaTopup
	price.HargaBuyback = request.HargaBuyback

	newPrice, err := s.repository.Save(price)

	if err != nil {
		return newPrice, err
	}

	return newPrice, nil
}
