package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/go-playground/validator/v10"

	"github.com/Alheor/gophermart/internal/response"
)

type RegisterForm struct {
	Login    string `validate:"required,min=3,max=255"`
	Password string `validate:"required,min=6,max=255"`
}

type LoginForm struct {
	Login    string `validate:"required,max=255"`
	Password string `validate:"required,max=255"`
}

type UserOrderForm struct {
	OrderID string `validate:"required,min=2"`
}

type UserWithdrawOrderForm struct {
	Order string  `validate:"required"`
	Sum   float32 `validate:"required"`
}

func ParseRegisterRequest(r *http.Request) (*RegisterForm, *response.Error) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError}
	}

	form := new(RegisterForm)

	if string(reqBody) == `` {
		reqBody = []byte(`{}`)
	}

	err = json.Unmarshal(reqBody, form)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError}
	}

	err = getValidator().Struct(form)
	if err == nil {
		return form, nil
	}

	return nil, parseError(err)
}

func ParseLoginRequest(r *http.Request) (*LoginForm, *response.Error) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError}
	}

	form := new(LoginForm)

	if string(reqBody) == `` {
		reqBody = []byte(`{}`)
	}

	err = json.Unmarshal(reqBody, form)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError}
	}

	err = getValidator().Struct(form)
	if err == nil {
		return form, nil
	}

	return nil, parseError(err)
}

func ParseAddUserOrderRequest(r *http.Request) (*UserOrderForm, *response.Error) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError}
	}

	form := new(UserOrderForm)
	form.OrderID = string(reqBody)

	err = getValidator().Struct(form)
	if err != nil {
		return nil, parseError(err)
	}

	if form.OrderID == `` {
		return nil, &response.Error{Code: http.StatusBadRequest}
	}

	err = goluhn.Validate(form.OrderID)
	if err != nil {
		return nil, &response.Error{Code: http.StatusUnprocessableEntity}
	}

	return form, nil
}

func ParseAddUserWithdrawOrderRequest(r *http.Request) (*UserWithdrawOrderForm, *response.Error) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError}
	}

	if string(reqBody) == `` {
		reqBody = []byte(`{}`)
	}

	form := new(UserWithdrawOrderForm)

	err = json.Unmarshal(reqBody, form)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError}
	}

	err = getValidator().Struct(form)
	if err != nil {
		return nil, parseError(err)
	}

	if form.Order == `` {
		return nil, &response.Error{Code: http.StatusBadRequest, Field: `order`}
	}

	if form.Sum <= 0 {
		return nil, &response.Error{Code: http.StatusBadRequest, Field: `sum`}
	}

	err = goluhn.Validate(form.Order)
	if err != nil {
		return nil, &response.Error{Code: http.StatusUnprocessableEntity, Field: `order`}
	}

	return form, nil
}

func getValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}

func parseError(err error) *response.Error {
	var invalidValidationError *validator.ValidationErrors
	if errors.As(err, &invalidValidationError) {
		return &response.Error{Code: http.StatusInternalServerError}
	}

	for _, err := range err.(validator.ValidationErrors) {
		return &response.Error{Code: http.StatusBadRequest, Field: strings.ToLower(err.Field()), Mess: err.Tag()}
	}

	return &response.Error{Code: http.StatusInternalServerError}
}
