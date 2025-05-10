package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Alheor/gophermart/internal/httprequest"
	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/models"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/Alheor/gophermart/internal/userauth"
)

func AddWithdrawOrder(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "AddWithdrawOrder" handler`)

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		sendAPIResponse(resp, &models.APIResponse{Error: `invalid body`, StatusCode: http.StatusInternalServerError})
		return
	}

	var form *models.UserWithdrawOrder

	form, err = httprequest.ParseUserWithdrawOrderRequest(reqBody)
	if err != nil {
		var rErr httprequest.Error
		if errors.As(err, &rErr) {
			sendAPIResponse(resp, &models.APIResponse{Error: rErr.Field + `: ` + rErr.Mess, StatusCode: rErr.Code})
			return
		}

		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusInternalServerError})
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	user := userauth.GetUser(ctx)
	if user == nil {
		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusUnauthorized})
		return
	}

	err = repository.GetWithdrawOrderRepository().AddWithdrawOrder(ctx, user, form)
	if err != nil {

		var errNotEnoughMemory *models.ErrNotEnoughMemory
		if errors.As(err, &errNotEnoughMemory) {
			resp.WriteHeader(http.StatusPaymentRequired)
			return
		}

		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusOK})
}

func GetUserWithdraw(resp http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	user := userauth.GetUser(ctx)
	if user == nil {
		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusUnauthorized})
		return
	}

	list, err := repository.GetWithdrawOrderRepository().GetWithdrawals(ctx, user)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
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
