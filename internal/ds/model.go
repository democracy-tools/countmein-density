package ds

import "fmt"

type RegisterRequest struct {
	Phone      string `json:"phone" datastore:"phone"`
	Name       string `json:"name" datastore:"name"`
	Preference string `json:"preference" datastore:"preference"`
	Time       int64  `json:"time" datastore:"time"`
}

type User struct {
	Id         string `json:"id" datastore:"id"`
	Phone      string `json:"phone" datastore:"phone"`
	Name       string `json:"name" datastore:"name"`
	Preference string `json:"preference" datastore:"preference"`
	Time       int64  `json:"time" datastore:"time"`
}

type Volunteer struct {
	Id              string `json:"id" datastore:"id"`
	DemonstrationId string `json:"demonstration_id" datastore:"demonstration_id"`
	Polygon         string `json:"polygon" datastore:"polygon"`
	Location        string `json:"location" datastore:"location"`
	Time            int64  `json:"time" datastore:"time"`
}

func GetVolunteerId(demonstrationId string, userId string) string {

	return fmt.Sprintf("%s$%s", demonstrationId, userId)
}
