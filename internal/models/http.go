package models

type RegistrationForm struct {
	Login    string `json:"login" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type LoginForm struct {
	Login    string `json:"login" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=255"`
}

type APIResponse struct {
	Result     string `json:"result,omitempty"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"-"`
}
