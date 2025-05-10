package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/Alheor/gophermart/internal/httphandler"
	"github.com/Alheor/gophermart/internal/logger"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHTTPHandler обработка сжатых запросов
func GzipHTTPHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get(httphandler.HeaderAcceptEncoding), httphandler.HeaderContentEncodingGzip) {
			f(resp, req)
			return
		}

		ct := req.Header.Get(httphandler.HeaderContentType)
		if ct != httphandler.HeaderContentTypeJSON && ct != httphandler.HeaderContentTypeXGzip {
			f(resp, req)
			return
		}

		if ct == httphandler.HeaderContentTypeXGzip {
			var data []byte

			data, err := io.ReadAll(req.Body)
			if err != nil {
				f(resp, req)
				return
			}

			data, err = GzipDecompress(data)
			if err != nil {
				logger.Error(`gzip decompress error:`, err)
				f(resp, req)
				return
			}

			req.Body = io.NopCloser(bytes.NewReader(data))
		}

		gz, err := gzip.NewWriterLevel(resp, gzip.BestSpeed)
		if err != nil {
			f(resp, req)
			logger.Error(`gzip error:`, err)
			return
		}

		defer gz.Close()

		resp.Header().Set(httphandler.HeaderContentEncoding, httphandler.HeaderContentEncodingGzip)

		f(gzipWriter{resp, gz}, req)
	}
}

// GzipDecompress разжатие данных
func GzipDecompress(data []byte) ([]byte, error) {

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var b bytes.Buffer

	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
