package request

type CreatePriceRequest struct {
	AdminID      string `json:"admin_id" binding:"required"`
	HargaTopup   int    `json:"harga_topup" binding:"required"`
	HargaBuyback int    `json:"harga_buyback" binding:"required"`
}
