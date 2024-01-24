package model

type Transaction struct {
	ID            string
	Norek         string
	Type          string
	Gram          int
	HargaTopup    int
	HargaBuyback  int
	SaldoTerakhir int
	Date          int
}

func (Transaction) TableName() string {
	return "transaksi"
}
