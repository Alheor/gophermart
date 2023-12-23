package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Alheor/gophermart/internal/accural"
	"github.com/Alheor/gophermart/internal/auth"
	"github.com/Alheor/gophermart/internal/entity"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/Alheor/gophermart/internal/request"
	"github.com/Alheor/gophermart/internal/response"
	"net/http"
	"time"
)

func AddUserOrder(w http.ResponseWriter, r *http.Request) {

	form, reqErr := request.ParseAddUserOrderRequest(r)
	if reqErr != nil {
		response.SendErrorResponse(w, reqErr)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user := auth.GetUserFromContext(ctx)

	err := repository.GetOrderRepository().AddOrder(ctx, user, form.OrderID)
	if err != nil {

		var errO *entity.UniqueErrByOrder
		if errors.As(err, &errO) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		var errUO *entity.UniqueErrByUserAndOrder
		if errors.As(err, &errUO) {
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accural.Sync(form.OrderID)

	w.Header().Set(response.HeaderContentTypeName, response.HeaderContentTypeJSONValue)
	w.WriteHeader(http.StatusAccepted)
}

func GetUserOrders(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	user := auth.GetUserFromContext(ctx)

	list, err := repository.GetOrderRepository().GetOrders(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
