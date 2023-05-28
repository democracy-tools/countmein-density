package action

import (
	"errors"
	"fmt"
	"time"

	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/whatsapp"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CreateUser struct {
	dsc   ds.Client
	wac   whatsapp.Client
	phone string
	name  string
}

func NewCreateUser(dsc ds.Client, wac whatsapp.Client, phone string, name string) Request {

	return &CreateUser{
		dsc:   dsc,
		wac:   wac,
		phone: phone,
		name:  name,
	}
}

func (a *CreateUser) Run() (string, error) {

	if err := doCreateUser(a.dsc, a.wac, a.phone, a.name); err != nil {
		return "", fmt.Errorf("failed to add user %s (%s) with %d", a.name, a.phone, err)
	}
	return fmt.Sprintf("user added: %s (%s)", a.name, a.phone), nil
}

func doCreateUser(dsc ds.Client, wac whatsapp.Client, phone string, name string) error {

	err := validateUser(dsc, phone, name)
	if err != nil {
		return err
	}

	id := uuid.NewString()
	err = dsc.Put(ds.KindUser, id, &ds.User{
		Id:         id,
		Phone:      phone,
		Name:       name,
		Preference: "",
		Time:       time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	return wac.SendOnboardingTemplate(phone, id)
}

func validateUser(dsc ds.Client, phone string, name string) error {

	user, err := ds.GetUserByPhone(dsc, phone)
	if err != nil {
		return err
	}
	if user != nil {
		message := fmt.Sprintf("user %s (%s) already exist", name, phone)
		logrus.Infof(message)
		return errors.New(message)
	}

	return nil
}
