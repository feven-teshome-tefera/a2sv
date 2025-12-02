package models

type User struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"-"`
	Role     string `bson:"role" json:"role"`
}
