package formatter

import "cek-mutasi-service/model"

type TransactionFormatter struct {
	Norek         string  `json:"norek"`
	Type          string  `json:"type"`
	Gram          string  `json:"gram"`
	HargaTopup    int     `json:"harga_topup"`
	HargaBuyback  int     `json:"harga_buyback"`
	SaldoTerakhir float64 `json:"saldo_terakhir"`
	Date          int     `json:"date"`
}

func FormatTransaction(transaction model.Transaction) TransactionFormatter {
	formatter := TransactionFormatter{
		Norek:         transaction.Norek,
		Type:          transaction.Type,
		Gram:          transaction.Gram,
		HargaTopup:    transaction.HargaTopup,
		HargaBuyback:  transaction.HargaBuyback,
		SaldoTerakhir: transaction.SaldoTerakhir,
		Date:          transaction.Date,
	}

	return formatter
}

func FormatTransactions(transactions []model.Transaction) []TransactionFormatter {
	transactionsFormatter := []TransactionFormatter{}

	for _, transaction := range transactions {
		productFormatter := FormatTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, productFormatter)
	}

	return transactionsFormatter
}
