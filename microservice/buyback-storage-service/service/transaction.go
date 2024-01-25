package service

import (
	"buyback-storage-service/helper"
	"buyback-storage-service/model"
	"buyback-storage-service/repository"
)

type TransactionService interface {
	StoreTransaction(request model.Transaction) (model.Transaction, error)
}

type transactionService struct {
	repository repository.TransactionRepository
}

func NewTransaction(repository repository.TransactionRepository) *transactionService {
	return &transactionService{repository}
}

func (s *transactionService) StoreTransaction(request model.Transaction) (model.Transaction, error) {
	transaction := model.Transaction{}
	transaction.ID = helper.GenShortId()
	transaction.Date = request.Date
	transaction.Gram = request.Gram
	transaction.Type = "BUYBACK"
	transaction.Norek = request.Norek
	transaction.HargaTopup = request.HargaTopup
	transaction.HargaBuyback = request.HargaBuyback
	transaction.SaldoTerakhir = request.SaldoTerakhir

	newTransaction, err := s.repository.Save(transaction)

	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}
