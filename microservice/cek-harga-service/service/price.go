package service

import (
	"cek-harga-service/model"
	"cek-harga-service/repository"
	"errors"
)

type Service interface {
	Find() (model.Price, error)
	FindById(adminId string) (model.Price, error)
	StorePrice(request model.Price) (model.Price, error)
}

type service struct {
	repository repository.PriceRepository
}

func NewService(repository repository.PriceRepository) *service {
	return &service{repository}
}

func (s *service) FindById(adminId string) (model.Price, error) {
	getPrice, err := s.repository.FindById(adminId)

	if err == nil {
		return getPrice, errors.New("admin already input price.")
	}

	return getPrice, nil
}

func (s *service) Find() (model.Price, error) {
	getPrice, err := s.repository.Find()

	if err != nil {
		return getPrice, errors.New("admin does not input the price.")
	}

	return getPrice, nil
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
