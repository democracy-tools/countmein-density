package ds

import "fmt"

type Demonstration struct {
	Id   string `json:"id" datastore:"id"`
	Name string `json:"name" datastore:"name"`
}

type User struct {
	Id         string `json:"id" datastore:"id"`
	Phone      string `json:"phone" datastore:"phone"`
	Name       string `json:"name" datastore:"name"`
	Preference string `json:"preference" datastore:"preference"`
	Time       int64  `json:"time" datastore:"time"`
}

type Volunteer struct {
	UserId          string `json:"user_id" datastore:"user_id"`
	DemonstrationId string `json:"demonstration_id" datastore:"demonstration_id"`
	Polygon         string `json:"polygon" datastore:"polygon"`
	Location        string `json:"location" datastore:"location"`
	Time            int64  `json:"time" datastore:"time"`
}

type Observation struct {
	Time          int64   `json:"time" datastore:"time"`
	User          string  `json:"user_id" datastore:"user_id"`
	Demonstration string  `json:"demonstration" datastore:"demonstration"`
	Polygon       string  `json:"polygon" datastore:"polygon"`
	Density       float32 `json:"density" datastore:"density"`
	Latitude      float32 `json:"latitude" datastore:"latitude"`
	Longitude     float32 `json:"longitude" datastore:"longitude"`
}

func GetVolunteerId(demonstrationId string, userId string) string {

	return fmt.Sprintf("%s$%s", demonstrationId, userId)
}
