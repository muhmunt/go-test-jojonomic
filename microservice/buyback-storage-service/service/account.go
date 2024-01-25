package service

import (
	"buyback-storage-service/helper"
	"buyback-storage-service/model"
	"buyback-storage-service/repository"
	"errors"
)

type AccountService interface {
	UpdateOrInsertAccount(request model.Transaction) (model.Account, error)
	FindById(noRek string) (model.Account, error)
}

type accountService struct {
	repository repository.AccountRepository
}

func NewAccount(repository repository.AccountRepository) *accountService {
	return &accountService{repository}
}

func (s *accountService) FindById(norek string) (model.Account, error) {
	getAccount, err := s.repository.FindById(norek)

	if err != nil {
		return getAccount, errors.New("account not found.")
	}

	return getAccount, nil
}

func (s *accountService) UpdateOrInsertAccount(request model.Transaction) (model.Account, error) {
	account := model.Account{}

	getAccount, err := s.FindById(request.Norek)

	if err == nil {
		getAccount.Saldo = helper.SubDecimal(getAccount.Saldo, request.Gram)
		newAccount, err := s.repository.Update(getAccount.Norek, getAccount)

		if err != nil {
			return newAccount, err
		}

		return newAccount, nil
	}

	requestGram, _ := helper.DecimalFromString(request.Gram)
	account.Norek = request.Norek
	account.Saldo = requestGram

	newAccount, err := s.repository.Create(account)

	if err != nil {
		return newAccount, err
	}

	return newAccount, nil
}
