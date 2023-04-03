package ds

type RegisterRequest struct {
	Name  string `json:"name" datastore:"name"`
	Email string `json:"email" datastore:"email"`
	Time  int64  `json:"time" datastore:"time"`
}

type User struct {
	Email string `json:"email" datastore:"email"`
	Name  string `json:"name" datastore:"name"`
	Time  int64  `json:"time" datastore:"time"`
}
