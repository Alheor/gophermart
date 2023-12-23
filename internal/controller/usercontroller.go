package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Alheor/gophermart/internal/auth"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/Alheor/gophermart/internal/request"
	"github.com/Alheor/gophermart/internal/response"
	"net/http"
	"time"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	form, err := request.ParseRegisterRequest(r)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	repErr := repository.GetUserRepository().CreateUser(ctx, form)
	if repErr != nil {

		var uErr *repository.UniqueErr
		if errors.As(repErr, &uErr) {
			response.SendErrorResponse(w,
				&response.Error{Field: `login`, Mess: repository.ErrUserAlreadyExist, Code: http.StatusConflict},
			)
			return
		}

		response.SendErrorResponse(w,
			&response.Error{Field: `login`, Mess: repErr.Error(), Code: http.StatusInternalServerError},
		)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

	form, reqErr := request.ParseLoginRequest(r)
	if reqErr != nil {
		response.SendErrorResponse(w, reqErr)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	user, queryErr := repository.GetUserRepository().FindUser(ctx, form)
	if queryErr != nil {
		http.Error(w, queryErr.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name: auth.CookiesName,

			Value: auth.PrepareCookie(user.Id),
		},
	)

	w.WriteHeader(http.StatusOK)
}

func GetUserBalance(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user := auth.GetUserFromContext(ctx)
	user, err := repository.GetUserRepository().GetUserById(ctx, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawByte, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(rawByte)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
