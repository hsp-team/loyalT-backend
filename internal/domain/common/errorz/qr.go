package errorz

import "errors"

var (
	ErrUserByQrNotFound     = errors.New("user by qr not found")
	ErrUserScanLimitReached = errors.New("user scan day limit reached")
)
