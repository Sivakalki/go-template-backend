package user_services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rest_with_mongo/db/users"
	"rest_with_mongo/repository/kafka"
	"rest_with_mongo/utils/hash"
	"rest_with_mongo/utils/jwt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	Create(ctx context.Context, user *users.User) (*users.User, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*users.User, error)
	GetUserByEmail(ctx context.Context, email string) ( bool,error)
	GetUserByEmailFull(ctx context.Context, email string) ( *users.User,error)
	GetUserByName(ctx context.Context, name string) ( bool,error)
}


type UserService struct{
	userRepo UserRepo
	jwtGen *jwt.ApxJwt
	producer *kafka.Producer
}

func NewUserService(repo UserRepo, jwtGen *jwt.ApxJwt, producer *kafka.Producer)(*UserService){
	return &UserService{userRepo: repo, jwtGen: jwtGen,producer: producer}
}


type InputUser struct{
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`	
}

func ToUser(user *InputUser)(*users.User){
	return  &users.User{
		Username: user.Username,
		Email: user.Email,
		Password: user.Password,
	}
}


func(s *UserService) CreateUser(ctx context.Context, user *InputUser)(*users.User, error){
	
if user.Username == "" || user.Email == "" || user.Password == "" {
		return nil,errors.New("missing required fields")
	}

	
	if user.Password != user.ConfirmPassword {
		return nil,errors.New("passwords are not matching")
	}
	
	u, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if u == true {
		return nil, errors.New("email already registered")
	}	
	u, err = s.userRepo.GetUserByName(ctx, user.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if u == true  {
		return nil, errors.New("username already registered")
	}
		

	hashed_pwd,err := hash.Encrypt(user.Password)
	if err!= nil{
		return nil,errors.New("Unable to hash the password")
	}

	user.Password = hashed_pwd

	user2, err := s.userRepo.Create(ctx, ToUser(user))
	if err != nil {
		return nil,err
	}

	value := fmt.Sprintf("{email:%s, username:%s}", user.Email, user.Username)

	err = s.producer.Publish(ctx, user.Email, value)
	if err != nil {
        log.Fatal("failed to publish:", err)
    }
	fmt.Println("user log is published successfully, please wait some time mail will be sent")

	// s.logger.Info("User created successfully", zap.String("email", user.Email), zap.Uint("id", *id))
	return user2, nil
}



func (svc *UserService) Login(ctx context.Context, email string, password string)(string, error){
	if(email == ""){
		return "", errors.New("email should not bye empty")
	}
	user,err :=  svc.userRepo.GetUserByEmailFull(ctx, email)
	if err!=nil{
		if err == mongo.ErrNoDocuments {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	if(!hash.Compare(user.Password, password)){
		return "", errors.New("invalid email or password")
	}

	token,err := svc.jwtGen.GenerateJwtToken(user.ID.Hex(),24*time.Hour)

	if err!=nil{
		return "",err
	}

	return token,nil

}


