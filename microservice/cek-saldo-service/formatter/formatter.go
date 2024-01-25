package formatter

import "cek-saldo-service/model"

type AccountFormatter struct {
	Norek string  `json:"norek"`
	Saldo float64 `json:"saldo"`
}

func FormatAccount(account model.Account) AccountFormatter {
	formatter := AccountFormatter{
		Norek: account.Norek,
		Saldo: account.Saldo,
	}

	return formatter
}
