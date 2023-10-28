package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
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
