package entity

type User struct {
	ID        int     `json:"-"`
	Login     string  `json:"-"`
	Balance   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
