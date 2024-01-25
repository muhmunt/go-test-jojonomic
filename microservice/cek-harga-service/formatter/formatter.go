package formatter

import "cek-harga-service/model"

type PriceFormatter struct {
	HargaTopup   int `json:"harga_topup"`
	HargaBuyback int `json:"harga_buyback"`
}

func FormatAccount(price model.Price) PriceFormatter {
	formatter := PriceFormatter{
		HargaTopup:   price.HargaTopup,
		HargaBuyback: price.HargaBuyback,
	}

	return formatter
}
