package ds

type RegisterRequest struct {
	Phone string `json:"phone" datastore:"phone"`
	Name  string `json:"name" datastore:"name"`
	Time  int64  `json:"time" datastore:"time"`
}

type User struct {
	Id    string `json:"id" datastore:"id"`
	Phone string `json:"phone" datastore:"phone"`
	Name  string `json:"name" datastore:"name"`
	Time  int64  `json:"time" datastore:"time"`
}

type Volunteer struct {
	Id              string `json:"id" datastore:"id"`
	DemonstrationId string `json:"demonstration_id" datastore:"demonstration_id"`
	Polygon         string `json:"polygon" datastore:"polygon"`
	Location        string `json:"location" datastore:"location"`
}

type Preference struct {
	UserId  string `json:"id" datastore:"id"`
	Polygon string `json:"polygon" datastore:"polygon"`
}
