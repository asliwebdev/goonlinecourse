package models

type User struct {
	ID    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Role  string `json:"role" bson:"role"`
	Email string `json:"email" bson:"email"`
}