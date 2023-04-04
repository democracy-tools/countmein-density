package internal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/democracy-tools/countmein-density/internal"
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHandle_CreateUser(t *testing.T) {

	r, err := http.NewRequest(http.MethodPost, "/users", nil)
	require.NoError(t, err)
	q := r.URL.Query()
	q.Add("token", uuid.NewString())
	r.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()

	internal.NewHandle(ds.NewInMemoryClient()).CreateUser(w, r)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
}
