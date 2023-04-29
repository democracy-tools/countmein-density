package infra

import (
	"strings"
	"testing"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserPreference(t *testing.T) {

	// env.Initialize()
	t.Skip("infra")
	dsc := ds.NewClientWrapper(env.Project)
	updateUserPreference(t, dsc, "+972 12-345-6789", "W13")
}

func updateUserPreference(t *testing.T, dsc ds.Client, phone string, preference string) {

	phone = phoneConvention(t, phone)

	var users []ds.User
	require.NoError(t, dsc.GetFilter(ds.KindUser, []ds.FilterField{{Name: "phone", Operator: "=", Value: phone}}, &users), phone)
	require.Len(t, users, 1)

	require.NoError(t, dsc.Put(ds.KindUser, users[0].Id, &ds.User{
		Id:         users[0].Id,
		Name:       users[0].Name,
		Phone:      phone,
		Preference: preference,
		Time:       users[0].Time,
	}), phone)
}

func phoneConvention(t *testing.T, phone string) string {

	phone = strings.ReplaceAll(strings.ReplaceAll(strings.TrimPrefix(phone, "+"), " ", ""), "-", "")
	require.Len(t, phone, 12)

	return phone
}
