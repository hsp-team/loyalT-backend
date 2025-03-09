package errorz

import "errors"

var (
	ErrRewardNotFound = errors.New("reward not found")
	ErrNotEnoughCoins = errors.New("not enough coins")
)
