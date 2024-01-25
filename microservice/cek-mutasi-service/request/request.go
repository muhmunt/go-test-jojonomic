package request

type GetTransactionRequest struct {
	Norek     string `json:"norek" binding:"required"`
	StartDate int    `json:"start_date" binding:"required"`
	EndDate   int    `json:"end_date" binding:"required"`
}
