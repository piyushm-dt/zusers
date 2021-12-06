package employee

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employee struct {
	ID primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName string `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	DOJ string `json:"date_of_joining" bson:"date_of_joining"`
	Skills []string `json:"skills" bson:"skills"`
	Designation string `json:"designation" bson:"designation"`
}
