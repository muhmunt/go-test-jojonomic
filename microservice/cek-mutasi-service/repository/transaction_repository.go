package repository

import (
	"cek-mutasi-service/model"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindById(noRek string, startDate, endDate int) ([]model.Transaction, error)
	Save(transaction model.Transaction) (model.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Save(transaction model.Transaction) (model.Transaction, error) {
	err := r.db.Create(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *transactionRepository) FindById(noRek string, startDate, endDate int) ([]model.Transaction, error) {
	var getTransaction []model.Transaction

	err := r.db.Where("norek = ? AND date BETWEEN ? AND ?", noRek, startDate, endDate).
		Find(&getTransaction).Error

	if err != nil {
		return getTransaction, err
	}

	return getTransaction, nil
}
