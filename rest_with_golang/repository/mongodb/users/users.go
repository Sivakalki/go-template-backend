package user_repo

import (
	"context"
	"rest_with_mongo/db/users"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository{
	return &UserRepository{collection: db.Collection("users")}
}

func (r *UserRepository) Create(ctx context.Context, user *users.User)(*users.User, error){
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, user)

	if err!=nil{
		return nil, err
	}

	return user, nil
}


func (r *UserRepository) GetById(ctx context.Context, id primitive.ObjectID)(*users.User, error){
	var user users.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	if err!=nil{
		return nil, err
	}
	return &user, nil
}


func (r *UserRepository) GetUserByEmail(ctx context.Context, email string)(bool, error){
	// var count int
	// var err error
	count,err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err!=nil{
		return false, err
	}

	if count>0{
		return true, nil
	}
	return false,nil
}


func (r *UserRepository) GetUserByName(ctx context.Context, name string)(bool, error){
	// var count int
	// var err error
	count,err := r.collection.CountDocuments(ctx, bson.M{"username": name})
	if err!=nil{
		return false, err
	}

	if count>0{
		return true, nil
	}
	return false,nil
}


func (r *UserRepository) GetUserByEmailFull(ctx context.Context, email string)(*users.User, error){
	var user users.User
	err := r.collection.FindOne(ctx, bson.M{"email":email}).Decode(&user)

	if err!=nil{
		return nil, err
	}

	return &user, nil
}