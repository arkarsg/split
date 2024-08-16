package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/arkarsg/splitapp/db/mock"
	db "github.com/arkarsg/splitapp/db/sqlc"
	"github.com/arkarsg/splitapp/token"
	u "github.com/arkarsg/splitapp/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetTransactionAPI(t *testing.T) {
	transaction := randomTransaction()

	testTable := []struct {
		name          string
		transactionID int64
		setUpAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			transactionID: transaction.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuth(t, request, tokenMaker, "bearer", "xxx", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetTransactionById(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(transaction, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				requireBodyMatchTransaction(t, r.Body, transaction)
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/transaction/%v", tc.transactionID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			tc.setUpAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomTransaction() db.Transaction {
	return db.Transaction{
		ID:       u.RandomInt(1, 1000),
		Amount:   u.RandomAmount(),
		Currency: db.CurrencySGD,
		Title:    u.RandomString(6),
		PayerID:  u.RandomInt(1, 10),
	}
}

func requireBodyMatchTransaction(t *testing.T, body *bytes.Buffer, transaction db.Transaction) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	var haveTransaction db.Transaction
	err = json.Unmarshal(data, &haveTransaction)
	assert.NoError(t, err)
	assert.Equal(t, transaction, haveTransaction)
}
