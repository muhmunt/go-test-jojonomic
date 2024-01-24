package service

import (
	"topup-storage-service/helper"
	"topup-storage-service/model"
	"topup-storage-service/repository"
)

type AccountService interface {
	UpdateAccount(norek string, request model.Account) (model.Account, error)
}

type accountService struct {
	repository repository.AccountRepository
}

func NewAccount(repository repository.AccountRepository) *accountService {
	return &accountService{repository}
}

func (s *accountService) UpdateAccount(norek string, request model.Account) (model.Account, error) {
	account := model.Account{}
	account.Norek = helper.GenShortId()

	newAccount, err := s.repository.Update(account)

	if err != nil {
		return newAccount, err
	}

	return newAccount, nil
}
