package service

import (
	"errors"
	"cek-mutasi-service/model"
	"cek-mutasi-service/repository"
)

type PriceService interface {
	Find() (model.Price, error)
	FindById(adminId string) (model.Price, error)
	StorePrice(request model.Price) (model.Price, error)
}

type priceService struct {
	repository repository.PriceRepository
}

func NewPrice(repository repository.PriceRepository) *priceService {
	return &priceService{repository}
}

func (s *priceService) FindById(adminId string) (model.Price, error) {
	getPrice, err := s.repository.FindById(adminId)

	if err == nil {
		return getPrice, errors.New("admin already input price.")
	}

	return getPrice, nil
}

func (s *priceService) Find() (model.Price, error) {
	getPrice, err := s.repository.Find()

	if err != nil {
		return getPrice, errors.New("admin does not input the price.")
	}

	return getPrice, nil
}

func (s *priceService) StorePrice(request model.Price) (model.Price, error) {
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
