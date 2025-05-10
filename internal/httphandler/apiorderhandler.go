package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Alheor/gophermart/internal/accural"
	"github.com/Alheor/gophermart/internal/httprequest"
	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/models"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/Alheor/gophermart/internal/userauth"
)

func AddUserOrders(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "AddUserOrders" handler`)

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		sendAPIResponse(resp, &models.APIResponse{Error: `invalid body`, StatusCode: http.StatusInternalServerError})
		return
	}

	var form *models.UserOrderForm

	form, err = httprequest.ParseUserOrderRequest(reqBody)
	if err != nil {
		var rErr httprequest.Error
		if errors.As(err, &rErr) {
			sendAPIResponse(resp, &models.APIResponse{Error: rErr.Field + `: ` + rErr.Mess, StatusCode: rErr.Code})
			return
		}

		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusInternalServerError})
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	user := userauth.GetUser(ctx)
	if user == nil {
		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusUnauthorized})
		return
	}

	err = repository.GetOrderRepository().AddOrder(ctx, user, form.OrderID)
	if err != nil {

		var errO *models.UniqueErrByOrder
		if errors.As(err, &errO) {
			sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusConflict})
			return
		}

		var errUO *models.UniqueErrByUserAndOrder
		if errors.As(err, &errUO) {
			sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusOK})
			return
		}

		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusInternalServerError})
		return
	}

	accural.Sync(form.OrderID)

	resp.Header().Set(HeaderContentTypeName, HeaderContentTypeJSONValue)
	resp.WriteHeader(http.StatusAccepted)
}

func GetUserOrders(resp http.ResponseWriter, req *http.Request) {

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	user := userauth.GetUser(ctx)
	if user == nil {
		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusUnauthorized})
		return
	}

	list, err := repository.GetOrderRepository().GetOrders(ctx, user)
	if err != nil {
		logger.Error(`Get orders error: `, err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(list) == 0 {
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	rawByte, err := json.Marshal(list)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Set(HeaderContentTypeName, HeaderContentTypeJSONValue)

	_, err = resp.Write(rawByte)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
}
