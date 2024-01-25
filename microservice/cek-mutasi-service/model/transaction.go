package model

type Transaction struct {
	ID            string  `json:"id"`
	Norek         string  `json:"norek"`
	Type          string  `json:"type"`
	Gram          string  `json:"gram"`
	HargaTopup    int     `json:"harga_topup"`
	HargaBuyback  int     `json:"harga_buyback"`
	SaldoTerakhir float64 `json:"saldo_terakhir"`
	Date          int     `json:"date"`
}

func (Transaction) TableName() string {
	return "transaksi"
}
