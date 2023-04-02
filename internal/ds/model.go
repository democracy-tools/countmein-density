package ds

type User struct {
	Email    string `json:"email" datastore:"email"`
	Name     string `json:"name" datastore:"name"`
	Time     int64  `json:"time" datastore:"time"`
	Verified int64  `json:"verified" datastore:"verified"`
}
