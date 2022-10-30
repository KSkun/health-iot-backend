package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserObject struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`    // Unique index
	Name     string             `bson:"name" json:"name"` // Unique index
	Password []byte             `bson:"password" json:"-"`
}
