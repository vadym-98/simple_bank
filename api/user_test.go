package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vadym-98/simple_bank/db/mock"
	db "github.com/vadym-98/simple_bank/db/sqlc"
	mocktoken "github.com/vadym-98/simple_bank/token/mock"
	"github.com/vadym-98/simple_bank/util"
	"github.com/vadym-98/simple_bank/util/faker"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, pwd string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg: arg, password: pwd}
}

func TestCreateUserAPI(t *testing.T) {
	const pwd = "mysecret"
	u := faker.NewUser().Get()
	uParams := db.CreateUserParams{
		Username: u.Username,
		FullName: u.FullName,
		Email:    u.Email,
	}
	stdUreq := createUserRequest{
		Username: u.Username,
		Password: pwd,
		FullName: u.FullName,
		Email:    u.Email,
	}
	uResp := userResponse{
		Username: u.Username,
		FullName: u.FullName,
		Email:    u.Email,
	}

	testCases := []struct {
		name          string
		body          createUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: stdUreq,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(uParams, pwd)).
					Times(1).
					Return(u, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchStruct[userResponse](t, recorder.Body, uResp)
			},
		},
		{
			name: "InternalError",
			body: stdUreq,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TooLongPassword",
			body: createUserRequest{
				Username: u.Username,
				Password: util.RandomString(73),
				FullName: u.FullName,
				Email:    u.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: stdUreq,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: createUserRequest{
				Username: "invalid#123",
				Password: pwd,
				FullName: u.FullName,
				Email:    u.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: createUserRequest{
				Username: u.Username,
				Password: pwd,
				FullName: u.FullName,
				Email:    "asdas",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: createUserRequest{
				Username: u.Username,
				Password: "123",
				FullName: u.FullName,
				Email:    u.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			req, err := http.NewRequest(http.MethodPost, "/users", createBody(t, tc.body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	const pwd = "mysecret"
	u := faker.NewUser().Get()
	hashedPwd, err := util.HashPassword(pwd)
	require.NoError(t, err)
	u.HashedPassword = hashedPwd

	lReq := loginUserRequest{
		Username: u.Username,
		Password: pwd,
	}

	testCases := []struct {
		name          string
		body          loginUserRequest
		setUpServer   func(t *testing.T, store *mockdb.MockStore, maker *mocktoken.MockMaker) *Server
		buildStubs    func(store *mockdb.MockStore, maker *mocktoken.MockMaker, duration time.Duration)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: lReq,
			setUpServer: func(t *testing.T, store *mockdb.MockStore, maker *mocktoken.MockMaker) *Server {
				return newTestServer(t, store)
			},
			buildStubs: func(store *mockdb.MockStore, maker *mocktoken.MockMaker, duration time.Duration) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(u.Username)).
					Times(1).
					Return(u, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: loginUserRequest{
				Username: "$invalid",
				Password: pwd,
			},
			setUpServer: func(t *testing.T, store *mockdb.MockStore, maker *mocktoken.MockMaker) *Server {
				return newTestServer(t, store)
			},
			buildStubs: func(store *mockdb.MockStore, maker *mocktoken.MockMaker, duration time.Duration) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPwd",
			body: loginUserRequest{
				Username: u.Username,
				Password: "123",
			},
			setUpServer: func(t *testing.T, store *mockdb.MockStore, maker *mocktoken.MockMaker) *Server {
				return newTestServer(t, store)
			},
			buildStubs: func(store *mockdb.MockStore, maker *mocktoken.MockMaker, duration time.Duration) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: lReq,
			setUpServer: func(t *testing.T, store *mockdb.MockStore, maker *mocktoken.MockMaker) *Server {
				return newTestServer(t, store)
			},
			buildStubs: func(store *mockdb.MockStore, maker *mocktoken.MockMaker, duration time.Duration) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(u.Username)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NoUserFound",
			body: lReq,
			setUpServer: func(t *testing.T, store *mockdb.MockStore, maker *mocktoken.MockMaker) *Server {
				return newTestServer(t, store)
			},
			buildStubs: func(store *mockdb.MockStore, maker *mocktoken.MockMaker, duration time.Duration) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(u.Username)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidPwd",
			body: loginUserRequest{
				Username: u.Username,
				Password: "invalidPwd",
			},
			setUpServer: func(t *testing.T, store *mockdb.MockStore, maker *mocktoken.MockMaker) *Server {
				return newTestServer(t, store)
			},
			buildStubs: func(store *mockdb.MockStore, maker *mocktoken.MockMaker, duration time.Duration) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(u.Username)).
					Times(1).
					Return(u, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "FailedCreateToken",
			body: lReq,
			setUpServer: func(t *testing.T, store *mockdb.MockStore, maker *mocktoken.MockMaker) *Server {
				s := newTestServer(t, store)
				s.tokenMaker = maker
				return s
			},
			buildStubs: func(store *mockdb.MockStore, maker *mocktoken.MockMaker, duration time.Duration) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(u.Username)).
					Times(1).
					Return(u, nil)

				maker.EXPECT().
					CreateToken(gomock.Eq(u.Username), gomock.Eq(duration)).
					Times(1).
					Return("", errors.New("failed to encode payload to []byte"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			maker := mocktoken.NewMockMaker(ctrl)

			server := tc.setUpServer(t, store, maker)

			tc.buildStubs(store, maker, server.config.AccessTokenDuration)
			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPost, "/users/login", createBody(t, tc.body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
