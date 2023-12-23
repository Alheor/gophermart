package entity

type AccrualOrder struct {
	Order   string  `json:"order"`
	Accrual float32 `json:"accrual,omitempty"`
	Status  string  `json:"status"`
}
