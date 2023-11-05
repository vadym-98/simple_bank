package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/vadym-98/simple_bank/db/sqlc"
	"github.com/vadym-98/simple_bank/util"
	"io"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	cfg := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(cfg, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode) // for different log output
	os.Exit(m.Run())
}

func createBody(t *testing.T, b interface{}) io.Reader {
	data, err := json.Marshal(b)
	require.NoError(t, err)

	return bytes.NewReader(data)
}

func requireBodyMatchStruct[E any](t *testing.T, body *bytes.Buffer, expectedStruct E) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var unmarshalledStruct E
	err = json.Unmarshal(data, &unmarshalledStruct)
	require.NoError(t, err)
	require.Equal(t, expectedStruct, unmarshalledStruct)
}
