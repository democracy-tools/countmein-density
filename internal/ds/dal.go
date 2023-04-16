package ds

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func GetUserByPhone(dsc Client, phone string) (*User, error) {

	var users []User
	err := dsc.GetFilter(KindUser, []FilterField{{Name: "phone", Operator: "=", Value: phone}}, &users)
	if err != nil {
		logrus.Errorf("failed to get user by phone '%s' with '%v'", phone, err)
		return nil, err
	}
	if len(users) > 1 {
		err := fmt.Errorf("more than 1 user with the same phone '%s'", phone)
		logrus.Error(err.Error())
		return nil, err
	}

	return nil, nil
}
