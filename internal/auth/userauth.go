package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"

	"github.com/Alheor/gophermart/internal/config"
	"github.com/Alheor/gophermart/internal/entity"
)

type contextKey string

const CookiesName = `authKey`
const ContextValueName contextKey = `xAuthUser`

func GetSignature(id string) []byte {
	h := hmac.New(sha256.New, []byte(config.Options.SignatureKey))
	h.Write([]byte(CookiesName))
	h.Write([]byte(id))

	return h.Sum(nil)
}

func ParseCookie(cookie *http.Cookie) (int, error) {

	cookieValue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return 0, errors.New(`invalid cookie`)
	}

	value := string(cookieValue[sha256.Size:])
	if value == `` {
		return 0, errors.New(`invalid cookie`)
	}

	signature := cookieValue[:sha256.Size]
	expectedSignature := GetSignature(value)
	if !hmac.Equal(signature, expectedSignature) {
		return 0, errors.New(`invalid cookie`)
	}

	userID, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New(`invalid cookie`)
	}

	return userID, nil
}

func PrepareCookie(userID int) string {
	stringUserID := strconv.Itoa(userID)
	cookieValue := string(GetSignature(stringUserID)) + stringUserID

	return base64.StdEncoding.EncodeToString([]byte(cookieValue))
}

func GetUserFromContext(ctx context.Context) *entity.User {

	authUser := ctx.Value(ContextValueName)
	if authUser == nil {
		return nil
	}

	return authUser.(*entity.User)
}
