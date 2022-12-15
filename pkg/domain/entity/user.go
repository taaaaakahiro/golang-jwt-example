package entity

type User struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
