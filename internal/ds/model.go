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
	Id              string  `json:"id" datastore:"id"`
	DemonstrationId string  `json:"demonstration_id" datastore:"demonstration_id"`
	Polygon         string  `json:"polygon" datastore:"polygon"`
	Latitude        float64 `json:"latitude" datastore:"latitude"`
	Longitude       float64 `json:"longitude" datastore:"longitude"`
}
