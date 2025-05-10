package httphandler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/models"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/Alheor/gophermart/internal/userauth"
)

func GetUserBalance(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "GetUserBalance" handler`)

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	user := userauth.GetUser(ctx)
	if user == nil {
		sendAPIResponse(resp, &models.APIResponse{StatusCode: http.StatusUnauthorized})
		return
	}

	user, err := repository.GetUserRepository().GetUserByID(ctx, user.ID)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	rawByte, err := json.Marshal(user)
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
