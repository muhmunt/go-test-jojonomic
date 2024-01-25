package request

type CreateTopupRequest struct {
	Norek string `json:"norek" binding:"required"`
	Harga int    `json:"harga" binding:"required"`
	Gram  string `json:"gram" binding:"required"`
}
