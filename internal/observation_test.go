package internal_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/democracy-tools/countmein-density/internal"
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/stretchr/testify/require"
)

func TestHandle_CreateObservation(t *testing.T) {

	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(
		&internal.Observation{
			Time:    time.Now().Unix(),
			User:    "israel",
			Polygon: "A18B",
			Density: 1.5,
		}))

	r, err := http.NewRequest(http.MethodPost, "/observations", bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	internal.NewHandle(ds.NewInMemoryClient()).CreateObservation(w, r)

	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
}

func TestHandle_GetObservations(t *testing.T) {

	r, err := http.NewRequest(http.MethodGet, "/observations", nil)
	require.NoError(t, err)
	r.Header.Add("Accept", "application/json")
	w := httptest.NewRecorder()

	client := ds.NewInMemoryClient().(*ds.InMemoryClient)
	now := time.Now().Unix()
	client.SetGetByTimeDst([]internal.Observation{
		{
			Time:    now - 10,
			User:    "israel",
			Polygon: "A8",
			Density: 1.5,
		},
		{
			Time:    now,
			User:    "israel",
			Polygon: "A8",
			Density: 2,
		},
		{
			Time:    now - 5,
			User:    "israel",
			Polygon: "A8",
			Density: 5,
		},
	})

	internal.NewHandle(client).GetObservations(w, r)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	var res map[string][]internal.Observation
	require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&res))
	require.Len(t, res["observations"], 1)
	require.Equal(t, float32(2), res["observations"][0].Density)
}
