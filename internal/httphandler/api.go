package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/models"
)

const (
	//HeaderAcceptEncoding header "Accept-Encoding" name
	HeaderAcceptEncoding = `Accept-Encoding`

	//HeaderContentEncodingGzip header Content-Encoding value gzip
	HeaderContentEncodingGzip = `gzip`

	//HeaderContentType header "Content-Type" name
	HeaderContentType = `Content-Type`

	//HeaderContentTypeJSON header Content-Type value application/json
	HeaderContentTypeJSON = `application/json`

	//HeaderContentTypeXGzip header Content-Type value application/x-gzip
	HeaderContentTypeXGzip = `application/x-gzip`

	//HeaderContentEncoding header "Content-Encoding" name
	HeaderContentEncoding = `Content-Encoding`

	HeaderContentTypeJSONValue = `application/json`

	HeaderContentTypeName = `Content-Type`
)

func sendAPIResponse(respWr http.ResponseWriter, resp *models.APIResponse) {
	rawByte, err := json.Marshal(resp)
	if err != nil {
		logger.Error(`response marshal error`, err)
		respWr.WriteHeader(http.StatusInternalServerError)
		return
	}

	respWr.Header().Add(HeaderContentType, HeaderContentTypeJSON)
	respWr.WriteHeader(resp.StatusCode)

	_, err = respWr.Write(rawByte)
	if err != nil {
		logger.Error(`write response error`, err)
		respWr.WriteHeader(http.StatusInternalServerError)
		return
	}
}
