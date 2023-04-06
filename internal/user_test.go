package internal_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/democracy-tools/countmein-density/internal"
	"github.com/democracy-tools/countmein-density/internal/ds"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatapp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHandle_CreateUser(t *testing.T) {

	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(&struct {
		Token string `json:"token"`
	}{Token: uuid.NewString()}))
	r, err := http.NewRequest(http.MethodPost, "/users", &buf)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	internal.NewHandle(ds.NewInMemoryClient(), whatsapp.NewInMemoryClient()).CreateUser(w, r)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
}
