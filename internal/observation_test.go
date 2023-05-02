package internal_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/democracy-tools/countmein-density/internal"
	"github.com/democracy-tools/countmein-density/internal/ds"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/stretchr/testify/require"
)

func TestHandle_CreateObservation(t *testing.T) {

	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(
		&ds.Observation{
			Time:    time.Now().Unix(),
			User:    "israel",
			Polygon: "A18B",
			Density: 1.5,
		}))

	r, err := http.NewRequest(http.MethodPost, "/observations", bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	internal.NewHandle(ds.NewInMemoryClient(), whatsapp.NewInMemoryClient()).CreateObservation(w, r)

	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
}

func TestHandle_GetObservations(t *testing.T) {

	r, err := http.NewRequest(http.MethodGet, "/observations", nil)
	require.NoError(t, err)
	r.Header.Add("Accept", "application/json")
	w := httptest.NewRecorder()

	dsc := ds.NewInMemoryClient().(*ds.InMemoryClient)
	now := time.Now().Unix()
	dsc.SetGetFilterDelegate(
		func(kind ds.Kind, filters []ds.FilterField, dst interface{}) error {
			reflect.ValueOf(dst).Elem().Set(reflect.ValueOf([]ds.Observation{
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
			}))

			return nil
		})

	internal.NewHandle(dsc, whatsapp.NewInMemoryClient()).GetObservations(w, r)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	var res map[string][]ds.Observation
	require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&res))
	require.Len(t, res["observations"], 1)
	require.Equal(t, float32(2), res["observations"][0].Density)

	// buf := new(strings.Builder)
	// io.Copy(buf, w.Result().Body)
	// a := buf.String()
	// fmt.Println(a)
}
