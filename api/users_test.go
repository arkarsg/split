package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/arkarsg/splitapp/db/mock"
	db "github.com/arkarsg/splitapp/db/sqlc"
	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	user := randomUser()

	testTable := []struct {
		name          string
		userID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				requireBodyMatchUser(t, r.Body, user)
			},
		},
		{
			name:   "Not Found",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, r.Code)
			},
		},
		{
			name:   "Internal Server Error",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, r.Code)
			},
		},
		{
			name:   "Invalid ID",
			userID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(0)).
					Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/users/%v", testCase.userID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)
			server.router.ServeHTTP(recorder, req)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestCreateUser(t *testing.T) {
	createUserArgs := db.CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	}

	testTable := []struct {
		name          string
		user          db.CreateUserParams
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			user: createUserArgs,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(createUserArgs)).
					Times(1).
					Return(db.User{
						ID:       0,
						Username: createUserArgs.Username,
						Email:    createUserArgs.Email,
					}, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				requireBodyMatchCreateUserParmas(t, r.Body, createUserArgs)
			},
		},
		{
			name: "Internal Server Error",
			user: createUserArgs,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(createUserArgs)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, r.Code)
			},
		},
		{
			name: "Bad Request (Invalid email)",
			user: db.CreateUserParams{
				Username: u.RandomUser(),
				Email:    "invalid_email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(createUserArgs)).
					Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/users")
			buf := jsonify(testCase.user)
			req, err := http.NewRequest(http.MethodPost, url, &buf)
			assert.NoError(t, err)
			server.router.ServeHTTP(recorder, req)
			testCase.checkResponse(t, recorder)
		})
	}
}

func jsonify(obj any) bytes.Buffer {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(obj)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}

func randomUser() db.User {
	return db.User{
		ID:       u.RandomInt(1, 1000),
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	assert.NoError(t, err)

	assert.Equal(t, user, gotUser)
}

func requireBodyMatchCreateUserParmas(t *testing.T, body *bytes.Buffer, args db.CreateUserParams) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	assert.NoError(t, err)

	assert.Equal(t, args.Username, gotUser.Username)
	assert.Equal(t, args.Email, gotUser.Email)
}
