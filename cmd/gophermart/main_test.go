package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Alheor/gophermart/internal/accural"
	"github.com/Alheor/gophermart/internal/config"
	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/models"
	"github.com/Alheor/gophermart/internal/repository"
	"github.com/Alheor/gophermart/internal/router"
	"github.com/Alheor/gophermart/internal/userauth"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type want struct {
	code         int
	responseBody string
	headerName   string
	headerValue  string
	cookieName   string
}

type test struct {
	name           string
	requestURL     string
	requestBody    []byte
	cookie         *http.Cookie
	requestHeaders map[string]string
	method         string
	want           want
}

var user1 = &models.User{ID: 344}
var user2 = &models.User{ID: 345}

func TestCreateUser(t *testing.T) {

	cfg := config.Load()
	err := logger.Init()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = repository.Init(ctx, cfg.DatabaseUri)
	assert.NoError(t, err)

	prepareDB(t)

	tests := []test{
		{
			name:        `positive test: register auth POST`,
			requestURL:  `/api/user/register`,
			requestBody: []byte(`{"login": "test", "password":"password"}`),
			method:      http.MethodPost,
			want: want{
				code:       http.StatusOK,
				cookieName: models.CookiesName,
			},
		}, {
			name:        `negative test: register auth POST`,
			requestURL:  `/api/user/register`,
			requestBody: []byte(`{"login": "test", "password":"password1"}`),
			method:      http.MethodPost,
			want: want{
				code: http.StatusConflict,
			},
		},
	}

	runTests(t, tests)
}

func TestLogin(t *testing.T) {
	cfg := config.Load()
	err := logger.Init()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = repository.Init(ctx, cfg.DatabaseUri)
	assert.NoError(t, err)

	prepareDB(t)

	tests := []test{
		{
			name:        `positive test: register auth POST`,
			requestURL:  `/api/user/register`,
			requestBody: []byte(`{"login": "test", "password":"password"}`),
			method:      http.MethodPost,
			want: want{
				code: http.StatusOK,
			},
		}, {
			name:        `positive test: login auth POST`,
			requestURL:  `/api/user/login`,
			requestBody: []byte(`{"login": "test", "password":"password"}`),
			method:      http.MethodPost,
			want: want{
				code:       http.StatusOK,
				cookieName: models.CookiesName,
			},
		}, {
			name:        `negative test: login auth POST`,
			requestURL:  `/api/user/login`,
			requestBody: []byte(`{"login": "test", "password":"password1"}`),
			method:      http.MethodPost,
			want: want{
				code: http.StatusUnauthorized,
			},
		},
	}

	runTests(t, tests)
}

func TestAddOrder(t *testing.T) {
	cfg := config.Load()
	err := logger.Init()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = repository.Init(ctx, cfg.DatabaseUri)
	assert.NoError(t, err)

	prepareDB(t)
	addUsersToDB(t)

	accural.InitConnector(cfg.AccrualAddr)
	accural.InitService(ctx)

	tests := []test{
		{
			name:        `positive test: add order POST`,
			requestURL:  `/api/user/orders`,
			requestBody: []byte(`123455`),
			method:      http.MethodPost,
			cookie:      prepareCookie(user1),
			want: want{
				code: http.StatusAccepted,
			},
		}, {
			name:        `positive test: add order twice POST`,
			requestURL:  `/api/user/orders`,
			requestBody: []byte(`123455`),
			method:      http.MethodPost,
			cookie:      prepareCookie(user1),
			want: want{
				code: http.StatusOK,
			},
		}, {
			name:        `negative test: add invalid order POST`,
			requestURL:  `/api/user/orders`,
			requestBody: []byte(`test`),
			method:      http.MethodPost,
			cookie:      prepareCookie(user1),
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		}, {
			name:        `negative test: add empty order POST`,
			requestURL:  `/api/user/orders`,
			requestBody: []byte(``),
			method:      http.MethodPost,
			cookie:      prepareCookie(user1),
			want: want{
				code: http.StatusBadRequest,
			},
		}, {
			name:        `negative test: add order unauthorized user POST`,
			requestURL:  `/api/user/orders`,
			requestBody: []byte(`123455`),
			method:      http.MethodPost,
			want: want{
				code: http.StatusUnauthorized,
			},
		}, {
			name:        `negative test: add order other user POST`,
			requestURL:  `/api/user/orders`,
			requestBody: []byte(`123455`),
			method:      http.MethodPost,
			cookie:      prepareCookie(user2),
			want: want{
				code: http.StatusConflict,
			},
		},
	}

	runTests(t, tests)
}

func TestGetOrder(t *testing.T) {
	cfg := config.Load()
	err := logger.Init()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = repository.Init(ctx, cfg.DatabaseUri)
	assert.NoError(t, err)

	prepareDB(t)
	addUsersToDB(t)

	accural.InitConnector(cfg.AccrualAddr)
	accural.InitService(ctx)

	tests := []test{
		{
			name:        `positive test: add order POST`,
			requestURL:  `/api/user/orders`,
			requestBody: []byte(`123455`),
			method:      http.MethodPost,
			cookie:      prepareCookie(user1),
			want: want{
				code: http.StatusAccepted,
			},
		}, {
			name:        `positive test: get order GET`,
			requestURL:  `/api/user/orders`,
			requestBody: []byte(`123455`),
			method:      http.MethodGet,
			cookie:      prepareCookie(user1),
			want: want{
				code: http.StatusOK,
			},
		},
	}

	runTests(t, tests)
}

func TestGetUserBalance(t *testing.T) {
	cfg := config.Load()
	err := logger.Init()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = repository.Init(ctx, cfg.DatabaseUri)
	assert.NoError(t, err)

	prepareDB(t)
	addUsersToDB(t)

	tests := []test{
		{
			name:       `positive test: get user balance GET`,
			requestURL: `/api/user/balance`,
			method:     http.MethodGet,
			cookie:     prepareCookie(user1),
			want: want{
				code:         http.StatusOK,
				responseBody: `{"current":13.4,"withdrawn":11.2}`,
			},
		},
	}

	runTests(t, tests)
}

func TestAddWithdrawOrder(t *testing.T) {
	cfg := config.Load()
	err := logger.Init()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = repository.Init(ctx, cfg.DatabaseUri)
	assert.NoError(t, err)

	prepareDB(t)
	addUsersToDB(t)

	tests := []test{
		{
			name:        `positive test: add withdraw order POST`,
			requestURL:  `/api/user/balance/withdraw`,
			requestBody: []byte(`{"order": "123455","sum": 7}`),
			method:      http.MethodPost,
			cookie:      prepareCookie(user1),
			want: want{
				code: http.StatusOK,
			},
		}, {
			name:        `negative test: add withdraw order with not enough memory POST`,
			requestURL:  `/api/user/balance/withdraw`,
			requestBody: []byte(`{"order": "123455","sum": 1751}`),
			method:      http.MethodPost,
			cookie:      prepareCookie(user1),
			want: want{
				code: http.StatusPaymentRequired,
			},
		}, {
			name:        `negative test: add withdraw order with invalid order POST`,
			requestURL:  `/api/user/balance/withdraw`,
			requestBody: []byte(`{"order": "333","sum": 1}`),
			method:      http.MethodPost,
			cookie:      prepareCookie(user1),
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		}, {
			name:       `positive test: get user balance GET`,
			requestURL: `/api/user/balance`,
			method:     http.MethodGet,
			cookie:     prepareCookie(user1),
			want: want{
				code:         http.StatusOK,
				responseBody: `{"current":6.4,"withdrawn":18.2}`,
			},
		},
	}

	runTests(t, tests)
}

func runTests(t *testing.T, tests []test) {

	ts := httptest.NewServer(router.GetRoutes())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			req, err := http.NewRequest(test.method, ts.URL+test.requestURL, bytes.NewReader(test.requestBody))
			require.NoError(t, err)

			if test.cookie != nil {
				req.AddCookie(test.cookie)
			}

			for hName, hVal := range test.requestHeaders {
				req.Header.Set(hName, hVal)
			}

			client := ts.Client()
			transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
			transport.DisableCompression = true
			client.Transport = transport

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.want.code, resp.StatusCode)

			if test.want.cookieName != `` {
				cookieExists := false

				for _, value := range resp.Cookies() {
					if value.Name == test.want.cookieName {
						cookieExists = true
						break
					}
				}

				assert.True(t, cookieExists)
			}

			resBody, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err)

			if test.want.responseBody != `` {
				assert.Equal(t, test.want.responseBody, string(resBody))
			}
		})
	}
}

func prepareDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := repository.GetConnection().Conn.Exec(ctx, `DELETE FROM "order"; DELETE FROM "withdrawal"; DELETE FROM "user";`)
	require.NoError(t, err)
}

func addUsersToDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := repository.GetConnection().Conn.Exec(ctx,
		`INSERT INTO "user" (id, login, pass, balance, withdrawn) VALUES (@id1, @login1, @pass1, @balance1, @withdrawn1),(@id2, @login2, @pass2, @balance2, @withdrawn2)`,
		pgx.NamedArgs{"id1": user1.ID, "login1": `test1`, "pass1": `test1`, "balance1": 13.4, "withdrawn1": 11.2, "id2": user2.ID, "login2": `test2`, "pass2": `test2`, "balance2": 130, "withdrawn2": 50},
	)

	require.NoError(t, err)
}

func prepareCookie(user *models.User) *http.Cookie {
	cookieValue := string(userauth.GetSignature(strconv.Itoa(user.ID))) + strconv.Itoa(user.ID)

	return &http.Cookie{
		Name:  models.CookiesName,
		Value: base64.StdEncoding.EncodeToString([]byte(cookieValue)),
	}
}
