package models

import "time"

type UserWithdrawOrder struct {
	Order string  `validate:"required"`
	Sum   float32 `validate:"required"`
}

type ErrNotEnoughMemory struct {
	error
}

type WithdrawalOrder struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
