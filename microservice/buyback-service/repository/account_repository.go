package repository

import (
	"buyback-service/model"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Update(account model.Account) (model.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccount(db *gorm.DB) *accountRepository {
	return &accountRepository{db}
}

func (r *accountRepository) Update(account model.Account) (model.Account, error) {
	err := r.db.Save(&account).Error

	if err != nil {
		return account, err
	}

	return account, nil
}
