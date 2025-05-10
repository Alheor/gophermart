package models

type contextKeyXAuthUser string

const CookiesName = `authKey`
const ContextValueName contextKeyXAuthUser = `xAuthUser`

type User struct {
	ID        int     `json:"-"`
	Login     string  `json:"-"`
	Balance   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type UserCookie struct {
	User User
	Sign []byte
}
