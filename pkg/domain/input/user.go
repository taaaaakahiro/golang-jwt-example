package input

type User struct {
	ID       int    `bson:"id"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
}
