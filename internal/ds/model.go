package ds

type User struct {
	Email    string `json:"email" datastore:"email"`
	Name     string `json:"name" datastore:"name"`
	Time     int64  `json:"time" datastore:"time"`
	Verified bool   `json:"verified" datastore:"verified"`
}
