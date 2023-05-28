package action_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/democracy-tools/countmein-density/internal/action"
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/stretchr/testify/require"
)

func TestReport_Run(t *testing.T) {

	const name, phone = "Israel", "123456789"

	dsc := ds.NewInMemoryClient()
	dsc.(*ds.InMemoryClient).SetGetFilterDelegate(func(kind ds.Kind, filters []ds.FilterField, dst interface {
	}) error {
		if kind == ds.KindUser {
			reflect.ValueOf(dst).Elem().Set(reflect.ValueOf([]ds.User{{
				Id:         "1",
				Name:       name,
				Phone:      phone,
				Preference: "",
				Time:       time.Now().Unix(),
				Role:       ds.RoleAdmin,
			}}))
		}

		return nil
	})

	message, err := action.Create(dsc, whatsapp.NewInMemoryClient(),
		whatsapp.Contact{
			WaID:    phone,
			Profile: whatsapp.ContactProfile{Name: name}},
		whatsapp.Message{
			Type: whatsapp.TypeText,
			Text: whatsapp.MessageText{Body: "Thanks4 110,000 https://www.facebook.com/14706725/posts/pfbid06j9pV3uWriQ1234vVWW4FUJ2aNA1fCiJREBnRB6l/?mibextid=crbbb"}}).Run()
	require.NoError(t, err)
	require.NotEmpty(t, message)
}
