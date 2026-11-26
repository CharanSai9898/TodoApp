package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string               `bson:"email" json:"email"`
	Password string               `bson:"password,omitempty" json:"-"`
	GoogleID string               `bson:"googleId,omitempty" json:"googleId,omitempty"`
	Todos    []primitive.ObjectID `bson:"todos" json:"todos"`
}
