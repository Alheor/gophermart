package router

import (
	"context"
	"github.com/Alheor/gophermart/internal/auth"
	"github.com/Alheor/gophermart/internal/controller"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type HTTPMiddleware func(f http.HandlerFunc) http.HandlerFunc

func Init() *chi.Mux {
	r := chi.NewRouter()
	prepareRoutes(r)

	return r
}

func prepareRoutes(r *chi.Mux) {
	r.Post("/api/user/register", middlewareConveyor(controller.RegisterUser))
	r.Post("/api/user/login", middlewareConveyor(controller.LoginUser))
	r.Post("/api/user/orders", middlewareConveyor(controller.AddUserOrder, WithUserAuth))
	r.Get("/api/user/orders", middlewareConveyor(controller.GetUserOrders, WithUserAuth))
	r.Get("/api/user/balance", middlewareConveyor(controller.GetUserBalance, WithUserAuth))
	r.Post("/api/user/balance/withdraw", middlewareConveyor(controller.AddWithdrawOrder, WithUserAuth))
	r.Get("/api/user/withdrawals", middlewareConveyor(controller.GetUserWithdrawals, WithUserAuth))
}

func middlewareConveyor(h http.HandlerFunc, middlewares ...HTTPMiddleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}

func WithUserAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var userCookie *http.Cookie

		for _, cookie := range r.Cookies() {
			if cookie.Name == auth.CookiesName {
				userCookie = cookie
				break
			}
		}

		if userCookie == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := auth.ParseCookie(userCookie)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
		defer cancel()

		user, err := repository.GetUserRepository().GetUserByID(ctx, userID)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctxWithUser := context.WithValue(r.Context(), auth.ContextValueName, user)

		f(w, r.Clone(ctxWithUser))
	}
}
