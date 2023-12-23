package entity

type User struct {
	Id        int     `json:"-"`
	Login     string  `json:"-"`
	Balance   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
