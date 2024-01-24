package repository

import (
	"input-harga-storage-service/model"

	"gorm.io/gorm"
)

type PriceRepository interface {
	Find() (model.Price, error)
	Save(price model.Price) (model.Price, error)
}

type priceRepository struct {
	db *gorm.DB
}

func NewPrice(db *gorm.DB) *priceRepository {
	return &priceRepository{db}
}

func (r *priceRepository) Save(price model.Price) (model.Price, error) {
	err := r.db.Create(&price).Error

	if err != nil {
		return price, err
	}

	return price, nil
}

func (r *priceRepository) Find() (model.Price, error) {
	var getPrice model.Price

	err := r.db.Find(&getPrice).Error

	if err != nil {
		return getPrice, err
	}

	return getPrice, nil
}
