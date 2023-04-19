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
	count := len(users)
	if count > 1 {
		err := fmt.Errorf("more than 1 user with the same phone '%s'", phone)
		logrus.Error(err.Error())
		return nil, err
	}
	if count == 1 {
		return &users[0], nil
	}

	return nil, nil
}

func GetKaplanDemonstration(dsc Client) (*Demonstration, error) {

	var demonstration Demonstration
	err := dsc.Get(KindDemonstration, DemonstrationKaplan, &demonstration)
	if err != nil {
		return nil, err
	}

	return &demonstration, nil
}
