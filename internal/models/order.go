package models

import "time"

type UniqueErrByOrder struct {
	error
}

type UniqueErrByUserAndOrder struct {
	error
}

type UserOrderForm struct {
	OrderID string `validate:"required,min=2"`
}

type Order struct {
	Number     string    `json:"number"`
	Accrual    float32   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
	Status     string    `json:"status"`
}
