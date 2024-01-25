package model

type Transaction struct {
	ID            string
	Norek         string
	Type          string
	Gram          string
	HargaTopup    int
	HargaBuyback  int
	SaldoTerakhir float64
	Date          int
}

func (Transaction) TableName() string {
	return "transaksi"
}
