package dto

type QRGetResponse struct {
	Data string `json:"data"`
}

type UserQRScanRequest struct {
	Code string `param:"code" validate:"required,len=8"`
}

type UserQRScanResponse struct {
	Username string `json:"username"`
	Balance  uint   `json:"balance"`
}

type UserEnrollCoinRequest struct {
	Code string `json:"code" validate:"required,len=8"`
}
