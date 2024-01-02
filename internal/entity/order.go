package entity

import "time"

type UniqueErrByOrder struct {
	error
}

type UniqueErrByUserAndOrder struct {
	error
}

type Order struct {
	Number     string    `json:"number"`
	Accrual    float32   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
	Status     string    `json:"status"`
}
