package internal_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/democracy-tools/countmein-density/internal"
	"github.com/democracy-tools/countmein-density/internal/ds"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatapp"
	"github.com/stretchr/testify/require"
)

func TestHandle_Register(t *testing.T) {

	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(internal.Register{
		Name:  "Israel",
		Phone: "0514123456",
	}))
	r, err := http.NewRequest(http.MethodPost, "/register", &buf)
	require.NoError(t, err)

	dsc := ds.NewInMemoryClient().(*ds.InMemoryClient)
	dsc.SetGetFilterDelegate(func(kind ds.Kind, filterFieldName string, filterOperator string, filterValue interface{}, dst interface{}) error {
		return datastore.ErrNoSuchEntity
	})

	w := httptest.NewRecorder()

	internal.NewHandle(dsc, whatsapp.NewInMemoryClient()).Register(w, r)

	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
}
