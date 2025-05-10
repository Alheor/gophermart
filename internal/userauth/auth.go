package userauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/Alheor/gophermart/internal/models"
)

var signatureKey []byte

func Init(key string) {
	signatureKey = []byte(key)
}

// GetUser get user from context
func GetUser(ctx context.Context) *models.User {
	authUser := ctx.Value(models.ContextValueName)
	if authUser != nil {
		return authUser.(*models.User)
	}

	return nil
}

func AuthHTTPHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		cookie, _ := req.Cookie(models.CookiesName)
		if cookie == nil {
			f(resp, req)
			return
		}

		userCookie := parseCookie(cookie)
		if userCookie == nil {
			f(resp, req)
			return
		}

		ctxWithUser := context.WithValue(req.Context(), models.ContextValueName, &userCookie.User)
		f(resp, req.Clone(ctxWithUser))
	}
}

func parseCookie(cookie *http.Cookie) *models.UserCookie {

	if len(cookie.Value) < sha256.Size {
		return nil
	}

	cookieValue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil
	}

	signature := cookieValue[:sha256.Size]
	value := string(cookieValue)[sha256.Size:]

	if len(value) == 0 {
		return nil
	}

	expectedSignature := GetSignature(value)

	if !hmac.Equal(signature, expectedSignature) {
		return nil
	}

	id, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}

	return &models.UserCookie{User: models.User{ID: id}, Sign: signature}
}

func GetSignature(id string) []byte {
	h := hmac.New(sha256.New, signatureKey)
	h.Write([]byte(models.CookiesName))
	h.Write([]byte(id))

	return h.Sum(nil)
}

func AddCookieToNewUser(resp http.ResponseWriter, user *models.User) {
	id := strconv.Itoa(user.ID)

	cookieValue := string(GetSignature(id)) + id

	http.SetCookie(resp,
		&http.Cookie{
			Name:  models.CookiesName,
			Value: base64.StdEncoding.EncodeToString([]byte(cookieValue)),
		},
	)
}
