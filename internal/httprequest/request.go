package httprequest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Alheor/gophermart/internal/models"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/go-playground/validator/v10"
)

func ParseRegistrationRequest(reqBody []byte) (*models.RegistrationForm, error) {

	var form models.RegistrationForm
	err := json.Unmarshal(reqBody, &form)
	if err != nil {
		return nil, err
	}

	err = getValidator().Struct(form)
	if err != nil {
		return nil, parseError(err)
	}

	return &form, nil
}

func ParseLoginRequest(reqBody []byte) (*models.LoginForm, error) {

	var form models.LoginForm
	err := json.Unmarshal(reqBody, &form)
	if err != nil {
		return nil, err
	}

	err = getValidator().Struct(form)
	if err != nil {
		return nil, parseError(err)
	}

	return &form, nil
}

func ParseUserOrderRequest(reqBody []byte) (*models.UserOrderForm, error) {

	var form models.UserOrderForm
	form.OrderID = string(reqBody)

	err := getValidator().Struct(form)
	if err != nil {
		return nil, parseError(err)
	}

	err = goluhn.Validate(form.OrderID)
	if err != nil {
		return nil, Error{Code: http.StatusUnprocessableEntity}
	}

	return &form, nil
}

func ParseUserWithdrawOrderRequest(reqBody []byte) (*models.UserWithdrawOrder, error) {

	var form models.UserWithdrawOrder
	err := json.Unmarshal(reqBody, &form)
	if err != nil {
		return nil, err
	}

	err = getValidator().Struct(form)
	if err != nil {
		return nil, parseError(err)
	}

	err = goluhn.Validate(form.Order)
	if err != nil {
		return nil, Error{Code: http.StatusUnprocessableEntity}
	}

	return &form, nil
}

func getValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}

func parseError(err error) Error {
	var invalidValidationError *validator.ValidationErrors
	if errors.As(err, &invalidValidationError) {
		return Error{Code: http.StatusInternalServerError}
	}

	for _, err := range err.(validator.ValidationErrors) {
		return Error{Code: http.StatusBadRequest, Field: strings.ToLower(err.Field()), Mess: err.Tag()}
	}

	return Error{Code: http.StatusInternalServerError}
}
