package dto

import "github.com/google/uuid"

type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required,plaintext"`
	Email    string `json:"email" validate:"required,email,plaintext"`
	Password string `json:"password" validate:"required,plaintext"`
}

type UserRegisterResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserGetResponse struct {
	ID         uuid.UUID              `json:"id"`
	Name       string                 `json:"name"`
	Email      string                 `json:"email"`
	Statistics UserStatisticsResponse `json:"statistics"`
}

type UserStatisticsResponse struct {
	UserQrScannedCount uint64 `json:"user_qr_scanned_count"`
	CouponsBought      uint64 `json:"coupons_bought"`
}
