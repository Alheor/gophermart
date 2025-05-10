package httprequest

type Error struct {
	error
	Code  int    `json:"-"`
	Field string `json:"field"`
	Mess  string `json:"mess"`
}
