package entity

type User struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
