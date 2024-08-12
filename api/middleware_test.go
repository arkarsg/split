package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arkarsg/splitapp/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func addAuth(t *testing.T, req *http.Request, tokenMaker token.TokenMaker, authType string, username string, duration time.Duration) {
	token, err := tokenMaker.CreateToken(username, duration)
	assert.NoError(t, err)
	authHeader := fmt.Sprintf("%s %s", authType, token)
	req.Header.Add(authorizationHeaderKey, authHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testTable := []struct {
		name          string
		setUpAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		checkResponse func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuth(t, request, tokenMaker, "bearer", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
			},
		},
		{
			name: "No Auth",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			},
		},
		{
			name: "Unsupported Auth",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuth(t, request, tokenMaker, "unsupported", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			},
		},
		{
			name: "Invalid Auth Format",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuth(t, request, tokenMaker, "", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			},
		},
		{
			name: "Expired Token",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuth(t, request, tokenMaker, "bearer", "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(c *gin.Context) {
					c.String(http.StatusOK, "pong")
				},
			)
			r := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			assert.NoError(t, err)
			tc.setUpAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(r, req)
			tc.checkResponse(t, r)
		})
	}
}
