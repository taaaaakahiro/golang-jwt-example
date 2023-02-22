package input

type User struct {
	UserID   string `bson:"user_id"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
}
