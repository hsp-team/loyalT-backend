package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type StatisticRepository struct {
	conn Connection
}

func NewStatisticRepository(conn Connection) *StatisticRepository {
	return &StatisticRepository{
		conn: conn,
	}
}

func (r *StatisticRepository) InsertCoinBalanceChange(ctx context.Context, change CoinBalanceChange) error {
	query := `
		INSERT INTO coin_balance_changes (
			business_id, user_id, program_id,
			balance_change, reason, coupon_id,
			timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	var couponID interface{} = nil
	if change.CouponID != nil {
		couponID = *change.CouponID
	}

	return r.conn.Exec(ctx, query,
		change.BusinessID,
		change.UserID,
		change.ProgramID,
		change.BalanceChange,
		change.Reason.String(),
		couponID,
		change.Timestamp,
	)
}

// GetUserActivityStats returns the number of QR scans and coupons bought by a user
func (r *StatisticRepository) GetUserActivityStats(ctx context.Context, userID uuid.UUID) (UserActivityStats, error) {
	query := `
		SELECT 
			countIf(reason = ?) as qr_scans,
			countIf(reason = ?) as coupons_bought
		FROM coin_balance_changes 
		WHERE user_id = ?
	`

	var stats UserActivityStats
	row := r.conn.QueryRow(ctx, query,
		ReasonQRScan.String(),
		ReasonCouponBuy.String(),
		userID,
	)

	err := row.Scan(&stats.QRScansCount, &stats.CouponsBought)
	if err != nil {
		return UserActivityStats{}, fmt.Errorf("failed to scan user activity stats: %w", err)
	}

	return stats, nil
}

// GetDailyTotalUniqueUsers returns the cumulative count of unique users who have ever scanned a QR code, by day
func (r *StatisticRepository) GetDailyTotalUniqueUsers(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time) ([]DailyTotalUniqueUsers, error) {
	query := `
		SELECT 
			dates.date,
			count(DISTINCT user_id) as total_unique_users
		FROM (
			SELECT arrayJoin(arrayMap(x -> toDate(?) + toIntervalDay(x), range(dateDiff('day', ?, ?)))) as date
		) dates
		LEFT JOIN (
			SELECT 
				user_id,
				toDate(timestamp) as event_date
			FROM coin_balance_changes 
			WHERE business_id = ?
				AND reason = ?
				AND timestamp <= ?
		) events ON 1=1
		WHERE dates.date >= events.event_date OR events.event_date IS NULL
		GROUP BY dates.date
		ORDER BY dates.date
	`

	rows, err := r.conn.Query(ctx, query,
		startDate,
		startDate,
		endDate,
		businessID,
		ReasonQRScan.String(),
		endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily total unique users: %w", err)
	}
	defer rows.Close()

	var result []DailyTotalUniqueUsers
	for rows.Next() {
		var stat DailyTotalUniqueUsers
		if err := rows.Scan(&stat.Date, &stat.TotalUniqueUsers); err != nil {
			return nil, fmt.Errorf("failed to scan daily total unique users: %w", err)
		}
		result = append(result, stat)
	}

	return result, nil
}

// GetDailyActiveUsers returns the count of unique users who scanned a QR code each day
func (r *StatisticRepository) GetDailyActiveUsers(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time) ([]DailyActiveUsers, error) {
	query := `
		SELECT 
			dates.date,
			count(DISTINCT if(dates.date = events.event_date, events.user_id, NULL)) as active_users
		FROM (
			SELECT arrayJoin(arrayMap(x -> toDate(?) + toIntervalDay(x), range(dateDiff('day', ?, ?)))) as date
		) dates
		LEFT JOIN (
			SELECT 
				toDate(timestamp) as event_date,
				user_id
			FROM coin_balance_changes 
			WHERE business_id = ?
				AND reason = ?
				AND timestamp BETWEEN ? AND ?
		) events ON 1=1
		GROUP BY dates.date
		ORDER BY dates.date
	`

	rows, err := r.conn.Query(ctx, query,
		startDate,
		startDate,
		endDate,
		businessID,
		ReasonQRScan.String(),
		startDate,
		endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily active users: %w", err)
	}
	defer rows.Close()

	var result []DailyActiveUsers
	for rows.Next() {
		var stat DailyActiveUsers
		if err := rows.Scan(&stat.Date, &stat.ActiveUsers); err != nil {
			return nil, fmt.Errorf("failed to scan daily active users: %w", err)
		}
		result = append(result, stat)
	}

	return result, nil
}

// GetBusinessStats returns overall business statistics including total unique users and growth
func (r *StatisticRepository) GetBusinessStats(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time) (BusinessStats, error) {
	query := `
		SELECT
			(
				SELECT count(DISTINCT user_id)
				FROM coin_balance_changes
				WHERE business_id = ? AND reason = ? AND timestamp <= ?
			) as total_unique_users,
			(
				SELECT count(DISTINCT user_id)
				FROM coin_balance_changes
				WHERE business_id = ? AND reason = ?
					AND timestamp BETWEEN ? AND ?
					AND user_id NOT IN (
						SELECT DISTINCT user_id
						FROM coin_balance_changes
						WHERE business_id = ?
							AND reason = ?
							AND timestamp < ?
					)
			) as new_users_in_period,
			(
				SELECT count(DISTINCT user_id)
				FROM coin_balance_changes
				WHERE business_id = ? AND reason = ?
					AND timestamp BETWEEN ? AND ?
			) as active_users_in_period
	`

	var stats BusinessStats
	row := r.conn.QueryRow(ctx, query,
		businessID, ReasonQRScan.String(), endDate,
		businessID, ReasonQRScan.String(), startDate, endDate,
		businessID, ReasonQRScan.String(), startDate,
		businessID, ReasonQRScan.String(), startDate, endDate,
	)

	err := row.Scan(&stats.TotalUniqueUsers, &stats.NewUsersInPeriod, &stats.ActiveUsersInPeriod)
	if err != nil {
		return BusinessStats{}, fmt.Errorf("failed to scan business stats: %w", err)
	}

	return stats, nil
}

// GetLoyaltyProgramStats returns total points spent in loyalty program for a specified period
func (r *StatisticRepository) GetLoyaltyProgramStats(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time) (LoyaltyProgramStats, error) {
	query := `
		SELECT
			(
				SELECT sum(balance_change)
				FROM coin_balance_changes
				WHERE business_id = ? AND reason = ?
					AND timestamp <= ?
			) as total_points_spent,
			(
				SELECT sum(balance_change)
				FROM coin_balance_changes
				WHERE business_id = ? AND reason = ?
					AND timestamp >= ? AND timestamp < ?
			) as points_spent_in_period,
			(
				SELECT count(*)
				FROM coin_balance_changes
				WHERE business_id = ? AND reason = ?
					AND timestamp <= ?
			) as total_coupons_purchased,
			(
				SELECT count(*)
				FROM coin_balance_changes
				WHERE business_id = ? AND reason = ?
					AND timestamp >= ? AND timestamp < ?
			) as coupons_purchased_in_period
	`

	var stats LoyaltyProgramStats
	row := r.conn.QueryRow(ctx, query,
		businessID, ReasonQRScan.String(), endDate,
		businessID, ReasonQRScan.String(), startDate, endDate,
		businessID, ReasonCouponBuy.String(), endDate,
		businessID, ReasonCouponBuy.String(), startDate, endDate,
	)

	err := row.Scan(
		&stats.TotalPointsReceived,
		&stats.PointsReceivedInPeriod,
		&stats.TotalCouponsPurchased,
		&stats.CouponsInPeriod,
	)
	if err != nil {
		return LoyaltyProgramStats{}, fmt.Errorf("failed to scan loyalty program stats: %w", err)
	}

	return stats, nil
}

// GetUserBusinessQRScansCount returns the number of times a user has scanned QR codes for a specific business
func (r *StatisticRepository) GetUserBusinessQRScansCount(
	ctx context.Context,
	businessID, userID uuid.UUID,
) (uint64, error) {
	query := `
		SELECT count(*)
		FROM coin_balance_changes 
		WHERE business_id = ?
		AND user_id = ?
		AND reason = ?`

	var count uint64
	err := r.conn.QueryRow(ctx, query,
		businessID,
		userID,
		ReasonQRScan.String(),
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get QR scans count: %w", err)
	}

	return count, nil
}
