package job

import (
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
)

func CreatePreferrers() error {

	dsc := ds.NewClientWrapper(env.Project)
	for _, currPreference := range getPreferences() {
		err := dsc.Put(ds.KindPreference, currPreference.UserId, &currPreference)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPreferences() []ds.Preference {

	return []ds.Preference{
		{
			UserId:  "",
			Polygon: "",
		},
		{
			UserId:  "",
			Polygon: "",
		},
		{
			UserId:  "",
			Polygon: "",
		},
		{
			UserId:  "",
			Polygon: "",
		},
		{
			UserId:  "",
			Polygon: "",
		},
	}
}
