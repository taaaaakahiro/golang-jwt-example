package entity

type User struct {
	UserID   string `json:"user_id" bson:"user_id"`
	Password string `json:"password" bson:"password"`
	Name     string `json:"name" bson:"name"`
}
