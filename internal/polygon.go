package internal

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/democracy-tools/countmein-density/internal/action"
	"github.com/democracy-tools/countmein-density/internal/ds"
)

func (h *Handle) GetAvailablePolygons(w http.ResponseWriter, r *http.Request) {

	demonstration, err := ds.GetKaplanDemonstration(h.dsc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := action.GetAvailablePolygons(h.dsc, demonstration.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string][]string{"polygons": toSortedPolygonSlice(res)})
}

func toSortedPolygonSlice(polygons map[string]string) []string {

	var res []string
	for curr := range polygons {
		res = append(res, curr)
	}
	sort.Strings(res)

	return res
}
