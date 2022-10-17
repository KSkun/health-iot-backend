package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserObject struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Password []byte             `bson:"password"`
}

func (m *mongoModel) CreateUser(name string, password string) (primitive.ObjectID, error) {
	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	// Insert
	ctx, cancel := defaultContext()
	defer cancel()
	res, err := m.colUser.InsertOne(ctx,
		UserObject{ID: primitive.NewObjectID(), Name: name, Password: passwordHash})
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (m *mongoModel) GetUserByName(name string) (UserObject, bool, error) {
	ctx, cancel := defaultContext()
	defer cancel()
	res := m.colUser.FindOne(ctx, bson.M{"name": name})
	if res.Err() == mongo.ErrNoDocuments {
		return UserObject{}, false, nil
	}
	if res.Err() != nil {
		return UserObject{}, false, res.Err()
	}
	obj := UserObject{}
	err := res.Decode(&obj)
	if err != nil {
		return UserObject{}, false, err
	}
	return obj, true, nil
}
