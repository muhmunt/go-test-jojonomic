package request

type CreateTopupRequest struct {
	Norek string `json:"norek" binding:"required"`
	Harga int    `json:"harga" binding:"required"`
	Gram  int    `json:"gram" binding:"required"`
}
