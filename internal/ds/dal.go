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

func GetVolunteers(dsc Client, demonstration string) ([]Volunteer, error) {

	var volunteers []Volunteer
	err := dsc.GetFilter(KindVolunteer,
		[]FilterField{{Name: "demonstration_id", Operator: "=", Value: demonstration}},
		&volunteers)
	if err != nil {
		return nil, err
	}

	return volunteers, nil
}

func GetObservations(dsc Client, demonstration string) ([]Observation, error) {

	var observations []Observation
	err := dsc.GetFilter(KindObservation,
		[]FilterField{{Name: "demonstration", Operator: "=", Value: demonstration}},
		&observations)
	if err != nil {
		return nil, err
	}

	return observations, nil
}
func IsAdmin(dsc Client, phone string) (bool, error) {

	user, err := GetUserByPhone(dsc, phone)
	if err != nil {
		return false, err
	}
	if user.Role == "admin" {
		return true, nil
	}

	return false, nil
}
