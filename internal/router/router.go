package router

import (
	"net/http"

	"github.com/Alheor/gophermart/internal/compress"
	"github.com/Alheor/gophermart/internal/httphandler"
	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/userauth"

	"github.com/go-chi/chi/v5"
)

type HTTPMiddleware func(f http.HandlerFunc) http.HandlerFunc

// GetRoutes Загрузка маршрутизации
func GetRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post(`/api/user/register`,
		middlewareConveyor(httphandler.UserRegistration, logger.LoggingHTTPHandler, compress.GzipHTTPHandler))

	r.Post(`/api/user/login`,
		middlewareConveyor(httphandler.UserLogin, logger.LoggingHTTPHandler, compress.GzipHTTPHandler))

	r.Post(`/api/user/orders`,
		middlewareConveyor(httphandler.AddUserOrders, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Get(`/api/user/orders`,
		middlewareConveyor(httphandler.GetUserOrders, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Get(`/api/user/balance`,
		middlewareConveyor(httphandler.GetUserBalance, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Post(`/api/user/balance/withdraw`,
		middlewareConveyor(httphandler.AddWithdrawOrder, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Get(`/api/user/withdrawals`,
		middlewareConveyor(httphandler.GetUserWithdraw, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	return r
}

func middlewareConveyor(h http.HandlerFunc, middlewares ...HTTPMiddleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}
