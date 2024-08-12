package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/arkarsg/splitapp/db/mock"
	db "github.com/arkarsg/splitapp/db/sqlc"
	u "github.com/arkarsg/splitapp/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	account, pwd := createRandomAccount(t)
	testTable := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  account.Username,
				"password":  pwd,
				"full_name": account.FullName,
				"email":     account.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Username: account.Username,
					FullName: account.FullName,
					Email:    account.Email,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), EqAccount(arg, pwd)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"username":  account.Username,
				"password":  pwd,
				"full_name": account.FullName,
				"email":     account.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Username: account.Username,
					FullName: account.FullName,
					Email:    account.Email,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), EqAccount(arg, pwd)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid Username",
			body: gin.H{
				"username":  "$invalid",
				"password":  pwd,
				"full_name": account.FullName,
				"email":     account.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Email",
			body: gin.H{
				"username":  account.Username,
				"password":  pwd,
				"full_name": account.FullName,
				"email":     "account.email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Password too short",
			body: gin.H{
				"username":  account.Username,
				"password":  "short",
				"full_name": account.FullName,
				"email":     account.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := "/account"
			buf := jsonify(testCase.body)
			req, err := http.NewRequest(http.MethodPost, url, &buf)
			assert.NoError(t, err)
			server.router.ServeHTTP(recorder, req)
			testCase.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, user db.Account) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, gotAccount.Username)
	assert.Equal(t, user.FullName, gotAccount.FullName)
	assert.Equal(t, user.Email, gotAccount.Email)
	assert.Empty(t, gotAccount.HashedPassword)
}

// Custom gomock matcher
func EqAccount(args db.CreateAccountParams, pwd string) gomock.Matcher {
	return &AccountMatcher{args: args, plainPassword: pwd}
}

type AccountMatcher struct {
	args          db.CreateAccountParams
	plainPassword string
}

func (a *AccountMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateAccountParams)
	if !ok {
		return false
	}

	err := u.CheckPasswordHash(a.plainPassword, arg.HashedPassword)
	if err != nil {
		return false
	}

	a.args.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(a.args, arg)
}

func (a *AccountMatcher) String() string {
	return fmt.Sprintf("Matches %v with %v", a.args, a.plainPassword)
}

func createRandomAccount(t *testing.T) (account db.Account, pwd string) {
	pwd = u.RandomString(10)
	hashedPwd, err := u.HashPassword(pwd)
	assert.NoError(t, err)

	account = db.Account{
		Username:       u.RandomUser(),
		HashedPassword: hashedPwd,
		FullName:       u.RandomString(10),
		Email:          u.RandomEmail(),
	}
	return
}
