package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserObject struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Password []byte             `bson:"password"`
}

func (m *MongoModel) CreateUser(name string, password string) (string, error) {
	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// Insert
	ctx, cancel := defaultContext()
	defer cancel()
	res, err := m.colUser.InsertOne(ctx,
		UserObject{ID: primitive.NewObjectID(), Name: name, Password: passwordHash})
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
