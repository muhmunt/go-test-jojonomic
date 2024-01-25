package service

import (
	"buyback-service/helper"
	"buyback-service/model"
	"buyback-service/repository"
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

	newTransaction, err := s.repository.Save(transaction)

	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}
