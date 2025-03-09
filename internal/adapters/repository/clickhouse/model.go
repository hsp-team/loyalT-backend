package clickhouse

import (
	"time"

	"github.com/google/uuid"
)

type Reason string

const (
	ReasonQRScan    Reason = "qr_scan"
	ReasonCouponBuy Reason = "coupon_buy"
)

func (r Reason) String() string {
	return string(r)
}

type CoinBalanceChange struct {
	BusinessID    uuid.UUID
	UserID        uuid.UUID
	ProgramID     uuid.UUID
	BalanceChange int64
	Reason        Reason
	CouponID      *uuid.UUID
	Timestamp     time.Time
}

type UserActivityStats struct {
	QRScansCount  uint64
	CouponsBought uint64
}

type DailyTotalUniqueUsers struct {
	Date             time.Time
	TotalUniqueUsers uint64
}

type DailyActiveUsers struct {
	Date        time.Time
	ActiveUsers uint64
}

type BusinessStats struct {
	TotalUniqueUsers    uint64
	NewUsersInPeriod    uint64
	ActiveUsersInPeriod uint64
}

// LoyaltyProgramStats represents statistics about points spent in loyalty program
type LoyaltyProgramStats struct {
	TotalPointsReceived    int64
	PointsReceivedInPeriod int64
	TotalCouponsPurchased  uint64
	CouponsInPeriod        uint64
}
