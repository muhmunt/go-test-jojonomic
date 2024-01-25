package model

type Account struct {
	Norek string
	Saldo float64
}

func (Account) TableName() string {
	return "rekening"
}
