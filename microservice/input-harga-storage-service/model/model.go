package model

type Price struct {
	AdminID      string
	HargaTopup   int
	HargaBuyback int
}

func (Price) TableName() string {
	return "harga"
}