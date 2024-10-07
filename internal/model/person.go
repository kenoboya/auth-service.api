package model

type Person struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Age       *int8  `json:"age" bson:"age"`
	Phone     string `json:"phone" bson:"phone"`
}
