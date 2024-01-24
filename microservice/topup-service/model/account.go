package model

type Account struct {
	Norek         string
	SaldoTerakhir int
}

func (Account) TableName() string {
	return "rekening"
}
