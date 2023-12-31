package api

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vadym-98/simple_bank/db/mock"
	db "github.com/vadym-98/simple_bank/db/sqlc"
	"github.com/vadym-98/simple_bank/token"
	"github.com/vadym-98/simple_bank/util"
	"github.com/vadym-98/simple_bank/util/faker"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateTransferAPI(t *testing.T) {
	u1 := faker.NewUser().Get()
	u2 := faker.NewUser().Get()

	a1 := faker.NewAccount().WithOwner(u1.Username).WithCurrency(util.USD).Get()
	a2 := faker.NewAccount().WithOwner(u2.Username).WithCurrency(util.USD).Get()
	transfer := faker.NewTransfer().WithFromAccountID(a1.ID).WithToAccountID(a2.ID).Get()
	tr := db.TransferTxResult{
		Transfer:    transfer,
		FromAccount: a1,
		ToAccount:   a2,
		FromEntry:   faker.NewEntry().WithAccountID(a1.ID).WithAmount(-transfer.Amount).Get(),
		ToEntry:     faker.NewEntry().WithAccountID(a2.ID).WithAmount(transfer.Amount).Get(),
	}
	stdTransReq := transferRequest{
		FromAccountID: a1.ID,
		ToAccountID:   a2.ID,
		Amount:        transfer.Amount,
		Currency:      util.USD,
	}

	testCases := []struct {
		name          string
		body          transferRequest
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: stdTransReq,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, u1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: a1.ID,
						ToAccountID:   a2.ID,
						Amount:        transfer.Amount,
					})).
					Times(1).
					Return(tr, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a1.ID)).
					Times(1).
					Return(a1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a2.ID)).
					Times(1).
					Return(a2, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchStruct[db.TransferTxResult](t, recorder.Body, tr)
			},
		},
		{
			name: "InternalError",
			body: stdTransReq,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, u1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: a1.ID,
						ToAccountID:   a2.ID,
						Amount:        transfer.Amount,
					})).
					Times(1).
					Return(db.TransferTxResult{}, sql.ErrTxDone)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a1.ID)).
					Times(1).
					Return(a1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a2.ID)).
					Times(1).
					Return(a2, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "FromAccountIDNotFound",
			body: stdTransReq,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, u1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: a1.ID,
						ToAccountID:   a2.ID,
						Amount:        transfer.Amount,
					})).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a2.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "GetFromAccountIDInternalError",
			body: stdTransReq,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, u1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: a1.ID,
						ToAccountID:   a2.ID,
						Amount:        transfer.Amount,
					})).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a2.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "FromAccountIDInvalidCurrency",
			body: transferRequest{
				FromAccountID: a1.ID,
				ToAccountID:   a2.ID,
				Amount:        transfer.Amount,
				Currency:      util.EUR,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, u1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: a1.ID,
						ToAccountID:   a2.ID,
						Amount:        transfer.Amount,
					})).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a1.ID)).
					Times(1).
					Return(a1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a2.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountIDInvalidCurrency",
			body: stdTransReq,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, u1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				toAccInvalidCurrency := a2
				toAccInvalidCurrency.Currency = util.EUR

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: a1.ID,
						ToAccountID:   a2.ID,
						Amount:        transfer.Amount,
					})).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a1.ID)).
					Times(1).
					Return(a1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccInvalidCurrency.ID)).
					Times(1).
					Return(toAccInvalidCurrency, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: transferRequest{
				FromAccountID: a1.ID,
				ToAccountID:   a2.ID,
				Amount:        transfer.Amount,
				Currency:      "invalid",
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, u1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: a1.ID,
						ToAccountID:   a2.ID,
						Amount:        transfer.Amount,
					})).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a1.ID)).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a2.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: stdTransReq,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, u2.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(a1.ID)).
					Times(1).
					Return(a1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPost, "/transfers", createBody(t, tc.body))
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
