package httphandler

import (
	"context"
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

func UserRegistration(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "UserRegistration" handler`)

	reqBody, err := io.ReadAll(req.Body)
	if err != nil || len(reqBody) == 0 {
		sendAPIResponse(resp, &models.APIResponse{Error: `invalid body`, StatusCode: http.StatusInternalServerError})
		return
	}

	var form *models.RegistrationForm

	form, err = httprequest.ParseRegistrationRequest(reqBody)
	if err != nil {
		var rErr httprequest.Error
		if errors.As(err, &rErr) {
			sendAPIResponse(resp, &models.APIResponse{Error: rErr.Field + `: ` + rErr.Mess, StatusCode: rErr.Code})
			return
		}

		logger.Error(`parse registration request error: `, err)
		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusInternalServerError})
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	user, err := repository.GetUserRepository().CreateUser(ctx, form)
	if err != nil {

		var uErr *repository.UniqueErr
		if errors.As(err, &uErr) {
			sendAPIResponse(resp, &models.APIResponse{Error: `User already exist`, StatusCode: http.StatusConflict})
			return
		}

		logger.Error(`get user error: `, err)
		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusInternalServerError})
		return
	}

	userauth.AddCookieToNewUser(resp, user)

	resp.WriteHeader(http.StatusOK)
}

func UserLogin(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "UserLogin" handler`)

	reqBody, err := io.ReadAll(req.Body)
	if err != nil || len(reqBody) == 0 {
		sendAPIResponse(resp, &models.APIResponse{Error: `invalid body`, StatusCode: http.StatusInternalServerError})
		return
	}

	var form *models.LoginForm

	form, err = httprequest.ParseLoginRequest(reqBody)
	if err != nil {
		var rErr httprequest.Error
		if errors.As(err, &rErr) {
			sendAPIResponse(resp, &models.APIResponse{Error: rErr.Field + `: ` + rErr.Mess, StatusCode: rErr.Code})
			return
		}

		logger.Error(`parse registration request error: `, err)
		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusInternalServerError})
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	user, err := repository.GetUserRepository().FindUser(ctx, form)
	if err != nil {
		logger.Error(`search user error: `, err)
		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusInternalServerError})
		return
	}

	if user == nil {
		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusUnauthorized})
		return
	}

	userauth.AddCookieToNewUser(resp, user)

	resp.WriteHeader(http.StatusOK)
}
