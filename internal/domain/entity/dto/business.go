package dto

import (
	"github.com/google/uuid"
	"time"
)

type BusinessRegisterRequest struct {
	Name        string `json:"name" validate:"required,gte=3,lte=100,plaintext"`
	Email       string `json:"email" validate:"required,email,plaintext"`
	Password    string `json:"password" validate:"required,plaintext"`
	Description string `json:"description" validate:"omitempty,lte=150,plaintext"`
}

type BusinessRegisterResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Description string    `json:"description,omitempty"`
}

type BusinessLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type BusinessLoginResponse struct {
	Token string `json:"token"`
}

type BusinessReturn struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Description string    `json:"description,omitempty"`
}

type BusinessUpdateRequest struct {
	Name        string `json:"name" validate:"required,gte=3,lte=100,plaintext"`
	Description string `json:"description" validate:"omitempty,lte=150,plaintext"`
}

type BusinessStatsRequest struct {
	StartDate time.Time `query:"start_date" validate:"required"`
	EndDate   time.Time `query:"end_date" validate:"required,gtefield=StartDate"`
}

type BusinessStatsResponse struct {
	TotalUsers  uint64 `json:"total_users"`
	NewUsers    uint64 `json:"new_users"`
	ActiveUsers uint64 `json:"active_users"`
}

type BusinessCoinProgramStatsRequest struct {
	StartDate time.Time `query:"start_date" validate:"required"`
	EndDate   time.Time `query:"end_date" validate:"required,gtefield=StartDate"`
}

type BusinessCoinProgramStatsResponse struct {
	TotalPointsReceived    int64  `json:"total_points_received"`
	PointsReceivedInPeriod int64  `json:"points_received_in_period"`
	TotalCouponsPurchased  uint64 `json:"total_coupons_purchased"`
	CouponsInPeriod        uint64 `json:"coupons_purchased_in_period"`
}

type BusinessStatsDailyTotalUsersRequest struct {
	StartDate time.Time `query:"start_date" validate:"required"`
	EndDate   time.Time `query:"end_date" validate:"required,gtefield=StartDate"`
}

type BusinessStatsDailyTotalUsersResponse struct {
	Date       time.Time `json:"date"`
	TotalUsers uint64    `json:"total_users"`
}

type BusinessStatsDailyActiveUsersRequest struct {
	StartDate time.Time `query:"start_date" validate:"required"`
	EndDate   time.Time `query:"end_date" validate:"required,gtefield=StartDate"`
}

type BusinessStatsDailyActiveUsersResponse struct {
	Date        time.Time `json:"date"`
	ActiveUsers uint64    `json:"active_users"`
}
