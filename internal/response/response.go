package response

import (
	"encoding/json"
	"github.com/Alheor/gophermart/internal/logger"
	"net/http"
)

const (
	HeaderContentTypeJSONValue = `application/json`
	HeaderContentTypeName      = `Content-Type`
)

type Error struct {
	Code  int    `json:"-"`
	Field string `json:"field"`
	Mess  string `json:"mess"`
}

func SendErrorResponse(w http.ResponseWriter, error *Error) {
	w.Header().Set(HeaderContentTypeName, HeaderContentTypeJSONValue)

	if error.Code == http.StatusInternalServerError {
		http.Error(w, error.Mess, http.StatusInternalServerError)
		return
	}

	rawByte, err := json.Marshal(error)
	if err != nil {
		logger.GetLogger().Panic(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(error.Code)

	_, err = w.Write(rawByte)
	if err != nil {
		logger.GetLogger().Panic(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
