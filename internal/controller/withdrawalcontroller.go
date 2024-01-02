package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Alheor/gophermart/internal/auth"
	"github.com/Alheor/gophermart/internal/entity"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/Alheor/gophermart/internal/request"
	"github.com/Alheor/gophermart/internal/response"
)

func AddWithdrawOrder(w http.ResponseWriter, r *http.Request) {
	form, reqErr := request.ParseAddUserWithdrawOrderRequest(r)
	if reqErr != nil {
		response.SendErrorResponse(w, reqErr)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user := auth.GetUserFromContext(ctx)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := repository.GetWithdrawalOrderRepository().AddWithdrawalOrder(ctx, user, form)
	if err != nil {

		var errNotEnoughMemory *entity.ErrNotEnoughMemory
		if errors.As(err, &errNotEnoughMemory) {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user := auth.GetUserFromContext(ctx)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	list, err := repository.GetWithdrawalOrderRepository().GetWithdrawals(ctx, user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(list) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	rawByte, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(response.HeaderContentTypeName, response.HeaderContentTypeJSONValue)

	_, err = w.Write(rawByte)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
