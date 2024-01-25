package service

import (
	"cek-saldo-service/helper"
	"cek-saldo-service/model"
	"cek-saldo-service/repository"
	"cek-saldo-service/request"
	"errors"
)

type TransactionService interface {
	StoreTransaction(request model.Transaction) (model.Transaction, error)
	FindTransactionByNorek(request request.GetTransactionRequest) ([]model.Transaction, error)
}

type transactionService struct {
	repository repository.TransactionRepository
}

func NewTransaction(repository repository.TransactionRepository) *transactionService {
	return &transactionService{repository}
}

func (s *transactionService) FindTransactionByNorek(request request.GetTransactionRequest) ([]model.Transaction, error) {
	var getTransaction []model.Transaction
	getTransaction, err := s.repository.FindById(request.Norek, request.StartDate, request.EndDate)

	if err != nil {
		return getTransaction, errors.New("transaction not found.")
	}

	return getTransaction, nil
}

func (s *transactionService) StoreTransaction(request model.Transaction) (model.Transaction, error) {
	transaction := model.Transaction{}
	transaction.ID = helper.GenShortId()
	transaction.Date = request.Date
	transaction.Gram = request.Gram
	transaction.Type = "TOPUP"
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
