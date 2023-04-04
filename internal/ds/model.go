package ds

type RegisterRequest struct {
	Phone string `json:"phone" datastore:"phone"`
	Name  string `json:"name" datastore:"name"`
	Time  int64  `json:"time" datastore:"time"`
}

type User struct {
	Phone string `json:"phone" datastore:"phone"`
	Name  string `json:"name" datastore:"name"`
	Time  int64  `json:"time" datastore:"time"`
}
