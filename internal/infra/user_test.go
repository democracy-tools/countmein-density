package job

import (
	"testing"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserPreference(t *testing.T) {

	t.Skip("infra")
	dsc := ds.NewClientWrapper(env.Project)
	updateUserPreference(t, dsc, "123", "A14")
}

func updateUserPreference(t *testing.T, dsc ds.Client, phone string, preference string) {

	var users []ds.User
	require.NoError(t, dsc.GetFilter(ds.KindUser, "phone", "=", phone, users), phone)
	require.Len(t, users, 1)
	require.NoError(t, dsc.Put(ds.KindUser, users[0].Id, &ds.User{
		Id:         users[0].Id,
		Name:       users[0].Name,
		Phone:      phone,
		Preference: preference,
		Time:       users[0].Time,
	}), phone)
}
