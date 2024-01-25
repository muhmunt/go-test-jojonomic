package repository

import (
	"cek-mutasi-service/model"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Update(noRek string, account model.Account) (model.Account, error)
	Create(account model.Account) (model.Account, error)
	FindById(norek string) (model.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccount(db *gorm.DB) *accountRepository {
	return &accountRepository{db}
}

func (r *accountRepository) Update(noRek string, account model.Account) (model.Account, error) {
	err := r.db.Where("norek = ?", noRek).Save(&account).Error

	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *accountRepository) Create(account model.Account) (model.Account, error) {
	err := r.db.Create(&account).Error

	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *accountRepository) FindById(norek string) (model.Account, error) {
	var getAccount model.Account

	err := r.db.Where("norek = ?", norek).
		First(&getAccount).Error

	if err != nil {
		return getAccount, err
	}

	return getAccount, nil
}
